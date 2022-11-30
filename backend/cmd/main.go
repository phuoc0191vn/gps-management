package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	utilities "ctigroupjsc.com/phuocnn/gps-management/uitilities"
)

var (
	configPrefix string
	configSource string

	mode      string
	queueMode bool

	container *Container
)

func main() {
	flag.Parse()
	defer utilities.TimeTrack(time.Now(), fmt.Sprintf("GPS-Management API Service"))
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
			main()
		}
	}()

	// load env
	var config Config
	err := utilities.LoadEnvFromFile(&config, configPrefix, configSource)
	if err != nil {
		log.Fatalln(err)
	}

	// load container
	container, err = NewContainer(config)
	if err != nil {
		log.Fatalln(err)
	}

	if mode == "cmd" {
		RunCmd()
		return
	}

	go runQueue()

	container.Logger().Infof("Listen and serve GPS-Management API at %s\n", container.Config.Binding)
	container.Logger().Fatalln(http.ListenAndServe(container.Config.Binding, NewAPIv1(container)))
}

func runQueue() {
	if !queueMode {
		return
	}

	container.Logger().Infoln("Queue Listening")
	queues := []func(){
		GenerateReport,
	}

	for _, worker := range queues {
		go worker()
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.StringVar(&configPrefix, "configPrefix", "gpsmanagement", "config prefix")
	flag.StringVar(&configSource, "configSource", ".env", "config source")

	flag.StringVar(&mode, "mode", "server", "Mode: server | cmd")
	flag.BoolVar(&queueMode, "queue", true, "Enable schedule mode")
}
