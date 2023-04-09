package main

import (
	"hot-cache/common/logs"
	"hot-cache/global"
	"hot-cache/initialize"
	"hot-cache/service"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	err := initialize.InitGlobal("./config.yaml")
	if err != nil {
		logs.Error(err)
		return
	}

	addr, err := net.ResolveTCPAddr("tcp", global.AppConfig.Listen)
	if err != nil {
		logs.Error(err)
		return
	}
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Info("listening:", global.AppConfig.Listen)

	sv := service.NewService()

	go func() {
		for {
			conn, err := listen.AcceptTCP()
			if err != nil {
				break
			}
			logs.Info(conn.RemoteAddr(), "connected")
			err = conn.SetKeepAlive(true)
			if err != nil {
				logs.Warning(err)
			}
			err = conn.SetKeepAlivePeriod(time.Duration(global.AppConfig.Heartbeat) * time.Second)
			if err != nil {
				logs.Warning(err)
			}
			go sv.Handle(conn)
		}
		logs.WarningWithoutStack("listening stop")
	}()

	tmp := <-sig
	err = listen.Close()
	if err != nil {
		logs.Warning(err)
	}

	logs.Info("over by", tmp)
}
