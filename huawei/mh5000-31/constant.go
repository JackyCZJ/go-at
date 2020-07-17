package mh5000_31

const (
	//IMEI IMEI :)
	IMEI = "AT+GSN"
	//ATI ask minor information of devices.
	ATI = "ATI"
	//HCSQ ask what type of signal is using , and  it's strength
	HCSQ = "AT^HCSQ?"
	//EONS Query network operator name and sim card EFSPN info
	EONS = "AT^EONS"
	//SYSINFOEX ask system(Network) Info
	SYSINFOEX = "AT^SYSINFOEX"
	//COPS Query is it register to operator.
	COPS = "AT+COPS?"
	//CPIN Query and management Pin code.
	CPIN = "AT+CPIN"
	//LENDC? Query it is NSA Signal
	LENDC = "AT^LENDC?"
	//C5GREG  5g register status
	C5GREG = "AT+C5GREG"
	//CGPADDR Get addr after dial
	CGPADDR = "AT+CGPADDR"

	NDISDUP = "AT^NDISDUP"
)

//Network mode status
const (
	//No service
	NoService = "NOSERVICE"
	//GSM  GSM/GRPS/EDGE mode
	GSM = "GSM"
	//WCDMA	WCDMA/HSDPA/HSPA mode
	WCDMA = "WCDMA"
	//LTE <-mode
	LTE = "LTE"
	//NR : 5G!!!!!
	NR = "NR"
)
