package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/YouLiao/SDK"
	"io/ioutil"
	"net/http"
)

const (
	_TOKEN = "e8UvMPlzYWjGw0WyO09F/LYZCl76fQsv4AK1mY/IaVM="
)

var (
	addr   string
	client *http.Client
)

func init() {
	addr = "127.0.0.1:7030"

	client = &http.Client{}
}

func SetHost(host, port string) {
	addr = fmt.Sprintf("%s:%s", host, port)
}

func InfoOfPost(id int64) (map[string]interface{}, error) {
	urlStr := fmt.Sprintf("http://%s/api/v1/info/posts/%d", addr, id)

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

		var data map[string]interface{}
		err = json.Unmarshal(dataRaw, &data)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		return nil, errors.New(res.Status)
	}
}

func InfoOfSection(userId, id int64) (map[string]interface{}, error) {
	urlStr := fmt.Sprintf("http://%s/api/v1/info/sections/%d?user_id=%d", addr, id, userId)

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

		var data map[string]interface{}
		err = json.Unmarshal(dataRaw, &data)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		return nil, errors.New(res.Status)
	}
}

func InfoOfUser(id int64) (map[string]interface{}, error) {
	urlStr := fmt.Sprintf("http://%s/api/v1/info/users/%d", addr, id)

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

		var data map[string]interface{}
		err = json.Unmarshal(dataRaw, &data)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		return nil, errors.New(res.Status)
	}
}

func GetPostComments(id, userId int64, page, count, cursor int64, sort string) ([]map[string]interface{}, error) {
	urlStr := fmt.Sprintf("http://%s/api/v1/posts/%d/comments?app=Appwill&ida=1234&user_id=%d&page=%d&count=%d&cursor=%d&sort=%s", addr, id, userId, page, count, cursor, sort)

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
			Code    string                   `json:"code"`
			Message string                   `json:"message"`
			Data    []map[string]interface{} `json:"data"`
		}
		err = json.Unmarshal(dataRaw, &data)
		if err != nil {
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

func GetSubscribedSections(id int64, page, count, cursor int64) ([]int64, error) {
	urlStr := fmt.Sprintf("http://%s/api/v1/users/%d/sections?page=%d&count=%d&cursor=%d", addr, id, page, count, cursor)

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
			Code    string  `json:"code"`
			Message string  `json:"message"`
			Data    []int64 `json:"data"`
		}
		err = json.Unmarshal(dataRaw, &data)
		if err != nil {
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

func IsSectionSubscribed(sectionId, userId int64) (bool, error) {
	urlStr := fmt.Sprintf("http://%s/api/v1/users/%d/sections/%d/subscribed", addr, userId, sectionId)

	res, err := http.Get(urlStr)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return true, nil
	} else if res.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		return false, errors.New(res.Status)
	}
}

func IsUserFollowed(target, userId int64) (bool, error) {
	urlStr := fmt.Sprintf("http://%s/api/v1/users/%d/followed?user_id=%d", addr, target, userId)

	res, err := http.Get(urlStr)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return true, nil
	} else if res.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		return false, errors.New(res.Status)
	}
}

func GetUserExtraInfo(target, userId int64) (map[string]interface{}, error) {
	urlStr := fmt.Sprintf("http://%s/api/v1/users/%d/extra_info?user_id=%d", addr, target, userId)

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

func GetActionById(id int64) (interface{}, error) {
	urlStr := fmt.Sprintf("http://%s/api/v1/actions/%d", addr, id)

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

		var data sdk.ResponseData
		err = json.Unmarshal(dataRaw, &data)
		if err != nil {
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

func GetRecentActions(userId int64, action int32, offset int, page, count, cursor int64) ([]int64, string, error) {
	urlStr := fmt.Sprintf("http://%s/api/v1/recent_actions?user_id=%d&action=%d&offset=%d&page=%d&count=%d&cursor=%d", addr, userId, action, offset, page, count, cursor)

	res, err := http.Get(urlStr)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		dataRaw, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, "", err
		}

		var data struct {
			Code    string  `json:"code"`
			Message string  `json:"message"`
			Data    []int64 `json:"data"`
			Info    struct {
				Cursor   string `json:"cursor"`
				NextPage int64 `json:"next_page"`
			} `json:"info"`
		}
		err = json.Unmarshal(dataRaw, &data)
		if err != nil {
			return nil, "", err
		}
		if data.Code == "0" {
			return data.Data, data.Info.Cursor, nil
		} else {
			return nil, "", errors.New(data.Message)
		}
	} else {
		return nil, "", errors.New(res.Status)
	}
}
