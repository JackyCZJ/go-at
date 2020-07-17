package mh5000_31

import (
	"errors"
	"net"

	goat "github.com/JackyCZJ/go-at"
)

type Modem struct {
	devices *HUAWEI
}

func GetModem(huawei *HUAWEI) goat.Modem {
	return &Modem{
		devices: huawei,
	}
}

//Get Devices info
func (m *Modem) DeviceInfo() (goat.Product, error) {
	return m.devices.ATI()
}

//Query Signal strength
func (m *Modem) SignalStrength() (data []goat.HCSQs, err error) {
	h, err := m.devices.HCSQ()
	if err != nil {
		return nil, err
	}
	for i := range h {
		data = append(data, *h[i].HCSQs)
	}
	return
}

//What Network provider now?
func (m *Modem) CurrentNetwork() (string, error) {
	if r, err := m.devices.COPS(); err != nil {
		return "", err
	} else {
		return r.Operator, nil
	}
}

//is it register to operator
func (m *Modem) IsRegistered() bool {
	if _, err := m.devices.COPS(); err != nil {
		return false
	}
	return true
}

//iS NSA or SA?
func (m *Modem) IsNSA() bool {
	if l, err := m.devices.LENDC(); err != nil {
		return false
	} else {
		return l.Pscell
	}
}

//Ipv4 or ipv6 and it's addr
func (m *Modem) PDPAddr() (net.IP, error) {
	return m.devices.CGPADDR()
}

//Dial operator
func (m *Modem) Dial() error {
	if ok, err := m.devices.NDISDUP(1, 1); err != nil {
		return err
	} else {
		if ok {
			return nil
		}
		return errors.New("Dial up failed. ")
	}
}

/*Pin management*/
//PinLockStatus Query pin lock status
func (m *Modem) PinLockStatus() bool {
	if ok, err := m.devices.CPIN(); err != nil {
		return false
	} else {
		return ok
	}
}
