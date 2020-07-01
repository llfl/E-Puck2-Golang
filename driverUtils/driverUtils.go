package driverutils

import (
	"golang.org/x/exp/io/i2c"
	// "bytes"
	// "binary"
	"fmt"
	"time"
	"github.com/imroc/biu"
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

//SendCmd send command to epuck
func (e *EPuckHandle) SendCmd(Actuator []uint8) bool {
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
		for i := 0; i < ActuatorSize; i++ {
			fmt.Println("[",i,"]",biu.ToBinaryString(Actuator[i]))
		}
		return true
	}
	return false
}

//Stop epuck stop
func (e *EPuckHandle) Stop() bool {
	var Actuator = make([]uint8, ActuatorSize)
	return e.SendCmd(Actuator)
}

//FreeForward free forward
func (e *EPuckHandle) FreeForward(rspeed int, lspeed int) bool {
	var Actuator = make([]uint8, ActuatorSize)
	RSL := uint8(rspeed)
	RSH := uint8(rspeed>>8)
	LSL := uint8(lspeed)
	LSH := uint8(lspeed>>8)
	Actuator[0] = RSL
	Actuator[1] = RSH
	Actuator[2] = LSL
	Actuator[3] = LSH
	return e.SendCmd(Actuator)
}

//Forward go forward 
func (e *EPuckHandle) Forward(speed int) bool {
	return e.FreeForward(speed,speed)
}

//FreeSpin freespin
func (e *EPuckHandle) FreeSpin(flag bool) bool {
	var rspeed,lspeed int
	if flag {
		rspeed = 64
		lspeed = -64
	}else{
		rspeed = -64
		lspeed = 64
	}
	return e.FreeForward(rspeed,lspeed)
}

//Spin epuck spin around
func (e *EPuckHandle) Spin(degree float32) bool {
	var f bool
	if degree < 0 {
		f = LEFT
	}else{
		f = RIGHT
	}
	if e.FreeSpin(f){
		t := time.Duration(degree * RATIO * 1000)
		time.Sleep(t * time.Millisecond)
		return e.Stop()
	}
	return false
	
}
