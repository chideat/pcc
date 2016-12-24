package models

import (
	"regexp"
	"time"

	"github.com/chideat/glog"
	. "github.com/chideat/pcc/pig/models"
	"github.com/chideat/pcc/sdk/pig"
	"github.com/chideat/pcc/user/modules/auth"
	. "github.com/chideat/pcc/user/modules/error"
	"github.com/dgrijalva/jwt-go"
)

const (
	TOKEN_EXPIRE = 60
)

var (
	name_reg     *regexp.Regexp
	password_reg *regexp.Regexp
)

func init() {
	name_reg, _ = regexp.Compile("^[^'\"\\s]{1,64}$")
	password_reg, _ = regexp.Compile("^[^\\s]{6,64}$")
}

type User struct {
	ID        int64 `json:"id" gorm:"unique_index:idx_user_id"`
	Ban       int   `json:"ban" gorm:"index:index_user_ban"`
	Level     int64 `json:"level" gorm:"index:index_user_level"`
	Anonymous bool  `json:"anonymous"`

	Icon  string `json:"icon"`
	Name  string `json:"name" gorm:"unique_index:idx_user_name"`
	Sex   string `json:"sex"`
	Desc  string `json:"desc" gorm:"type:text"`
	Phone string `json:"phone" gorm:"index:idx_user_phone"`
	Email string `json:"email" gorm:"index:idx_user_email"`

	Password string    `json:"password" gorm:"size:255;not null"`
	Salt     string    `json:"salt" gorm:"size:255;not null"`
	Token    string    `json:"token" gorm:"type:text;index:idx_user_token"`
	Expire   time.Time `json:"expire"`

	// access info
	Online        bool  `json:"online"`
	LastLoginUtc  int64 `json:"last_login_utc"`
	LastLogoutUtc int64 `json:"last_logout_utc"`

	ModifiedUtc int64 `json:"modified_utc"`
	CreatedUtc  int64 `json:"created_utc"`
}

func (user *User) UpdateToken() error {
	token := jwt.New(jwt.SigningMethodHS256)

	timestamp := time.Now()
	expireTimestamp := timestamp.Add(time.Hour * 24 * TOKEN_EXPIRE)
	tokenInfo := auth.UserTokenInfo{
		ID:      user.ID,
		Name:    user.Name,
		Created: timestamp,
		Expired: expireTimestamp,
	}
	tokenStr := tokenInfo.String()
	token.Claims["info"] = tokenStr

	var err error
	user.Token, err = token.SignedString([]byte(auth.JWT_AUTH_KEY))
	user.Expire = expireTimestamp
	return err
}

func (user *User) Login() {
	user.Online = true
	user.LastLoginUtc = time.Now().Local().Unix() / int64(time.Millisecond)
}

func (user *User) Logout() {
	user.Online = false
	user.LastLogoutUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)
}

func (user *User) ResetPassword(password string) {
	user.Password, user.Salt = auth.PasswordEncrypt(password)
}

func (user *User) Save() error {
	var err error

	user.ModifiedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)
	if user.ID == 0 {
		user.ID, err = pig.Int64(TYPE_USER)
		if err != nil {
			return NewUserError("100001", "系统错误")
		}
		user.CreatedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)
		err = db.Create(user).Error
	} else {
		err = db.Save(user).Error
	}
	return err
}

func (user *User) BaseInfo() map[string]interface{} {
	info := map[string]interface{}{}
	info["id"] = user.ID
	info["name"] = user.Name
	info["icon"] = user.Icon
	info["desc"] = user.Desc

	return info
}

func (user *User) Info() map[string]interface{} {
	info := user.BaseInfo()
	info["sex"] = user.Sex
	info["phone"] = user.Phone
	info["email"] = user.Email
	return info
}

func GetUserById(id int64) (*User, error) {
	user := User{}

	err := db.First(&user, id).Error
	if err == ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

func GetUserByName(name string) (*User, error) {
	user := &User{}

	err := db.Where("deleted=false and name=?", name).First(user).Error
	if err == ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func NewUser(name, password string) (*User, error) {
	if !name_reg.MatchString(name) {
		return nil, NewUserError("201002", "用户名无效")
	}
	if !password_reg.MatchString(password) {
		return nil, NewUserError("201003", "密码无效")
	}

	user, err := GetUserByName(name)
	if err != nil {
		glog.Error(err.Error())
		return nil, NewUserError("100001", "系统错误")
	} else if user != nil {
		return nil, NewUserError("201005", "用户名已经存在")
	}

	user.Name = name
	user.Password, user.Salt = auth.PasswordEncrypt(password)

	err = user.Save()
	if err != nil {
		return nil, NewUserError("100001", "创建用户失败")
	}
	return user, nil
}
