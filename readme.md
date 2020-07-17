GO-AT
-----

A way to send AT command or something else to modem.

Develop base huawei mh5000-31

### Usage
You can send some at commend that we packed.
 
 ```go
package main

import (
	"fmt"

	devices "git.esixcloud.net/jackyczj/go-at"
	"github.com/albenik/go-serial/v2"

	huawei "git.esixcloud.net/jackyczj/go-at/huawei/mh5000-31"
)

func main() {

	opt := []serial.Option{
		serial.WithDataBits(8),
		serial.WithStopBits(serial.OneStopBit),
		serial.WithBaudrate(devices.BaudRate),
		serial.WithParity(serial.OddParity),
	}
	h := new(huawei.HUAWEI)
	h.Init("/dev/ttyUSB1", opt...)
	res, err := h.HCSQ()
	if err != nil {
		fmt.Println("Got error result:", err)
		return
	}
	for i := range res {
		data, _ := res[i].Serialization()
		fmt.Println(string(data))
	}

}
```
or send some command.
```go
	opt := []serial.Option{
		serial.WithDataBits(8),
		serial.WithStopBits(serial.OneStopBit),
		serial.WithBaudrate(devices.BaudRate),
		serial.WithParity(serial.OddParity),
	}
	h := new(huawei.HUAWEI)
	h.Init("/dev/ttyUSB1", opt...)
	h.Cmd("...some command")
```