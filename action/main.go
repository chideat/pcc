package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/chideat/glog"
	. "github.com/chideat/pcc/action/modules/config"
	"github.com/chideat/pcc/action/routes"
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

func main() {
	httpAddr := Config.HttpAddress
	flag.StringVar(&httpAddr, "httpaddr", httpAddr, "http address")
	flag.BoolVar(&version, "version", false, "version info")

	if !flag.Parsed() {
		flag.Parse()
	}

	if version {
		Version()
		os.Exit(0)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	Config.Model = os.Getenv("DEBUG")

	go func() {
		time.Sleep(time.Millisecond * 50)
		// record pid file
		pidFilePath := path.Join(Config.LogPath, Config.Name+".pid")
		if Config.Model == "debug" {
			pidFilePath = path.Join(Config.LogPath, Config.Name+"_debug.pid")
		}
		pid := []byte(fmt.Sprintf("%d", os.Getpid()))
		ioutil.WriteFile(pidFilePath, pid, 0666)
		fmt.Printf("Server listen on address %s\n", httpAddr)

		<-signalChan
		os.Remove(pidFilePath)
		os.Exit(0)
	}()

	logFilePath := path.Join(Config.LogPath, "access.log")
	if Config.Model == "debug" {
		logFilePath = path.Join(Config.LogPath, "access_debug.log")
	}
	if logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err == nil {
		hd := handlers.CombinedLoggingHandler(logFile, handlers.ProxyHeaders(routes.Handler))
		if err = gracehttp.Serve(&http.Server{Addr: httpAddr, Handler: hd}); err != nil {
			glog.Error(err)
		}
	} else {
		panic(err)
	}
}
