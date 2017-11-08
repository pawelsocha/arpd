package main

import (
	"flag"
	logging "github.com/op/go-logging"
	"os"
)

var Log = logging.MustGetLogger("arpd")

var (
	ConfigFile  string
	Verbose     bool
	BindAddress string
)

func init() {
	flag.StringVar(&ConfigFile, "config", "/etc/lms/lms.ini", "Path to lms config file")
	flag.BoolVar(&Verbose, "verbose", false, "Enable verbose logging")
	flag.StringVar(&BindAddress, "bind", "localhost:1029", "Bind to address")
	flag.Parse()

	console := logging.NewLogBackend(os.Stdout, "", 0)
	formated := logging.NewBackendFormatter(
		console,
		logging.MustStringFormatter("[%{level:.1s}] %{message}"))

	leveled := logging.AddModuleLevel(formated)
	leveled.SetLevel(logging.INFO, "")

	if Verbose {
		leveled.SetLevel(logging.DEBUG, "")
	}

	Log.SetBackend(leveled)
}
