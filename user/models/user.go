package models

import (
	"fmt"
	"regexp"
	"time"

	"github.com/chideat/glog"
	. "github.com/chideat/pcc/user/modules/error"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
)

var (
	name_reg     *regexp.Regexp
	password_reg *regexp.Regexp
)

func init() {
	name_reg, _ = regexp.Compile("^[^'\"\\s]{1,64}$")
	password_reg, _ = regexp.Compile("^[^\\s]{6,64}$")
}

func (user *User) Save() error {
	defer user.cache()

	var err error

	user.ModifiedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)
	if user.Id == 0 || user.CreatedUtc == 0 {
		// user.Id = pig.Next(Conf.Group, pig.TYPE_USER)
		user.CreatedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)
		err = db.Create(user).Error
	} else {
		err = db.Save(user).Error
	}
	return err
}

func (user *User) cache() {
	data, _ := proto.Marshal(user)
	_, err := cache.Do("SET", fmt.Sprintf("index://users/%d", user.Id), data)
	if err != nil {
		glog.Error(err)
	}
	_, err = cache.Do("SET", fmt.Sprintf("index://users?name=%s", user.Name), data)
	if err != nil {
		glog.Error(err)
	}
}

func (user *User) Info() map[string]interface{} {
	info := map[string]interface{}{}
	info["id"] = user.Id
	info["name"] = user.Name

	return info
}

func GetUserById(id uint64) (*User, error) {
	user := User{}

	// get user from cache
	data, err := redis.Bytes(cache.Do("GET", fmt.Sprintf("index://users/%d", id)))
	if err == nil {
		err = proto.Unmarshal(data, &user)
		if err == nil {
			return &user, nil
		} else {
			glog.Error(err)
		}
	} else {
		glog.Error(err)
	}

	err = db.First(&user, id).Error
	if err == ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		go user.cache()

		return &user, nil
	}
}

func GetUserByName(name string) (*User, error) {
	user := User{}

	// get user from cache
	data, err := redis.Bytes(cache.Do("GET", fmt.Sprintf("index://users?name=%s", name)))
	if err == nil {
		err = proto.Unmarshal(data, &user)
		if err == nil {
			return &user, nil
		} else {
			glog.Error(err)
		}
	} else {
		glog.Error(err)
	}

	err = db.Where("name=?", name).First(&user).Error
	if err == ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		go user.cache()

		return &user, nil
	}
}

func NewUser(id uint64, name, password string) (*User, error) {
	// if !name_reg.MatchString(name) {
	// 	return nil, NewUserError("201002", "用户名无效")
	// }
	// if !password_reg.MatchString(password) {
	// 	return nil, NewUserError("201003", "密码无效")
	// }

	user := &User{}
	user.Id = id
	user.Name = name
	user.Password = password

	err := user.Save()
	if err != nil {
		return nil, NewUserError("100001", "创建用户失败")
	}
	return user, nil
}
