package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/chideat/pcc/sdk/pig"
	"github.com/chideat/pcc/sdk/user"
	"github.com/golang/protobuf/proto"

	. "github.com/chideat/pcc/pig/models"
)

func (action *LikeAction) _BeforeSave() error {
	if action.UserId == 0 {
		return fmt.Errorf("invalid user id")
	}
	if action.Target == 0 {
		return fmt.Errorf("invalid target id")
	}

	oldAction, err := GetLikeActionByUserAndTarget(action.UserId, action.Target)
	if err != nil || (oldAction != nil && oldAction.Id != action.Id) {
		return fmt.Errorf("重复点赞")
	}
	return nil
}

func (action *LikeAction) Save() error {
	err := action._BeforeSave()
	if err != nil {
		return err
	}

	if action.Id == 0 {
		action.Id, err = pig.Int64(TYPE_ACTION)
		if err != nil {
			return errors.New("系统错误")
		}
		action.CreatedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)
		err = db.Create(action).Error
	} else {
		err = db.Save(action).Error
	}

	return err
}

func (action *LikeAction) Delete() error {
	action.Deleted = true
	action.DeletedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)

	return action.Save()
}

func (action *LikeAction) Map() (map[string]interface{}, error) {
	output := map[string]interface{}{}
	output["id"] = action.Id
	output["target"] = action.Target
	output["mood"] = action.Mood.String()
	output["created_utc"] = action.CreatedUtc

	info, err := user.UserBaseInfo(action.UserId)
	if err != nil {
		return nil, err
	}
	output["user"] = info

	return output, nil
}

func (action *LikeAction) Bytes() []byte {
	data, _ := proto.Marshal(action)

	return data
}

func GetActionById(id int64) (*LikeAction, error) {
	if TYPE_ACTION != uint8(id&255) {
		return nil, errors.New("invalid id")
	}

	action := LikeAction{}
	err := db.Where("deleted=false").First(&action, id).Error
	if err == ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &action, nil
	}
}

func GetLikeActionByUserAndTarget(userId, target int64) (*LikeAction, error) {
	if TYPE_USER != uint8(userId&25) {
		return nil, errors.New("invalid user id")
	}

	action := LikeAction{}
	err := db.Where("user_id=? and target=?", userId, target).First(&action).Error
	if err == ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &action, nil
	}
}

func NewLikeAction(userId, target int64, mood LikeMood) *LikeAction {
	action := LikeAction{}

	action.UserId = userId
	action.Target = target
	action.Mood = mood

	return &action
}
