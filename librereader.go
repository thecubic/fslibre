package fslibre

import (
	"bytes"
	"fmt"
	"github.com/karalabe/hid"
	"log"
	// "os"
	"strconv"
	"strings"
)

type LibreReader struct {
	// devices all the way down
	deviceinfo *hid.DeviceInfo
	device     *hid.Device
	serial     string
	version    string
}

type libreresponse struct {
	text    string
	ok      bool
	cksm    uint64
	cksm_ok bool
}

// type historyrecord struct {
// 	sequence int
// }

type HistoryRecords struct {
	Records [][16]int
}

const UsbVendor uint16 = 0x1a61
const UsbDevice uint16 = 0x3650

// PUBLIC API: DEVICE MANAGEMENT

func New(di *hid.DeviceInfo) LibreReader {
	// make it
	var li LibreReader
	li = LibreReader{di, nil, "", ""}
	return li
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

func (lbr *LibreReader) Close() error {
	err := lbr.device.Close()
	log.Println("closed LibreReader")
	return err
}

// PUBLIC API: DEVICE COMMANDS

func (lbr *LibreReader) SerialNumber() (string, error) {
	resps, err := lbr.text_command("$sn?")
	if resps.ok {
		return resps.text, err
	} else if err != nil {
		return "", err
	} else {
		return "", fmt.Errorf("command failure")
	}
}

func (lbr *LibreReader) SwVersion() (string, error) {
	resps, err := lbr.text_command("$swver?")
	if resps.ok {
		return resps.text, err
	}
	return "", fmt.Errorf("command failure")
}

func (lbr *LibreReader) DateTime() (string, error) {
	date, err := lbr.text_command("$date?")
	if err != nil {
		return "", err
	}

	time, err := lbr.text_command("$time?")
	if err != nil {
		return "", err
	}

	fmt.Printf("date: %v\ntime: %v\n", date.text, time.text)
	return "", err
}

func (lbr *LibreReader) History() (*HistoryRecords, error) {
	var err error
	err = lbr.send_text_command("$history?")
	if err != nil {
		return nil, err
	}
	resp, err := lbr.history_recv()
	return resp, err
}

// TODO: not the same as history
// func (lbr *LibreReader) ArrHistory() (*HistoryRecords, error) {
// 	var err error
// 	err = lbr.send_text_command("$arresult?")
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp, err := lbr.history_recv()
// 	return resp, err
// }

// func (lbr *LibreReader) Arresult() (string, error) {
// 	resps, err := lbr.raw_text_command("$arresult?")
// 	return resps, err
// }

func (lbr *LibreReader) Dbrnum() (int, error) {
	resps, err := lbr.text_command("$dbrnum?")
	if err != nil {
		return -1, err
	}
	if resps.ok {
		if resps.text[:19] == "DB Record Number = " {
			return strconv.Atoi(resps.text[19:])
		} else {
			return -1, fmt.Errorf("response unknown")
		}
	}
	return -1, fmt.Errorf("command failure")
}

// TODO: broken, probably because they are not set and thus special

// func (lbr *LibreReader) PatientName() (string, error) {
// 	resps, err := lbr.text_command("$ptname?")
// 	if resps.ok {
// 		return resps.text, err
// 	}
// 	return "", fmt.Errorf("command failure")
// }
//
// func (lbr *LibreReader) PatientId() (string, error) {
// 	resps, err := lbr.text_command("$ptid?")
// 	if resps.ok {
// 		return resps.text, err
// 	}
// 	return "", fmt.Errorf("command failure")
// }

// PRIVATE

func (lbr *LibreReader) send_text_command(cmd string) error {
	packet := bytes.Repeat([]byte{0x00}, 65)
	packet[1] = 0x60
	packet[2] = byte(len(cmd))
	copy(packet[3:], cmd)
	return lbr.send(packet)
}

// func (lbr *LibreReader) raw_text_command(cmd string) (string, error) {
// 	var err error
// 	err = lbr.send_text_command(cmd)
// 	if err != nil {
// 		return "", err
// 	}
// 	resp, err := lbr.multi_recv()
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(resp), err
// }

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
		}
		if pktType == 53 {
			lbr.version = string(payload[:])
		}
	}
	return err
}

func (lbr *LibreReader) send(wbuf []byte) error {
	_, err := lbr.device.Write(wbuf)
	return err
}

func (lbr *LibreReader) recv() ([]byte, error) {
	rbuf := bytes.Repeat([]byte{0x00}, 64)
	_, err := lbr.device.Read(rbuf)
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
	return rbuf, err
}

func (lbr *LibreReader) history_recv() (*HistoryRecords, error) {
	var err error
	var nrecords int
	// 2000 should be enough
	var records [2000][16]int
	var c_records_cksum uint64
	var records_cksum uint64
	var c_full_cksum uint64
	var full_cksum uint64

	recordidx := 0
	failure := true
	rdone := false

	c_records_cksum = 0
	c_full_cksum = 0

	for !rdone {
		pkt, err := lbr.recv()
		if err != nil {
			return nil, err
		}
		// pktType, izepkt := packetize(pkt)
		resplen := int(pkt[1])
		resptext := string(pkt[2 : 2+resplen])
		var pktcksum uint64
		pktcksum = 0
		for _, chr := range pkt[2 : 2+resplen] {
			pktcksum += uint64(chr)
		}
		respfields := strings.Split(resptext, "\r\n")
		nfields := len(respfields)
		if nfields == 2 {
			var data [16]int
			for i, datum := range strings.Split(respfields[0], ",") {
				data[i], err = strconv.Atoi(datum)
			}
			records[recordidx] = data
			recordidx++
			c_records_cksum += pktcksum
			c_full_cksum += pktcksum
		} else if nfields == 4 {
			// only add the first field
			for _, chr := range pkt[2 : 4+len(respfields[0])] {
				c_full_cksum += uint64(chr)
			}

			endinfo := strings.Split(respfields[0], ",")
			nrecords, err = strconv.Atoi(endinfo[0])
			records_cksum, err = strconv.ParseUint(endinfo[1], 16, 32)

			if respfields[1][:5] == "CKSM:" {
				full_cksum, err = strconv.ParseUint(respfields[1][5:], 16, 32)
			}
			if respfields[2] == "CMD OK" {
				failure = false
			}
			rdone = true
		}
	}
	if failure {
		return nil, fmt.Errorf("general protocol read failure")
	}
	if nrecords != recordidx {
		return nil, fmt.Errorf("number of records mismatch: %v received, %v expected", recordidx, nrecords)
	}
	if c_full_cksum != full_cksum {
		return nil, fmt.Errorf("checksum mismatch: %v received, %v calculated", full_cksum, c_full_cksum)
	}
	if c_records_cksum != records_cksum {
		return nil, fmt.Errorf("records checksum mismatch: %v received, %v calculated", records_cksum, c_records_cksum)
	}
	return &HistoryRecords{records[:nrecords]}, err
}

func (lbr *LibreReader) text_command(cmd string) (*libreresponse, error) {
	var err error

	err = lbr.send_text_command(cmd)
	if err != nil {
		return nil, err
	}

	resp, err := lbr.recv()
	if err != nil {
		return nil, err
	}

	if resp[0] != 0x60 {
		log.Printf("response not text-response (%v)!", resp[0])
	}

	var cksum uint64

	resplen := int(resp[1])
	resptext := string(resp[2 : 2+resplen])
	respfields := strings.Split(resptext, "\r\n")
	respstruct := new(libreresponse)
	respstruct.text = respfields[0]

	if len(respfields) > 1 {
		cksum = 0
		for _, chr := range []byte(respfields[0]) {
			cksum += uint64(chr)
		}
		cksum += uint64('\r')
		cksum += uint64('\n')
		ckv := strings.Split(respfields[1], ":")
		if len(ckv) > 1 {
			rcksum, _ := strconv.ParseUint(ckv[1], 16, 32)
			respstruct.cksm = rcksum
			respstruct.cksm_ok = rcksum == cksum
		} else {
			for _, rstr := range respfields[1:] {
				log.Printf("line: %v\n", rstr)
			}
		}
	}

	if len(respfields) > 2 {
		if respfields[2] == "CMD OK" {
			respstruct.ok = true
		} else {
			respstruct.ok = false
			// err = fmt.Errorf("checkum mismatch: %v (recv) vs %v (calc)", rcksum, cksum)
		}
	}

	return respstruct, err
}
