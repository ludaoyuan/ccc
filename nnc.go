package main

import (
	"conf"
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	inchttp "server/http"
	"server/tcp"
	// _ "github.com/mkevac/debugcharts"
)

var (
	cpufile = "cpu.prof"
	memfile = "mem.prof"
)

func main() {
	// go func() {
	// 	log.Fatal(http.ListenAndServe(":8081", nil))
	// }()

	f, err := os.Create(cpufile)
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	// defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	// defer pprof.StopCPUProfile()

	f2, err := os.Create(memfile)
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f2.Close()
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f2); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func(c chan os.Signal, f *os.File) {
		<-c
		pprof.StopCPUProfile()
		f.Close()
		os.Exit(1)
	}(c, f)

	config := make(conf.Config)
	config.Init("/etc/nnc/nnc.conf")

	go tcp.Run(config["TCP"]["domain"], config["TCP"]["port"])
	inchttp.Run(config["HTTP"]["port"], config["HTTP"]["webroot"])
}
