package models

import (
	"errors"
	"fmt"
	"time"

	. "github.com/chideat/pcc/action/modules/config"
	"github.com/chideat/pcc/action/modules/pig"
	"github.com/golang/protobuf/proto"
)

func (action *LikeAction) _BeforeSave() error {
	if action.UserId == 0 {
		return fmt.Errorf("invalid user id")
	}
	if action.Target == 0 {
		return fmt.Errorf("invalid target id")
	}

	oldAction, err := GetLikeAction(action.UserId, action.Target)
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

	action.Index = uint64(time.Now().Local().UnixNano() / int64(time.Millisecond))
	if action.IsFriend {
		action.Index = (1 << 48) | action.Index
	}
	if action.Id == 0 {
		action.Id = pig.Next(Conf.Group, pig.TYPE_ACTION)
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

	return db.Save(action).Error
}

func (action *LikeAction) Broadcast(method RequestMethod) error {
	req := Request{Method: method}
	req.Data = action.Bytes()
	data, _ := proto.Marshal(&req)
	return producer.Publish("pcc.action.like", data)
}

func (action *LikeAction) Map() (map[string]interface{}, error) {
	output := map[string]interface{}{}
	output["id"] = action.Id
	output["target"] = action.Target
	output["created_utc"] = action.CreatedUtc
	output["is_friend"] = action.IsFriend
	user, err := GetUserById(action.UserId)
	if err != nil {
		return nil, err
	}
	output["user"] = user.Info()

	return output, nil
}

func (action *LikeAction) Bytes() []byte {
	data, _ := proto.Marshal(action)
	return data
}

func GetLikeActionById(id int64) (*LikeAction, error) {
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

func GetLikeAction(userId, target uint64) (*LikeAction, error) {
	if uint8(userId&25) != uint8(pig.TYPE_USER) {
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

func GetLikeActions(target int64, count int, cursor uint64) ([]*LikeAction, uint64, error) {
	var actions = []*LikeAction{}

	err := db.Where("deleted=false and target=? and index<", target, cursor).
		Order("index desc").Limit(count).Find(&actions).Error
	if err != nil {
		return nil, 0, err
	}
	if len(actions) == count {
		cursor = actions[len(actions)-1].Index
	} else {
		cursor = 0
	}
	return actions, cursor, nil
}

func NewLikeAction(userId, target uint64) (*LikeAction, error) {
	action := LikeAction{}

	action.Id = pig.Next(Conf.Group, pig.TYPE_ACTION)
	action.UserId = userId
	action.Target = target
	action.CreatedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)

	return &action, nil
}
