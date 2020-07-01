package driverutils

import (
	"golang.org/x/exp/io/i2c"
	// "bytes"
	// "binary"
	"fmt"
	"time"
)

const (
	// DefaultI2CDevices I2C设备
	DefaultI2CDevices = "/dev/i2c-4"

	// DefaultI2CAddress I2C设备的地址
	DefaultI2CAddress = 0x1F

	// ActuatorSize 定义控制器的长度 with checksum
	ActuatorSize = (19 + 1)

	// SensorDataSize 定义传感器数据的长度 with checksum
	// SensorDataSize = (46 + 1)

	// LEFT flag bit
	LEFT = false
	// RIGHT flag bit
	RIGHT = true
	//RATIO the ratio of spin time to degree
	RATIO = 0.0120133052
)

// var (
// 	//Actuator epuck actuator
// 	// Actuator = make([]uint8, ActuatorSize)

// 	//SensorData epuck sensor result
// 	SensorData = make([]uint8, SensorDataSize)
// )

// EPuckHandle 为EPuck的操作接口
type EPuckHandle struct{
	I2CDevices  string
	I2CAddress  int
	Device  *i2c.Device
}

// NewEPuckHandle new a handle for epuck
func NewEPuckHandle(opts ...Option) *EPuckHandle {
	options := newOptions(opts...)
	device, err := i2c.Open(&i2c.Devfs{Dev: options.I2CDevices}, options.I2CAddress)
	if err != nil {
		return nil	}
	return &EPuckHandle{
		I2CDevices:options.I2CDevices,
		I2CAddress:options.I2CAddress,
		Device:device,
	}
}

func (e *EPuckHandle) sendCmd(Actuator []uint8) bool {
	var checksum uint8
	for i := 0; i < (ActuatorSize - 1); i++ {
		checksum ^= Actuator[i]
	}
	Actuator[ActuatorSize-1] = checksum
	for trails := 0; trails < 3; trails++ {
		err := e.Device.Write(Actuator)
		if err != nil {
			continue
		}
		fmt.Println("write success!")
		return true
	}
	return false
}

func (e *EPuckHandle) forword(speed int) bool {
	SL := uint8(speed)
	SH := uint8(speed>>8)
	var Actuator = make([]uint8, ActuatorSize)
	Actuator[0] = SL
	Actuator[1] = SH
	Actuator[2] = SL
	Actuator[3] = SH
	return e.sendCmd(Actuator)
}

func (e *EPuckHandle) stop() bool {
	var Actuator = make([]uint8, ActuatorSize)
	return e.sendCmd(Actuator)
}

func (e *EPuckHandle) spin(degree float32) bool {
	var f bool
	if degree < 0 {
		f = LEFT
	}else{
		f = RIGHT
	}
	if e.freespin(f){
		t := time.Duration(degree * RATIO * 1000)
		time.Sleep(t * time.Millisecond)
		return e.stop()
	}
	return false
	
}

func (e *EPuckHandle) freespin(flag bool) bool {
	var Actuator = make([]uint8, ActuatorSize)
	if flag {
		Actuator[0] = 0xC0
		Actuator[1] = 0xFF
		Actuator[2] = 0x04
		Actuator[3] = 0x00
	}else{
		Actuator[0] = 0x04
		Actuator[1] = 0x00
		Actuator[2] = 0xC0
		Actuator[3] = 0xFF
	}
	return e.sendCmd(Actuator)
}

