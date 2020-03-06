package main

import (
	"io"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"bitbucket.org/8ox86/santak-monitor/pkg/santak"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/multi"
	"github.com/tarm/serial"
)

func init() {
	initLog()
}

func main() {
	var ttyDev = "/dev/ttyUSB0"
	if runtime.GOOS == "windows" {
		ttyDev = "COM4"
	}
	c := &serial.Config{Name: ttyDev, Baud: 2400, ReadTimeout: time.Millisecond * 500}

	log.Infof("try open port: %s\n", ttyDev)

	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err.Error())
	}

	result := make([]byte, 0)
	buf := make([]byte, 128)

	var wg sync.WaitGroup

	// After setting everything up!
	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		log.Info("Received an interrupt, stopping services...")
		log.Info("waiting exit ...")
		wg.Wait()
		log.Info("exited")
		s.Flush()
		s.Close()
		close(cleanupDone)
	}()

	showRatedInfo(s)
	// log.Info("try send test command ...")
	// //测试 10 秒钟后返回市电供电
	// _, err = s.Write([]byte("T\r"))
	// if err != nil {
	// 	log.Errorf("send err: %s", err.Error())
	// }

exit:
	for {
		//clear the slice first
		result = result[:0]

		select {
		case <-cleanupDone:
			break exit
		default:
			wg.Add(1)
			log.Info("sleep 3 start ...")
			time.Sleep(time.Second * 3)
			log.Info("sleep 3 end ...")
			log.Info("try send query command ...")
			n, err := s.Write([]byte("Q1\r"))
			if err != nil {
				log.Errorf("send err: %s", err.Error())
				wg.Done()
				continue
			}

			log.Info("try read info ...")
			for {
				n, err = s.Read(buf[0:])
				if err != nil {
					if err != io.EOF {
						log.Errorf("read err: %s", err.Error())
					}
					break
				} else {
					log.Infof("read data: %#v", string(buf[:n]))
				}
				if string(buf[0:n]) == "\r" {
					log.Info("hit cr ...")
					break
				}
				result = append(result, buf[:n]...)
			}
			wg.Done()
			log.Infof("done read info: %s", string(result))
		}

		// rt := *(**santak.RatingInfo)(unsafe.Pointer(&result))
		// log.Infof("RatingInfo: %#v\n", rt)

		//byte array to struct
		rs := *(**santak.QueryResult)(unsafe.Pointer(&result))

		//log.Infof("%q", result)
		log.Infof(" %s\n", " -- HuangYeWuDeng SanTak UPS Monitor -- ")
		log.Infof("输入侧电压(初级电压): %s v\n", rs.IPVoltage)
		log.Infof("输入侧故障电压: %s v\n", rs.IPFaultVoltage)
		log.Infof("输入侧频率: %s Hz\n", rs.IPFreq)

		log.Infof("输出侧电压(次级电压): %s v\n", rs.OPVoltage)
		log.Infof("输出侧电流负载: %s%%\n", rs.OPCurrentPercent)

		log.Infof("电池电压: %s v\n", rs.BatteryVoltage)
		log.Infof("温度: %s °C\n", rs.Temperature)
		log.Infof("buzzerActive: %#v\n", rs.Status.BuzzerActive)

		log.Infof("状态: %#v\n", rs.Status)

		//断电状态：
		//santak.UPSStatus{UtilityFail:0x31, BatteryLow:0x30, BypassBoostActive:0x30, UPSFailed:0x30, UPSType:0x31, TestActive:0x30, ShutdownActive:0x30, Reserved:0x31}
		if rs.Status.UtilityFail == '1' {
			log.Errorf("市电状态: 已断电\n")
		} else {
			log.Infof("市电状态: 正常供电\n")
		}

		if rs.Status.BatteryLow == '1' {
			log.Errorf("电池电压: 低\n")
		} else {
			log.Infof("电池电压: 正常\n")
		}
	} //end for
}

func initLog() {
	//default handler: text
	var logHandlers []log.Handler
	//only print to console when debug
	logHandlers = append(logHandlers, cli.New(os.Stdout))
	log.SetHandler(multi.New(logHandlers...))
}

func showRatedInfo(s *serial.Port) {
	log.Info("showRatedInfo begin ...")
	//测试 10 秒钟后返回市电供电
	_, err := s.Write([]byte("F\r"))
	if err != nil {
		log.Errorf("showRatedInfo send err: %s", err.Error())
	}
	log.Info("showRatedInfo try read info ...")
	result := make([]byte, 0)
	buf := make([]byte, 128)
	for {
		n, err := s.Read(buf[0:])
		if err != nil {
			if err != io.EOF {
				log.Errorf("read err: %s", err.Error())
			}
			break
		} else {
			// log.Infof("read data: %#v", string(buf[:n]))
		}
		if string(buf[0:n]) == "\r" {
			log.Info("hit cr ...")
			break
		}
		result = append(result, buf[:n]...)
	}
	rt := *(**santak.RatingInfo)(unsafe.Pointer(&result))
	log.Infof("BatteryVoltage: %s\n", rt.BatteryVoltage)
	log.Infof("VoltageRating: %s\n", rt.VoltageRating)
	log.Infof("CurrentRating: %s\n", rt.CurrentRating)
	log.Infof("FrequencyRating: %s\n", rt.FrequencyRating)
}
