package santak

const StartByteChar = '('
const EndByteChar = '\r'

// 无效命令和信息的处理:
// 	收到无效的命令时,UPS 要将受到的内容原样返回。若命令 UPS 无法返回信息,则返回“@”

// UPS 开关量状态:<U>, <U>是以二进制数位表示法:<b7b6b5b4b3b2b1b0>, 并以 ASCII 码单位传输的一个状态量
// b7:1 表示 市电电压异常
// b6:1 表示 电池低电压
// b5:1 表示 Bypass 或 Buck Active
// b4:1 表示 UPS 故障
// b3:1 表示 UPS 为后备式(0 表示在线式)
// b2:1 表示 测试中
// b1:1 表示 关机有效
// b0:1 表示 蜂鸣器开
type UPSStatus struct {
	UtilityFail       byte //市电电压异常, utility power 市电，1 fail 表示 UPS 市电断了，电池供电
	BatteryLow        byte //电池低电压, 1 low
	BypassBoostActive byte //Bypass 或 Buck Active (1 : AVR 0:NORMAL)
	//当电网电压变化时，通过自动调整实现输出电压稳定并供给负载使用，称为AVR(AutomaticVoltageRegulation)
	UPSFailed      byte //UPS 故障, 1 failed
	UPSType        byte //表示 UPS 为后备式(0 表示在线式), 1 standby 0 online
	TestActive     byte //表示 测试中, 1 test in progress
	ShutdownActive byte //表示 关机有效, 1 shutdown active
	BuzzerActive   byte //表示 蜂鸣器开
}

// UPS 状态查询请求 Q1
//"(228.0 228.0 228.4 006 50.2 27.4 25.0 00001000"
//(228.0 228.0 228.4 017 50.0 27.4 25.0 00001001
type QueryResult struct {
	StartByte      byte
	IPVoltage      [5]byte //输入电压(I/P voltage):MMM.M, M 为0~9的整数，状态量单位为 Vac
	_              byte
	IPFaultVoltage [5]byte //输入故障电压(I/P fault voltage):NNN.N, N 为 0~9 的整数,状态量单位为 Vac
	// ** 对后备式 UPS 而言 **
	// 目的是为了标识引起后备式 UPS 转入逆变模式的瞬间毛刺电压。如有电压
	// 瞬变发生,输入电压将在电压瞬变前、后一个查询保持正常。 I/P 异常电压将把瞬
	// 变电压保持到下一个查询。查询完成后,I/P 异常电压将与 I/P 电压保持一致,直
	// 到发生新的瞬变。
	// ** 对在线式 UPS 而言 **
	// 目的是为了标识引起在线式 UPS 转入电池供电模式的短时输入异常。如有
	// 电压瞬变发生,输入电压将在电压瞬变前、后一个查询保持正常。 I/P 异常电压将
	// 把瞬变电压保持到下一个查询。查询完成后,I/P 异常电压将与 I/P 电压保持一致
	// 直到发生新的瞬变。
	_                byte
	OPVoltage        [5]byte //输出电压(O/P voltage):PPP.P, P 为 0~9 的整数,状态量单位为 Vac
	_                byte
	OPCurrentPercent [3]byte //输出电流(O/P current):QQQ, QQQ 是一个相对于最大允许电流的百分比,不是一个绝对值
	_                byte
	IPFreq           [4]byte //输入频率(I/P frequency):RR.R, R 为 0~9 的整数,状态量单位为 Hz
	_                byte
	BatteryVoltage   [4]byte //电池电压(Battery voltage):SS.S 或 S.SS, S 为 0~9 的整数
	// 对在线式单体电池电压显示方式为 S.SS Vdc
	// 对后备式总电池电压显示方式为 SS.S Vdc
	// ( UPS 类型将在 UPS 状态信息中获得 )
	_           byte
	Temperature [4]byte //环境温度(Temperature):TT.T, T 为 0~9 的整数,单位为 C
	_           byte
	Status      UPSStatus //UPS 开关量状态:<U>, <U>是以二进制数位表示法:<b7b6b5b4b3b2b1b0>,
	// 并以 ASCII 码单位传输的一个状态量
	// b7:1 表示 市电电压异常
	// b6:1 表示 电池低电压
	// b5:1 表示 Bypass 或 Buck Active
	// b4:1 表示 UPS 故障
	// b3:1 表示 UPS 为后备式(0 表示在线式)
	// b2:1 表示 测试中
	// b1:1 表示 关机有效
	// b0:1 表示 蜂鸣器开
	StopByte byte
}

//UPS 额定值信息
//这个功能是使 UPS 能回答额定值信息。每个信息段的之间有一个空格符。
//输入：F<CR>
//输出：#MMM.M QQQ SS.SS RR.R<CR>
//     #220.0 007 24.00 50.0
// 信息段格式定义如下:
// 额定电压:MMM.M
// 额定电流:QQQ
// 电池电压:SS.SS 或 SSS.S
// 额定频率:RR.R
type RatingInfo struct {
	_               byte
	VoltageRating   [5]byte
	_               byte
	CurrentRating   [3]byte
	_               byte
	BatteryVoltage  [5]byte
	_               byte
	FrequencyRating [4]byte
	StopByte        byte
}

// func (r RatingInfo) String() string {
// 	return fmt.Sprintf("VoltageRating: %s, CurrentRating: %s, BatteryVoltage: %s, FrequencyRating: %s",
// 		r.VoltageRating, r.CurrentRating, r.BatteryVoltage, r.FrequencyRating)
// }
