package service

import (
	"fmt"
	"sync"

	"github.com/chideat/glog"
	"github.com/chideat/pcc/article/models"
	. "github.com/chideat/pcc/article/modules/config"
	"github.com/golang/protobuf/proto"
	"github.com/nsqio/go-nsq"
)

var (
	consumer *nsq.Consumer
)

func init() {
	var err error

	// config
	config := nsq.NewConfig()
	config.MaxAttempts = 2
	consumer, err = nsq.NewConsumer("pcc.action.like", "article", config)
	if err != nil {
		glog.Panic(err)
	}

	consumer.AddHandler(&LikeActionHandler{errorMessages: map[uint64]error{}, lock: &sync.RWMutex{}})
	consumer.SetLogger(nil, nsq.LogLevelError)
	err = consumer.ConnectToNSQLookupd(Conf.MQ.ConsumerHTTPAddress)
	if err != nil {
		panic(err)
	}
}

type LikeActionHandler struct {
	errorMessages map[uint64]error
	lock          *sync.RWMutex
}

func (handler *LikeActionHandler) HandleMessage(msg *nsq.Message) error {
	req := models.Request{}
	err := proto.Unmarshal(msg.Body, &req)
	if err != nil {
		glog.Error(err)
		return nil
	}
	action := models.LikeAction{}
	err = proto.Unmarshal(req.Data, &action)
	if err != nil {
		glog.Error(err)
		return nil
	}
	if action.Id == 0 {
		glog.Error("invalid like action with empty id")
		return nil
	}

	switch req.Method {
	case models.RequestMethod_Add:
		err = action.Save()
	case models.RequestMethod_Delete:
		err = action.Delete()
	default:
		err = fmt.Errorf("unknow type %s", req.Method)
	}

	if err != nil {
		handler.lock.Lock()
		handler.errorMessages[action.Id] = fmt.Errorf("%s %s", req.Method, err.Error())
		handler.lock.Unlock()
	}
	return nil
}

func (handler *LikeActionHandler) LogFailedMessage(msg *nsq.Message) {
	req := models.Request{}
	err := proto.Unmarshal(msg.Body, &req)
	if err != nil {
		glog.Error(err)
		return
	}
	action := models.LikeAction{}
	err = proto.Unmarshal(req.Data, &action)
	if err != nil {
		glog.Error(err)
		return
	}

	handler.lock.Lock()
	err, ok := handler.errorMessages[action.Id]
	delete(handler.errorMessages, action.Id)
	handler.lock.Unlock()
	if ok {
		glog.Error("process like action %d failed with error %s", action.Id, err.Error())
	} else {
		glog.Error("process like action %d failed with unknown error", action.Id)
	}
}
