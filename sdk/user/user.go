package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chideat/glog"
)

var (
	addr   string
	client *http.Client
)

func init() {
	addr = "127.0.0.1:7020"

	client = &http.Client{}
}

func SetHost(host, port string) {
	addr = fmt.Sprintf("%s:%s", host, port)
}

func UserBaseInfo(id int64) (map[string]interface{}, error) {
	urlStr := fmt.Sprintf("http://%s/api/v1/users/%d", addr, id)
	res, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		dataRaw, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var data struct {
			Code    string                 `json:"code"`
			Message string                 `json:"message"`
			Data    map[string]interface{} `json:"data"`
		}
		err = json.Unmarshal(dataRaw, &data)
		if err != nil {
			glog.Error(err)
			return nil, err
		}
		if data.Code == "0" {
			return data.Data, nil
		} else {
			return nil, errors.New(data.Message)
		}
	} else {
		return nil, errors.New(res.Status)
	}
}
