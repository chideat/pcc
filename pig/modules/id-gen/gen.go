package gen

import (
	"fmt"

	"github.com/chideat/glog"
	"github.com/garyburd/redigo/redis"
)

const ()

type IDGen struct {
	pool      *redis.Pool
	prefix    string
	scriptMap string
}

func (gen *IDGen) Next(typ uint8) (uint64, error) {
	conn := gen.pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:%d", gen.prefix, typ)
	if id, err := redis.Uint64(conn.Do("INCR", key)); err != nil {
		return 0, err
	} else {
		return ((id << 8) | uint64(typ)), nil
	}
}

func (gen *IDGen) GetIdByRefer(typ uint8, refer string) (uint64, error) {
	conn := gen.pool.Get()
	defer conn.Close()

	script := fmt.Sprintf(`
	local refer_id_map = "%s:"..KEYS[1]..":refer-id"
	local id_refer_map = "%s:"..KEYS[1]..":id-refer"
	local id = redis.call("HGET", refer_id_map, KEYS[2])
	if id == false or id == nil then
		id = redis.call("INCR", KEYS[1])
		if id ~= nil and id ~= false then
		    id = bit.lshift(tonumber(id), 8) + tonumber(KEYS[1])
		    redis.call("HSET", refer_id_map, KEYS[2], id)
		    redis.call("HSET", id_refer_map, id, KEYS[2])
		end
	end
	return id
`, gen.prefix, gen.prefix)

	if id, err := redis.Uint64(conn.Do("EVAL", script, 2, typ, refer)); err != nil {
		return 0, err
	} else {
		return id, nil
	}
}

func (gen *IDGen) GetReferById(typ uint8, id int64) (string, error) {
	conn := gen.pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:%d:id-refer", gen.prefix, typ)
	refer, err := redis.String(conn.Do("HGET", key, id))
	if err == nil {
		return refer, nil
	} else {
		return "", err
	}
}

func (gen *IDGen) UpdateRefer(typ uint8, id int64, refer string) error {
	conn := gen.pool.Get()
	defer conn.Close()

	script := fmt.Sprintf(`
	local refer_id_map = "%s:"..KEYS[1]..":refer-id"
	local id_refer_map = "%s:"..KEYS[1]..":id-refer"

	redis.call("HSET", id_refer_map, KEYS[2], KEYS[3])
	redis.call("HSET", refer_id_map, KEYS[3], KEYS[2])

	return nil
`, gen.prefix, gen.prefix)

	_, err := conn.Do("EVAL", script, 3, typ, id, refer)
	return err
}

func (gen *IDGen) GetReferByUserId(typ uint8, userId int64) (string, error) {
	conn := gen.pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:%d:id-refer", gen.prefix, typ)
	refer, err := redis.String(conn.Do("HGET", key, userId))
	if err != nil {
		return "", err
	} else {
		return refer, nil
	}
}

func NewIDGen(addr string) *IDGen {
	gen := IDGen{}

	gen.prefix = "PIG"
	gen.pool = redis.NewPool(func() (redis.Conn, error) {
		conn, err := redis.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}

		_, err = conn.Do("PING")
		if err != nil {
			conn.Close()
			return nil, err
		}

		return conn, nil
	}, 10)

	// set redis config
	conn := gen.pool.Get()
	defer conn.Close()
	_, err := conn.Do("CONFIG", "SET", "SAVE", "30 1")
	if err != nil {
		glog.Panic(err)
	}

	return &gen
}
