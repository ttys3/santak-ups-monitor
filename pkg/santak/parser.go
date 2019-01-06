package santak

const StartByteChar = '('
const EndByteChar = '\r'

type UPSStatus struct {
	UtilityFail byte //utility power 市电，1 fail 表示 UPS 市电断了，电池供电
	BatteryLow byte //1 low
	BypassBoostActive byte
	UPSFailed byte //1 failed
	UPSType byte //1 standby 0 online
	TestActive byte //1 test in progress
	ShutdownActive byte //1 shutdown active
	Reserved byte //always 0
}

//"(228.0 228.0 228.4 006 50.2 27.4 25.0 00001000"
type QueryResult struct {
	StartByte byte
	IPVoltage [5]byte
	_ byte
	IPFaultVoltage [5]byte
	_ byte
	OPVoltage [5]byte
	_ byte
	OPLoad [3]byte
	_ byte
	IPFreq [4]byte //Hz
	_ byte
	BatteryVoltage [4]byte
	_ byte
	Temperature [4]byte
	_ byte
	Status UPSStatus
	StopByte byte
}