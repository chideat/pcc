package pig

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net/http"
)

const (
	_TOKEN = "H48K1VlJyqOcLqGmzEwSuuorW4qmHlQuDfY9bLSRKSw="
)

var (
	addr   string
	client *http.Client
)

func init() {
	addr = "127.0.0.1:7000"

	client = &http.Client{}
}

func SetHost(host, port string) {
	addr = fmt.Sprintf("%s:%s", host, port)
}

func Int64(typ uint8) (int64, error) {
	id, err := Uint64(typ)
	return int64(id), err
}

func Uint64(typ uint8) (uint64, error) {
	urlStr := fmt.Sprintf("http://%s/api/v1/id/%d", addr, typ)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", _TOKEN)

	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		var id uint64
		err = binary.Read(res.Body, binary.LittleEndian, &id)
		if err != nil {
			return 0, err
		}
		return id, nil
	} else {
		return 0, errors.New(res.Status)
	}
}
