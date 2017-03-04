package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/chideat/glog"
	"github.com/chideat/pcc/action/models"
	. "github.com/chideat/pcc/action/modules/config"
	"github.com/chideat/pcc/action/modules/pig"
	"github.com/chideat/pcc/action/routes"
	_ "github.com/chideat/pcc/action/service"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/handlers"
)

var (
	BuildTimestamp  string
	BuildCommit     string
	version         bool
	friendsFilePath string
	likeFilePath    string
)

func Version() {
	fmt.Printf("Build at %s, based on commit %s\n", BuildTimestamp, BuildCommit)
}

func importFriends(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = strings.Trim(line, "\r\n")
		parts := strings.Split(line, ",")

		action := models.FollowAction{}
		action.UserId, err = strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			fmt.Println(line, err)
			continue
		}
		action.Target, err = strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			fmt.Println(line, err)
			continue
		}
		action.Id = pig.Next(Conf.Group, pig.TYPE_ACTION)
		err = action.Broadcast(models.RequestMethod_Add)
		if err != nil {
			fmt.Println(line, err)
		}
	}
	return nil
}

func importLike(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	action := models.FollowAction{}
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = strings.Trim(line, "\r\n")
		parts := strings.Split(line, ":")

		action.Target, err = strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			fmt.Println(line, err)
			continue
		}

		for _, userIdStr := range strings.Split(strings.Trim(parts[1], "[]"), ",") {
			action.UserId, err = strconv.ParseUint(userIdStr, 10, 64)
			if err != nil {
				fmt.Println(action.Target, userIdStr, "ignore")
				continue
			}
			action.Id = pig.Next(Conf.Group, pig.TYPE_ACTION)
			err = action.Broadcast(models.RequestMethod_Add)
			if err != nil {
				fmt.Println(action.Target, action.UserId, "ignore")
				continue
			}
		}
	}
	return nil
}

func main() {
	flag.StringVar(&friendsFilePath, "f", "", "friends file path")
	flag.StringVar(&likeFilePath, "l", "", "like file path")
	flag.BoolVar(&version, "version", false, "version info")

	if !flag.Parsed() {
		flag.Parse()
	}

	if version {
		Version()
		os.Exit(0)
	}

	if friendsFilePath != "" {
		err := importFriends(friendsFilePath)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	if likeFilePath != "" {
		err := importLike(likeFilePath)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		time.Sleep(time.Millisecond * 50)
		// record pid file
		pidFilePath := path.Join(Conf.LogPath, Conf.Name+".pid")
		if Conf.IsDebug() {
			pidFilePath = path.Join(Conf.LogPath, Conf.Name+"_debug.pid")
		}
		pid := []byte(fmt.Sprintf("%d", os.Getpid()))
		ioutil.WriteFile(pidFilePath, pid, 0666)
		glog.Infof("Server listen on address %s\n", Conf.HTTPAddr)

		<-signalChan
		os.Remove(pidFilePath)
		os.Exit(0)
	}()

	logFilePath := path.Join(Conf.LogPath, "access.log")
	if Conf.IsDebug() {
		logFilePath = path.Join(Conf.LogPath, "access_debug.log")
	}
	if logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err == nil {
		hd := handlers.CombinedLoggingHandler(logFile, handlers.ProxyHeaders(routes.Handler))
		if err = gracehttp.Serve(&http.Server{Addr: Conf.HTTPAddr, Handler: hd}); err != nil {
			glog.Error(err)
		}
	} else {
		glog.Panic(err)
	}
}
