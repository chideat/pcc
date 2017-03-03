package pig

import "time"

const (
	_sys_timestamp_ uint64 = 1483200000000000000
)

func Next(group uint8, typ TYPE) uint64 {
	t := uint64(time.Now().UnixNano()) - _sys_timestamp_
	return (t << 16) | uint64(group)<<8 | uint64(typ)
}
