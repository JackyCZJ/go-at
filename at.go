package go_at

import (
	"net"
)

type Product struct {
	Manufacturer string
	Model        string
	Revision     string
	IMEI         string
	GCAP         string `json:"+GCAP:"`
}

type SysInfoEx struct {
	SrvStatus  int `json:"srv_status"`
	SrvDomain  int `json:"srv_domain"`
	RoamStatus int `json:"roam_status"`
	SimState   int `json:"sim_state"`
	//Lock state unsupported now
	LockState   int    `json:"lock_state,omitempty"`
	SysMode     int    `json:"sysmode"`
	SysModeName string `json:"sysmode_name"`
	SubMode     int    `json:"submode"`
	SubModeName string `json:"submode_name"`
}

type HCSQs struct {
	Sysmode string
	Rssi    string `json:"rssi,omitempty"`
	Rsrp    string `json:"rsrp,omitempty"`
	Ecio    string `json:"ecio,omitempty"`
	Sinr    string `json:"sinr,omitempty"`
	Rsrq    string `json:"rsrq,omitempty"`
	Rscp    string `json:"rscp,omitempty"`
}

//Todo: ⬇️ ⬇️ ⬇️
type Modem interface {
	//Get Devices info
	DeviceInfo() (Product, error)
	//Query Signal strength
	SignalStrength() ([]HCSQs, error)
	//What Network provider now?
	CurrentNetwork() (string, error)
	//is it register to operator
	IsRegistered() bool
	//iS NSA or SA?
	IsNSA() bool
	//Ipv4 or ipv6 and it's addr
	PDPAddr() (net.IP, error)
	//Dial operator
	Dial() error

	/*Pin management*/
	//PinLockStatus Query pin lock status
	PinLockStatus() bool
}
