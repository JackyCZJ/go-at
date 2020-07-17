//utils.go: Convert to human readable format
package mh5000_31

import (
	"fmt"
	"strconv"
	"strings"
)

//Wcdma interface
type Wcdma struct {
	*HCSQs
}

//RSSI wcdma impalement
func (w *Wcdma) RSSI() {
	w.Rssi = rssi(MustInt(w.Data[0]))
}

//RSRP rsrp  for wcdma only
func (Wcdma) RSCP(i int) string {
	if i == 0 {
		return "rscp < -120 dBm"
	}
	if i > 96 || i == 255 {
		return "unknown or can't be measuring"
	}
	if i == 96 {
		return "-25 dBm ≤ rscp"
	}
	return fmt.Sprintf("%x dBm < rscp < %x dBm", -120+i-1, -120+i)
}

//ECIO ECIO for wcdma only
func (Wcdma) ECIO(i int) string {
	if i == 0 {
		return "Ec/Io < -32 dB "
	}
	if i > 65 || i == 255 {
		return "unknown or can't be measuring"
	}
	if i == 65 {
		return "0 dB ≤ Ec/Io"
	}
	return fmt.Sprintf("%x dB < Ec/Io < %x dB", float64(-120)+float64(i)*0.5-float64(1), -120+float64(i)*0.5)
}

//Lte interface
type Lte struct {
	*HCSQs
}

//RSRP Lte impalement
func (l *Lte) RSRP() {
	l.Rsrp = rsrp(MustInt(l.Data[1]))
}

//Rssi Lte impalement
func (l *Lte) RSSI() {
	l.Rssi = rssi(MustInt(l.Data[0]))
}

//Sinr Lte impalement
func (l *Lte) SINR() {
	l.Sinr = sinr(MustInt(l.Data[2]))
}

//Rsrq Lte impalement
func (l *Lte) RSRQ() {
	l.Rsrq = rsrq(MustInt(l.Data[3]))
}

type Nr struct {
	*HCSQs
}

//Rsrp Nr(5g) impalement
func (n *Nr) RSRP() {
	n.Rsrp = rsrp(MustInt(n.Data[0]))
}

//Sinr Nr(5g) impalement
func (n *Nr) SINR() {
	n.Sinr = sinr(MustInt(n.Data[1]))
}

//Rsrq Nr(5g) impalement
func (n Nr) RSRQ() {
	n.Rsrq = rsrq(MustInt(n.Data[2]))
}

//sinr only for LTE , Nr(5G)
func sinr(i int) string {
	if i == 0 {
		return "sinr < -19.5 dB"
	}
	if i > 251 || i == 255 {
		return "unknown or can't be measuring"
	}
	if i == 251 {
		return "30 dBm ≤ sinr\n"
	}
	return fmt.Sprintf("%x dB < sinr < %x dB", -19.5+(float64(i)-1)*0.5, -19.5+float64(i)*0.5)
}

//rsrq only for LTE , Nr(5G)
func rsrq(i int) string {
	if i == 0 {
		return "rsrq < -19.5 dB"
	}
	if i > 34 || i == 255 {
		return "unknown or can't be measuring"
	}
	if i == 34 {
		return "-3 dBm ≤ rsrq"
	}
	return fmt.Sprintf("%x dB < rsrq < %x dB", -19.5+(float64(i)-1)*0.5, -19.5+float64(i)*0.5)
}

//rsrp only for LTE NR(5G)
func rsrp(i int) string {
	if i == 0 {
		return "rscp < -140 dBm"
	}
	if i > 97 || i == 255 {
		return "unknown or can't be measuring"
	}
	if i == 97 {
		return "-44 dBm ≤ rscp"
	}
	return fmt.Sprintf("%x dBm < rscp < %x dBm", -140+i-1, -140+i)
}

//Rssi only for gsm,wcdma,lte
func rssi(i int) string {
	if i == 0 {
		return "rssi < -120 dBm"
	}
	if i > 96 || i == 255 {
		return "unknown or can't be measuring"
	}
	if i == 96 {
		return "-25 dBm ≤ rssi"
	}
	return fmt.Sprintf("%x dBm < rssi < %x dBm", -120+i-1, -120+i)
}

//RemoveCrCL remove CR AND CL
func RemoveCrCL(str string) string {
	str = strings.TrimPrefix(str, "\r")
	str = strings.TrimSuffix(str, "\r")
	str = strings.TrimPrefix(str, "\n")
	str = strings.TrimSuffix(str, "\n")
	return str
}

//RemoveQuote remove " prefix and suffix
func RemoveQuote(str string) string {
	str = strings.TrimPrefix(str, "\"")
	str = strings.TrimSuffix(str, "\"")
	return str
}

func Result2Array(str, prefix, sep string) []string {
	return strings.Split(strings.TrimPrefix(str, prefix), sep)
}

func MustInt(s string) (i int) {
	i, _ = strconv.Atoi(s)
	return
}
