package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"runtime/trace"

	"github.com/natexornate/goray/goray"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var tracefile = flag.String("tracefile", "", "write trace profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *tracefile != "" {
		f, err := os.Create(*tracefile)
		if err != nil {
			log.Fatal(err)
		}
		trace.Start(f)
		defer trace.Stop()
	}

	fmt.Printf("Rendering Scene!\n")
	goray.Scene()

	return
}
