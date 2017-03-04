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
	"github.com/chideat/pcc/user/models"
	. "github.com/chideat/pcc/user/modules/config"
	"github.com/chideat/pcc/user/routes"
	"github.com/chideat/pcc/user/service"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/handlers"
)

var (
	BuildTimestamp string
	BuildCommit    string
	version        bool
)

func Version() {
	fmt.Printf("Build at %s, based on commit %s\n", BuildTimestamp, BuildCommit)
}

func importUsers(filePath string) error {
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

		user := models.User{}
		user.Id, err = strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			fmt.Println(line, err)
			continue
		}
		user.Name = parts[1]

		err = user.Save()
		if err != nil {
			fmt.Println(line, err)
		}
	}
	return nil
}

func main() {
	var userFilePath string

	flag.BoolVar(&version, "version", false, "version info")
	flag.StringVar(&userFilePath, "u", "", "user file path")

	if !flag.Parsed() {
		flag.Parse()
	}

	if userFilePath != "" {
		_, err := os.Stat(userFilePath)
		if !os.IsNotExist(err) {
			importUsers(userFilePath)
		}
		os.Exit(0)
	}

	if version {
		Version()
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

	go service.StartRPCService(Conf.RPCAddr)

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
