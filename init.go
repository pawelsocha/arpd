package main

import (
	"flag"

	Log "github.com/pawelsocha/kryptond/logging"
)

var (
	ConfigFile  string
	BindAddress string
)

func init() {
	flag.StringVar(&ConfigFile, "config", "/etc/lms/lms.ini", "Path to lms config file")
	flag.StringVar(&BindAddress, "bind", "localhost:1029", "Bind to address")
	flag.Parse()
	Log.SetLogLevel()
}
