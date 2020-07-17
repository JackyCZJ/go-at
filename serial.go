package go_at

import (
	"fmt"
	"sync"
	"time"

	"github.com/albenik/go-serial/v2"
)

const BaudRate = 115200
const defaultTimeout = 30 * time.Second

type Device struct {
	Name     string
	Port     *serial.Port
	PortName string
	Result   func(str string) (bool, error)
	Opt      []serial.Option
	*sync.Mutex
}

//Connect : Connect to devices with configure.
func (d *Device) Connect() (*serial.Port, error) {
	port, err := serial.Open(d.PortName, d.Opt...)
	if err != nil {
		if ports, err := serial.GetPortsList(); err != nil || len(ports) == 0 {
			return nil, err
		} else {
			fmt.Println("Invali port name, which port you want to connect with?")
			for i := range ports {
				fmt.Println(i+1, ":", ports[i])
			}
		}
		return nil, err
	}
	return port, nil
}

//SetPortName ,default will be /dev/ttyUSB1
func (d *Device) SetPortName(portName string) {
	if portName == "" {
		d.PortName = "/dev/ttyUSB1"
	}
	d.PortName = portName
}

//Cmd send command to devices and return result when result contains ok.
func (d *Device) Cmd(command string) (result string, err error) {
	d.Lock()
	defer d.Unlock()
	if d.PortName == "" {
		return result, fmt.Errorf("No port name are set. ")
	}
	if d.Port, err = d.Connect(); err != nil {
		return result, err
	}
	defer func() {
		if err := d.Port.Close(); err != nil {
			fmt.Println("Close port error: ", err.Error())
			return
		}
	}()
	//Inspired by github.com/xlab/at
	t := time.NewTimer(defaultTimeout)
	defer t.Stop()
	stop := make(chan struct{}, 1)
	defer close(stop)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
		select {
		case <-t.C:
			if err := d.Port.Close(); err != nil {
				panic(err)
			}
		case <-stop:
			return
		}
	}()
	if err := d.Port.ResetOutputBuffer(); err != nil {
		return "", err
	}
	_, err = d.Port.Write([]byte(command + "\r\n"))
	if err != nil {
		return result, err
	}
	buff := make([]byte, 1024)
	t.Reset(defaultTimeout)
	for {
		n, err := d.Port.Read(buff)
		if err != nil {
			return "", err
		}
		result = string(buff[:n])
		if ok, err := d.Result(result); ok {
			return result, nil
		} else {
			if err != nil {
				return "", err
			}
			continue
		}
	}
}
