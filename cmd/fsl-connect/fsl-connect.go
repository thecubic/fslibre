package main

import (
	"flag"
	"fmt"
	"github.com/karalabe/hid"
	"github.com/thecubic/fslibre"
	"log"
	"os"
)

var (
	f_sn   = flag.Bool("sn", true, "get serial number")
	f_swv  = flag.Bool("swv", true, "get software version")
	f_arr  = flag.Bool("arr", false, "get associated reports")
	f_hst  = flag.Bool("hst", true, "get history")
	f_dt   = flag.Bool("dt", false, "get datetime")
	f_dn   = flag.Bool("dn", false, "get dbrnum")
	f_ptid = flag.Bool("ptid", false, "get ptid")
	f_ptn  = flag.Bool("ptn", false, "get ptn")
)

func main() {
	var (
		err    error
		ptname string
		ptid   string
		dbrnum int
		dt     string
		arr    string
		swv    string
		sn     string
	)

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

		if *f_sn {
			sn, err = lbr.SerialNumber()
			if err != nil {
				log.Printf("error retrieving serial number: %v", err)
			} else {
				fmt.Printf("serial number: %v\n", sn)
			}
		}

		if *f_swv {
			swv, err = lbr.SwVersion()
			if err != nil {
				log.Printf("error retrieving swversion: %v", err)
			} else {
				fmt.Printf("swversion: %v\n", swv)
			}
		}

		if *f_arr {
			arr, err = lbr.Arresult()
			if err != nil {
				log.Printf("error retrieving arr: %v", err)
			} else {
				log.Printf("arr: %v", arr)
			}
		}

		if *f_dt {
			if err != nil {
				dt, err = lbr.DateTime()
				log.Printf("error retrieving dt: %v", err)
			} else {
				log.Printf("dt: %v", dt)
			}
		}

		if *f_hst {
			hst, err := lbr.History()
			if err != nil {
				log.Printf("error retrieving history: %v", err)
			} else {
				fmt.Printf("history: %v records\n", len(hst.Records))
				latest := hst.Records[0]
				fmt.Printf("latest record: %v\n", latest)
				fmt.Printf("idx:%v type:%v month:%v day:%v year:%v hour:%v minute:%v second:%v value:%v errors:%v",
					latest[0], latest[1], latest[2], latest[3], latest[4], latest[5], latest[6],
					latest[7], latest[13], latest[15])
			}
		}

		if *f_arr {
			arr, err := lbr.ArrHistory()
			if err != nil {
				log.Printf("error retrieving arrhistory: %v", err)
			} else {
				fmt.Printf("arrhistory: %v records\n", len(arr.Records))
				fmt.Printf("latest record: %v\n", arr.Records[0])
			}
		}

		if *f_dn {
			dbrnum, err = lbr.Dbrnum()
			if err != nil {
				log.Printf("error retrieving db record number: %v", err)
			} else {
				fmt.Printf("db record number: %v\n", dbrnum)
			}
		}

		if *f_ptn {
			ptname, err = lbr.PatientName()
			if err != nil {
				log.Printf("error retrieving patient name: %v", err)
			} else {
				log.Printf("patient name: %v", ptname)
			}
		}

		if *f_ptid {
			ptid, err = lbr.PatientId()
			if err != nil {
				log.Printf("error retrieving patient id: %v", err)
			} else {
				log.Printf("patient id: %v", ptid)
			}
		}

		err = lbr.Close()
		if err != nil {
			log.Printf("error during close: %v", err)
		}

	}
	log.Print("finished hidtest")
}
