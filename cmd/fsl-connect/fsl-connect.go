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
	edevs := hid.Enumerate(fslibre.UsbVendor, fslibre.UsbDevice)
	fmt.Fprintf(os.Stderr, "hm: %v\n", edevs)
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

		var sn string
		sn, err = lbr.SerialNumber()
		if err != nil {
			log.Printf("error retrieving serial number: %v", err)
		} else {
			fmt.Printf("serial number: %v\n", sn)
		}

		var swv string
		swv, err = lbr.SwVersion()
		if err != nil {
			log.Printf("error retrieving swversion: %v", err)
		} else {
			fmt.Printf("swversion: %v\n", swv)
		}

		//
		// var arr string
		// arr, err = lbr.Arresult()
		// if err != nil {
		// 	log.Printf("error retrieving arr: %v", err)
		// } else {
		// 	log.Printf("arr: %v", arr)
		// }

		// var dt string
		// dt, err = lbr.DateTime()
		// if err != nil {
		// 	log.Printf("error retrieving dt: %v", err)
		// } else {
		// 	log.Printf("dt: %v", dt)
		// }

		hst, err := lbr.History()
		if err != nil {
			log.Printf("error retrieving history: %v", err)
		} else {
			fmt.Printf("history: %v records\n", len(hst.Records))
			fmt.Printf("latest record: %v\n", hst.Records[0])
		}

		// arr, err := lbr.ArrHistory()
		// if err != nil {
		// 	log.Printf("error retrieving arrhistory: %v", err)
		// } else {
		// 	fmt.Printf("arrhistory: %v records\n", len(arr.Records))
		// 	fmt.Printf("latest record: %v\n", arr.Records[0])
		// }

		var dbrnum int
		dbrnum, err = lbr.Dbrnum()
		if err != nil {
			log.Printf("error retrieving db record number: %v", err)
		} else {
			fmt.Printf("db record number: %v\n", dbrnum)
		}

		// var ptname string
		// ptname, err = lbr.PatientName()
		// if err != nil {
		// 	log.Printf("error retrieving patient name: %v", err)
		// } else {
		// 	log.Printf("patient name: %v", ptname)
		// }
		//
		// var ptid string
		// ptid, err = lbr.PatientId()
		// if err != nil {
		// 	log.Printf("error retrieving patient id: %v", err)
		// } else {
		// 	log.Printf("patient id: %v", ptid)
		// }
		err = lbr.Close()
		if err != nil {
			log.Printf("error during close: %v", err)
		}

	}
	log.Print("finished hidtest")
}
