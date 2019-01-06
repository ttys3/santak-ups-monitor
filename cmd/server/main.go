package main

import (
	"bitbucket.org/8ox86/santak-monitor/pkg/santak"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/multi"
	"os"
	"time"
	"unsafe"

	"github.com/tarm/serial"
)

func init() {
	initLog()
}

func main() {
	c := &serial.Config{Name: "COM4", Baud: 2400, ReadTimeout: time.Second * 10 }
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err.Error())
	}

	n, err := s.Write([]byte("Q1\r"))
	if err != nil {
		log.Fatal(err.Error())
	}

	result := make([]byte, 0)
	buf := make([]byte, 128)

	for {
		n, err = s.Read(buf[0:])
		if err != nil {
			log.Fatal(err.Error())
		}
		//log.Printf("%#v", buf[:n])
		if string(buf[0:n]) == "\r" {
			break
		}
		//log.Printf("%q", buf[:n])
		result = append(result, buf[:n]...)
	}


	//byte array to struct
	rs := *(**santak.QueryResult)(unsafe.Pointer(&result))

	//log.Infof("%q", result)
	log.Infof(" %s\n", " -- HuangYeWuDeng SanTak UPS Monitor -- ")
	log.Infof("输入侧电压(初级电压): %s v\n", rs.IPVoltage)
	log.Infof("输入侧故障电压: %s v\n", rs.IPFaultVoltage)
	log.Infof("输出侧电压(次级电压): %s v\n", rs.OPVoltage)
	log.Infof("输出侧负载: %s%%\n", rs.OPLoad)
	log.Infof("输入侧频率: %s Hz\n", rs.IPFreq)
	log.Infof("电池电压: %s v\n", rs.BatteryVoltage)
	log.Infof("温度: %s °C\n", rs.Temperature)
	//log.Printf("状态: %#v\n", rs.Status)
	if rs.Status.UtilityFail == '1' {
		log.Errorf("市电状态: 已断电\n")
	} else {
		log.Infof("市电状态: 正常供电\n")
	}
}

func initLog() {
	//default handler: text
	var logHandlers []log.Handler
	//only print to console when debug
	logHandlers = append(logHandlers, cli.New(os.Stdout))
	log.SetHandler(multi.New(logHandlers...))
}