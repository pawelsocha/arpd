package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/pawelsocha/kryptond/config"
	"github.com/pawelsocha/kryptond/database"
	. "github.com/pawelsocha/kryptond/logging"
	routeros "github.com/pawelsocha/routeros"
)

type Arpd struct {
	done   chan bool
	result chan *routeros.Reply
	cache  map[string]string
}

func (a *Arpd) Run() {

	a.cache = make(map[string]string)
	go a.Collect()

	socket, err := net.Listen("tcp", BindAddress)
	if err != nil {
		Log.Errorf("Cant bind to %s. Error: %s", BindAddress, err.Error())
		os.Exit(1)
	}

	defer socket.Close()
	Log.Info("started")

	for {
		conn, err := socket.Accept()
		if err != nil {
			Log.Errorf("Error accepting: %s", err.Error())
			return
		}
		Log.Infof("Accept connection from: %s", conn.RemoteAddr().String())
		go a.handleConnection(conn)
	}
}

func (a *Arpd) Collect() {
	for {
		select {
		case <-a.done:
			return
		case result := <-a.result:
			a.processResult(result)
		}
	}
}

func (a *Arpd) Result() chan *routeros.Reply {
	return a.result
}

func (a *Arpd) processResult(result *routeros.Reply) {

	if len(result.Re) <= 0 {
		Log.Debug("Invalid result: %v", result)
		return
	}
	for _, v := range result.Re {
		Log.Debugf("Update: %s = %s", v.Map["mac-address"], v.Map["address"])
		a.cache[v.Map["mac-address"]] = v.Map["address"]
	}
}

func (a *Arpd) handleConnection(conn net.Conn) {
	buff := []byte{}
	for k, v := range a.cache {
		formated := []byte(strings.Replace(fmt.Sprintf("%s %s\n", v, k), ":", "", -1))
		buff = append(buff, formated...)
	}
	conn.Write(buff)
	conn.Close()
}

func main() {
	config, err := config.New(ConfigFile)

	if err != nil {
		Log.Critical("Can't read configuration. Error: ", err)
		return
	}

	routers, err := database.Database(config).GetRoutersList()
	if err != nil {
		Log.Critical("Can't get list of routers from database. Error:", err)
	}
	Log.Infof("routers: %v", routers)

	arpd := &Arpd{
		done:   make(chan bool),
		result: make(chan *routeros.Reply),
	}

	/* TODO:
	w.ExecuteCommand("/ip/arp/print", arpd.Result())
	*/
	arpd.Run()
}
