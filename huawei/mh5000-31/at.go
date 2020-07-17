package mh5000_31

import (
	"encoding/json"
	"fmt"
	"net"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	devices "github.com/JackyCZJ/go-at"

	"github.com/albenik/go-serial/v2"
)

//Init if device is nil , use default.
func (h *HUAWEI) Init(portName string, opt ...serial.Option) {
	defer func() {
		if r := recover(); r != nil {
			panic(r.(error))
		}
	}()
	if len(opt) < 1 {
		opt = []serial.Option{
			serial.WithDataBits(8),
			serial.WithStopBits(serial.OneStopBit),
			serial.WithBaudrate(devices.BaudRate),
			serial.WithParity(serial.OddParity),
		}
	}
	h.Device = &devices.Device{
		Opt:   opt,
		Mutex: &sync.Mutex{},
	}
	h.Device.SetPortName(portName)
	h.Device.Result = h.Result
	//AT ECHO ON!
	_, err := h.Device.Cmd("ATE")
	if err != nil {
		panic(err)
	}
}

type HUAWEI struct {
	*devices.Device
}

//Result Deal with buff result.
func (h *HUAWEI) Result(str string) (bool, error) {
	ok := strings.Contains(str, "OK")
	switch true {
	case strings.Contains(str, "TOO MANY PARAMETERS") || strings.Contains(str, "COMMAND NOT SUPPORT") || strings.Contains(str, "NO CARRIER"), strings.Contains(str, "ERROR"):
		return false, fmt.Errorf(str)
	case ok:
		return true, nil
	}
	return false, nil
}

//ATI ask minor information of devices.
func (h *HUAWEI) ATI() (devices.Product, error) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()
	res, err := h.Cmd(ATI)
	var p devices.Product
	if err != nil {
		return p, err
	}
	s := strings.Split(res, "\r")
	Manufacturer := "Manufacturer"
	Revision := "Revision"
	GCAP := "GCAP"
	for i := range s {
		switch true {
		case strings.Contains(s[i], Manufacturer):
			p.Manufacturer = strings.Trim(s[i], Manufacturer+":")
		case strings.Contains(s[i], "Model"):
			p.Model = strings.Trim(s[i], "Model:")
		case strings.Contains(s[i], Revision):
			p.Revision = strings.Trim(s[i], Revision+":")
		case strings.Contains(s[i], IMEI):
			p.IMEI = strings.Trim(s[i], IMEI+":")
		case strings.Contains(s[i], GCAP):
			p.GCAP = strings.Trim(s[i], "+"+GCAP+":")
		default:
		}
	}
	return p, nil
}

//IMEI Get IMEI
func (h *HUAWEI) IMEI() (str string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			return
		}
	}()
	str, err = h.Cmd(IMEI)
	if err != nil {
		panic(err)
	}
	s := strings.Split(str, "\r")
	if len(s) < 5 {
		panic(fmt.Errorf("Error result: %s ", str))
	}
	return s[1], nil
}

type HCSQs struct {
	Data []string `json:"-"`
	*devices.HCSQs
}

func (h *HCSQs) Serialization() ([]byte, error) {
	switch h.Sysmode {
	case LTE:
		var l Lte
		l.HCSQs = h
		l.RSRP()
		l.RSRQ()
		l.SINR()
		l.RSSI()
	case NR:
		var n Nr
		n.HCSQs = h
		n.SINR()
		n.RSRQ()
		n.RSRP()
	case WCDMA:
		var w Wcdma
		w.HCSQs = h
		w.RSSI()
		w.Rscp = w.RSCP(MustInt(w.Data[1]))
		w.Ecio = w.ECIO(MustInt(w.Data[2]))
	case GSM:
		var w Wcdma
		w.HCSQs = h
		w.RSSI()
	default:
	}
	return json.Marshal(h)
}

//HCSQ ask what type of signal is using , and it's strength
func (h *HUAWEI) HCSQ() ([]HCSQs, error) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()
	str, err := h.Cmd(HCSQ)
	var arr []HCSQs
	if err != nil {
		return arr, err
	}
	s := strings.Split(str, "\r")
	if len(s) < 3 {
		return arr, fmt.Errorf("Error result: %s ", str)
	}

	for i := range s {
		result := s[i]
		result = RemoveCrCL(result)
		if result == HCSQ {
			continue
		}
		if result == "OK" {
			continue
		}
		if len(result) > 0 {
			var h HCSQs
			ar := Result2Array(result, "^HCSQ: ", ",")
			h.Sysmode = RemoveQuote(ar[0])
			h.Data = ar[1:]
			arr = append(arr, h)
		}
	}

	return arr, nil
}

//EONS get operate network name and sim efspn information
func (h *HUAWEI) EONS() (str string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			return
		}
	}()
	str, err = h.Cmd(EONS + "?1")
	if err != nil {
		panic(err)
	}
	s := strings.Split(str, "\r")
	if len(s) < 5 {
		panic(fmt.Errorf("Error result: %s ", RemoveCrCL(str)))
	}
	result := s[1]
	result = RemoveCrCL(result)
	return result, nil
}

//SYSINFO show system info
func (h *HUAWEI) SYSINFO() (d devices.SysInfoEx, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			return
		}
	}()
	info := &devices.SysInfoEx{}
	if str, err := h.Cmd(SYSINFOEX); err != nil {
		panic(err)
	} else {
		if s := strings.Split(str, "\r"); len(s) > 3 {
			s[1] = RemoveCrCL(s[1])
			sys := "^SYSINFOEX:"
			if strings.HasPrefix(s[1], sys) {
				str = strings.Trim(s[1], sys)
				s := strings.Split(str, ",")
				for i := range s {
					switch i {
					case 0:
						info.SimState, _ = strconv.Atoi(s[i])
					case 1:
						info.SrvDomain, _ = strconv.Atoi(s[i])
					case 2:
						info.RoamStatus, _ = strconv.Atoi(s[i])
					case 3:
						info.SimState, _ = strconv.Atoi(s[i])
					case 4:
						//LockState now support by hardware now.
						if s[i] == "" {
							info.LockState = 0
							continue
						}
						info.LockState, _ = strconv.Atoi(s[i])
					case 5:
						info.SysMode, _ = strconv.Atoi(s[i])
					case 6:
						info.SysModeName = RemoveQuote(s[i])
					case 7:
						info.SubMode, _ = strconv.Atoi(s[i])
					case 8:
						info.SubModeName = RemoveQuote(s[i])
					}
				}
				return *info, nil
			}
		}
		return *info, nil
	}
}

type LenDc struct {
	//Enable
	Enable bool `json:"enable"`
	//Is it support endc(NR)?
	Available bool `json:"endc_available"`
	// false : not support , true , support.
	PlmnAvailable bool `json:"endc_plmn_available"`
	// false mean restricted , true means note restricted
	Restricted bool `json:"endc_restricted"`
	//Is is NR ENDC now
	Pscell bool `json:"nr_pscell"`
}

//LENDC: Got lendc result
func (h *HUAWEI) LENDC() (L LenDc, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			return
		}
	}()
	if result, err := h.Cmd(LENDC); err == nil {
		s := strings.Split(result, "\r")
		if len(s) < 2 {
			panic(fmt.Errorf("got some error: %s", s))
		}
		s[1] = RemoveCrCL(s[1])
		str := strings.TrimPrefix(s[1], "^LENDC: ")
		s = strings.Split(str, ",")
		var result LenDc
		for i, v := range s {
			var ok bool
			switch v {
			case "1":
				ok = true
			case "0":
				ok = false
			}
			switch i {
			case 0:
				result.Enable = ok
			case 1:
				result.Available = ok
			case 2:
				result.PlmnAvailable = ok
			case 3:
				result.Restricted = ok
			case 4:
				result.Pscell = ok
			}
		}
		//fmt.Println(result)
		return result, nil
	}
	return LenDc{}, nil
}

//Only use it when it using 5g status
type NrRegStatus struct {
	ReportStatus int `json:"n"`
	Stat         int `json:"stat"`
}

//C5GREG: AT+C5GREG Only use it when it using 5g networking.
func (h *HUAWEI) C5GREG() (n NrRegStatus, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			return
		}
	}()
	if res, err := h.Cmd(C5GREG + "?"); err == nil {
		s := strings.Split(res, "\r")
		if len(s) < 3 {
			panic(err)
		}
		s[1] = RemoveCrCL(s[1])
		s = Result2Array(s[1], "+C5GREG: ", ",")
		if len(s) < 2 {
			panic(fmt.Errorf("too few result to decode. "))
		}
		var nrs NrRegStatus
		nrs.ReportStatus, err = strconv.Atoi(s[0])
		if nrs.ReportStatus == 0 || err != nil {
			panic(err)
		}
		nrs.Stat, _ = strconv.Atoi(s[1])
		return nrs, nil
	} else {
		panic(err)
	}
}

//Plmn: Public Land Mobile Network
type Plmn struct {
	Mode      int    `json:"mode"`
	Format    int    `json:"format"`
	Operator  string `json:"oper"`
	RadioType int    `json:"rat"`
}

//GetMode human readable format
func (p *Plmn) GetMode() string {
	var opt = []string{
		"Auto mode",
		"manual mode",
		"None-reg networking",
		"only return format when +COPS?",
		"if manual mode fail , change to auto mode",
	}
	return opt[p.Mode]
}

//GetFormat human readable format
func (p *Plmn) GetFormat() string {
	var opt = []string{
		"long char operator information",
		"short char operator information",
		"number format operator information",
	}

	return opt[p.Format]
}

//GetRadioType human readable format
func (p *Plmn) GetRadioType() string {
	var opt = map[int]string{
		0:  GSM,
		2:  WCDMA,
		7:  LTE,
		12: NR,
	}
	if res, ok := opt[p.RadioType]; !ok {
		return "Unknown Network mode"
	} else {
		return res
	}
}

//Cops string
func (h *HUAWEI) COPS() (result Plmn, err error) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			err = r.(error)
			return
		}
	}()
	str, err := h.Cmd(COPS)
	if err != nil {
		panic(err)
	}
	s := strings.Split(str, "\r")
	if len(s) < 2 {
		panic(fmt.Errorf("too few result to decode. "))
	}
	s = Result2Array(s[1], "+COPS: ", ",")
	if len(s) == 4 {
		for i := range s {
			switch i {
			case 0:
				result.Mode, err = strconv.Atoi(s[i])
			case 1:
				result.Format, err = strconv.Atoi(s[i])
			case 2:
				result.Operator = RemoveQuote(s[i])
			case 3:
				result.RadioType, err = strconv.Atoi(s[i])
			}
		}
	}

	if err != nil {
		panic(err)
	}
	return
}

//CGPADDR
func (h *HUAWEI) CGPADDR() (ip net.IP, err error) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			err = r.(error)
			return
		}
	}()
	if str, err := h.Cmd(CGPADDR); err != nil {
		panic(err)
	} else {
		s := strings.Split(str, "\r")
		if len(s) == 2 {
			panic(fmt.Errorf("No Dail up yet. use AT^NDISDUP=1,1 to Dail up with operator manually. "))
		}
		if len(s) < 3 {
			panic(fmt.Errorf("too few result to decode. "))
		}
		s = Result2Array(s[1], "+CGPADDR: ", ",")
		ip = net.IP{}
		err = ip.UnmarshalText([]byte(RemoveQuote(s[1])))
		if err != nil {
			panic(err)
		}
	}
	return
}

//NDISDUP Dial up with config.
func (h *HUAWEI) NDISDUP(config int, on int) (ok bool, err error) {
	if on != 0 {
		on = 1
	}
	if _, err = h.Cmd(fmt.Sprintf(NDISDUP+"=,%x,%x", config, on)); err != nil {
		ok = false
	} else {
		ok = true
	}
	return
}

//CPIN  just check if it ready.
func (h *HUAWEI) CPIN() (ok bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
			err = r.(error)
			return
		}
	}()
	var str string
	if str, err = h.Cmd(CPIN + "?"); err != nil {
		panic(err)
	} else {
		ok = strings.Contains(str, "READY")
		if ok {
			return
		} else {
			return ok, fmt.Errorf("Got result: %s ", str)
		}
	}
}
