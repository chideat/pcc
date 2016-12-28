package v1

import "github.com/nsqio/go-nsq"

var (
	consumer *nsq.Consumer
)

func init() {
	// 	var err error
	//
	// 	// config
	// 	config := nsq.NewConfig()
	// 	consumer, err = nsq.NewConsumer(ACTION_TOPIC, "default", config)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	consumer.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
	// 		req := models.Request{}
	//
	// 		err := proto.Unmarshal(msg.Body, &req)
	// 		if err != nil || req.Action == nil {
	// 			glog.Error("invalid data")
	// 			return nil
	// 		}
	//
	// 		switch strings.ToUpper(req.Method) {
	// 		case "ADD":
	// 			if req.Action.Type == models.ActionType_Load {
	// 				return nil
	// 			}
	// 			err = req.Action.Save()
	// 			if err != nil {
	// 				glog.Error(err)
	// 			}
	// 		case "DEL":
	// 			action, err := models.GetActionById(req.Action.Id)
	// 			if err != nil {
	// 				glog.Error(err)
	// 			}
	// 			if action != nil && !action.Deleted {
	// 				err = action.Delete()
	// 				if err != nil {
	// 					glog.Error(err)
	// 				}
	// 			}
	// 		}
	// 		return nil
	// 	}))
	//
	// 	consumer.SetLogger(nil, nsq.LogLevelError)
	// 	err = consumer.ConnectToNSQLookupd(Config.Queue.LookupdAddress)
	// 	if err != nil {
	// 		panic(err)
	// 	}
}
