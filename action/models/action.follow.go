package models

import (
	"errors"
	"fmt"
	"time"

	. "github.com/chideat/pcc/action/modules/config"
	"github.com/chideat/pcc/action/modules/pig"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

func (action *FollowAction) _BeforeSave() error {
	if action.UserId == 0 {
		return fmt.Errorf("invalid user id")
	}
	if action.Target == 0 {
		return fmt.Errorf("invalid target id")
	}

	oldAction, err := GetFollowAction(action.UserId, action.Target)
	if err != nil || (oldAction != nil && oldAction.Id != action.Id) {
		return fmt.Errorf("已经关注")
	}
	return nil
}

func (action *FollowAction) Save() error {
	err := action._BeforeSave()
	if err != nil {
		return err
	}

	action.ModifiedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)
	if action.Id == 0 {
		action.Id = pig.Next(Conf.Group, pig.TYPE_ACTION)
		action.CreatedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)
		err = db.Create(action).Error
	} else {
		err = db.Save(action).Error
	}
	go action.cache()
	// save list
	_, err = cache.Do("ZADD", fmt.Sprintf("index://target/%d/follow", action.Target), action.ModifiedUtc, action.Id)
	return err
}

func (action *FollowAction) Delete() error {
	action.Deleted = true
	action.DeletedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)

	err := db.Save(action).Error
	if err != nil {
		return err
	}
	go action.cache()
	_, err = cache.Do("ZREM", fmt.Sprintf("index://target/%d/follow", action.Target), action.Id)
	return err
}

func (action *FollowAction) cache() {
	_, err := cache.Do("SET", fmt.Sprintf("index://target/%d", action.Id), action.Bytes())
	if err != nil {
		glog.Error(err)
	}
	_, err = cache.Do("SET", fmt.Sprintf("index://target?target_id=%d&user_id=%d", action.Target, action.UserId), action.Bytes())
	if err != nil {
		glog.Error(err)
	}
}

func (action *FollowAction) Broadcast(method RequestMethod) error {
	req := Request{Method: method}
	req.Data = action.Bytes()
	data, _ := proto.Marshal(&req)
	return producer.Publish("pcc.action.follow", data)
}

func (action *FollowAction) Map() (map[string]interface{}, error) {
	output := map[string]interface{}{}
	output["id"] = action.Id
	output["target"] = action.Target
	output["created_utc"] = action.CreatedUtc
	user, err := GetUserById(action.UserId)
	if err != nil {
		return nil, err
	}
	output["user"] = user.Info()

	return output, nil
}

func (action *FollowAction) Bytes() []byte {
	data, _ := proto.Marshal(action)
	return data
}

func GetFollowActionById(id uint64) (*FollowAction, error) {
	if TYPE_ACTION != uint8(id&255) {
		return nil, errors.New("invalid id")
	}

	action := FollowAction{}
	data, err := redis.Bytes(cache.Do("GET", fmt.Sprintf("index://target/%d", id)))
	if err == nil {
		err = proto.Unmarshal(data, &action)
		if err == nil {
			return &action, nil
		}
	}

	err = db.Where("deleted=false").First(&action, id).Error
	if err == ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		go action.cache()

		return &action, nil
	}
}

func GetFollowAction(userId, target uint64) (*FollowAction, error) {
	if uint8(userId&255) != uint8(pig.TYPE_USER) {
		return nil, errors.New("invalid user id")
	}

	action := FollowAction{}
	data, err := redis.Bytes(cache.Do("GET", fmt.Sprintf("index://target?target_id=%d&user_id=%d", action.Target, action.UserId)))
	if err == nil {
		err = proto.Unmarshal(data, &action)
		if err == nil {
			return &action, nil
		}
	}

	err = db.Where("user_id=? and target=?", userId, target).First(&action).Error
	if err == ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		go action.cache()

		return &action, nil
	}
}

func GetFollowActions(target int64, count int, cursor uint64) ([]*FollowAction, uint64, error) {
	var actions = []*FollowAction{}

	if cursor == 0 {
		cursor = uint64(1 << 63)
	}
	vals, err := redis.Values(cache.Do("ZREVRANGEBYSCORE", fmt.Sprintf("index://target/%d/follow", target),
		cursor, 0, "LIMIT", 0, count))

	if err == nil {
		for _, val := range vals {
			id, err := redis.Uint64(val, nil)
			if err != nil {
				glog.Error(err)
				continue
			}
			action, err := GetFollowActionById(id)
			if err != nil {
				glog.Error(err)
				continue
			}
			actions = append(actions, action)
			cursor = uint64(action.ModifiedUtc)
		}
		if len(vals) < count {
			cursor = 0
		}
		return actions, cursor, nil
	} else {
		err = db.Where("deleted=false and target=? and modified_utc<", target, cursor).
			Order("modified_utc desc").Limit(count).Find(&actions).Error
		if err != nil {
			return nil, 0, err
		}
		if len(actions) == count {
			cursor = uint64(actions[len(actions)-1].ModifiedUtc)
		} else {
			cursor = 0
		}
		return actions, cursor, nil
	}
}

func NewFollowAction(userId, target uint64) (*FollowAction, error) {
	action := FollowAction{}

	action.Id = pig.Next(Conf.Group, pig.TYPE_ACTION)
	action.UserId = userId
	action.Target = target
	action.CreatedUtc = time.Now().Local().UnixNano() / int64(time.Millisecond)

	return &action, nil
}
