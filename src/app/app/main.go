package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/golang/glog"
)

type Options struct {
	StaticDir string
	Hostname  string
	Port      int
	Address   string
}

var options Options

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s -dir [static_dir]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.IntVar(&options.Port, "port", 8090, "port")
	flag.StringVar(&options.Hostname, "host", "localhost", "hostname")
	flag.StringVar(&options.StaticDir, "dir", "app", "directory with statics[js,images,css]")
	flag.StringVar(&options.Address, "address", "http://google.com", "address to check")
}

func main() {
	// start web server with as many cpu as possible
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Parse flags
	flag.Parse()

	// run websocket hub
	go h.run()

	// run checker
	go runChecker(options.Address, &h)

	// start HTTP server
	LaunchServer(options)

	glog.Infof("Done")
}
