package main

import (
	"fmt"
	"github.com/karalabe/hid"
	"github.com/thecubic/fslibre"
	"log"
	"os"
)

func main() {
	var err error
	edevs := hid.Enumerate(fslibre.UsbDevice, fslibre.UsbDevice)
	if len(edevs) == 0 {
		fmt.Fprintf(os.Stderr, "no Freestyle Libre Reader Devices found!\n")
		os.Exit(1)
	}
	for _, edev := range edevs {
		lbr := fslibre.New(&edev)
		err = lbr.Open()
		if err != nil {
			log.Printf("error opening LibreReader: %v", err)
			break
		}

		err = lbr.Init()
		if err != nil {
			log.Printf("error initializing LibreReader: %v", err)
		}
	}
	log.Print("finished hidtest")
}
