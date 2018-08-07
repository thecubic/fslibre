package fslibre

import (
	"bytes"
	"github.com/karalabe/hid"
	"log"
)

type LibreReader struct {
	// devices all the way down
	deviceinfo *hid.DeviceInfo
	device     *hid.Device
	serial     string
	version    string
}

func New(di *hid.DeviceInfo) LibreReader {
	// make it
	var li LibreReader
	li = LibreReader{di, nil, "", ""}
	return li
}

const UsbVendor uint16 = 0x1a61
const UsbDevice uint16 = 0x3650

func packetize(rpkt []byte) (byte, []byte) {
	var length int
	length = int(rpkt[1])
	return rpkt[0], rpkt[2 : 2+length]
}

func (lbr *LibreReader) handshake() error {
	handshakeCmds := []byte{0x04, 0x05, 0x15, 0x01}
	var err error
	var pktType byte
	var payload []byte
	var rbuf []byte
	for _, cmd := range handshakeCmds {
		rbuf, err = lbr.rpc(cmd)
		if err != nil {
			log.Printf("failure in handshake step %v: %v", cmd, err)
			return err
		}
		pktType, payload = packetize(rbuf)
		if pktType == 0x06 {
			lbr.serial = string(payload[:])
			log.Printf("serial: %v", lbr.serial)
		}
		if pktType == 53 {
			lbr.version = string(payload[:])
			log.Printf("version: %v", lbr.version)
		}
		log.Printf("completed rpc %v -> %v (%v)", cmd, pktType, len(payload))
	}
	log.Print("handshake succeeded")
	return err
}

func (lbr *LibreReader) send(wbuf []byte) error {
	wcnt, err := lbr.device.Write(wbuf)
	if err != nil {
		log.Printf("error writing: %v", err)
	} else {
		log.Printf("wrote: %v", wcnt)
	}
	return err
}

func (lbr *LibreReader) recv() ([]byte, error) {
	rbuf := bytes.Repeat([]byte{0x00}, 64)
	rcnt, err := lbr.device.Read(rbuf)
	if err != nil {
		log.Printf("error reading: %v", err)
	} else {
		log.Printf("read: %v", rcnt)
	}
	return rbuf, err
}

func (lbr *LibreReader) rpc(cmd byte) ([]byte, error) {
	// device packets are 64bytes, and it needs a id prepended (0)
	var err error
	var rbuf []byte
	packet := bytes.Repeat([]byte{0x00}, 65)
	packet[1] = cmd
	err = lbr.send(packet)
	if err != nil {
		return rbuf, err
	}
	rbuf, err = lbr.recv()
	if err != nil {
		log.Printf("RPC %v failed", cmd)
	} else {
		log.Printf("<- %v", rbuf)
		log.Printf("RPC %v succeeded", cmd)
	}
	return rbuf, err
}

func (lbr *LibreReader) Init() error {
	log.Println("initializing LibreReader")
	err := lbr.handshake()
	return err
}

func (lbr *LibreReader) Open() error {
	log.Println("opening LibreReader")
	device, err := lbr.deviceinfo.Open()
	lbr.device = device
	return err
}

func (lbr *LibreReader) close() error {
	err := lbr.device.Close()
	log.Println("closed LibreReader")
	return err
}
