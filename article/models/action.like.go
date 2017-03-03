package models

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func (action *LikeAction) Save() error {
	key := fmt.Sprintf("index://articles/%d/like_count", action.Target)
	_, err := cache.Do("INCRBY", key, 1)
	return err
}

func (action *LikeAction) Delete() error {
	key := fmt.Sprintf("index://articles/%d/like_count", action.Target)
	_, err := cache.Do("INCRBY", key, -1)
	return err
}

func GetArticleLikeCount(id uint64) (int, error) {
	key := fmt.Sprintf("index://articles/%d/like_count", id)
	return redis.Int(cache.Do("GET", key))
}
