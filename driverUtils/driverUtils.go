package driverutils

import (
	"golang.org/x/exp/io/i2c"
	// "bytes"
	// "binary"
	
	"time"
	"fmt"
	// "github.com/imroc/biu"
)

const (
	// DefaultI2CDevices I2C设备
	DefaultI2CDevices = "/dev/i2c-4"

	// DefaultI2CAddress I2C设备的地址
	DefaultI2CAddress = 0x1F

	// DefaultGyroAddress 陀螺仪设备的地址
	DefaultGyroAddress = 0x68

	// ActuatorSize 定义控制器的长度 with checksum
	ActuatorSize = (19 + 1)

	// SensorDataSize 定义传感器数据的长度 with checksum
	SensorDataSize = (46 + 1)

	// LEFT flag bit
	LEFT = false
	// RIGHT flag bit
	RIGHT = true
	//RATIO the ratio of spin time to degree
	RATIO = 12
	// RATIO = 0.0120133052

	//NumSamplesCalibrationCONST calibrating the gyro
	NumSamplesCalibrationCONST = 20

	//GyroDPSCoefficientCONST 陀螺仪的时间参数
	GyroDPSCoefficientCONST = 0.0076293945
)

// var (
// 	//Actuator epuck actuator
// 	// Actuator = make([]uint8, ActuatorSize)

// 	//SensorData epuck sensor result
// 	SensorData = make([]uint8, SensorDataSize)
// )


type ePuckSensors struct{
	prox      [8]uint16
	ambient   [8]uint16
	mic       [4]uint16
	sel       uint8
	button    uint8
	motorStep [2]int16
	tv        uint8
}

type ePuckGyro struct{
	Values [3]int16
	offsets [3]int16
	calibrated bool
}

// EPuckHandle 为EPuck的操作接口
type EPuckHandle struct{
	// I2CDevices  string
	// I2CAddress  int
	// GyroAddress int
	Device  *i2c.Device
	GyroDevice  *i2c.Device
	GyroEnabled  bool
	Sensors  ePuckSensors
	Gyro ePuckGyro
}

// NewEPuckHandle new a handle for epuck
func NewEPuckHandle(opts ...Option) *EPuckHandle {
	options := newOptions(opts...)
	device, err := i2c.Open(&i2c.Devfs{Dev: options.I2CDevices}, options.I2CAddress)
	if err != nil {
		fmt.Println(err)
		return nil	
	}
	gyro := &i2c.Device{}
	gyroEnabled := false
	if options.Gyro {
		gyro, _ = i2c.Open(&i2c.Devfs{Dev: options.I2CDevices}, options.GyroAddress)
		gyroEnabled = true
		fmt.Println("Gyro Enabled",gyro)
	}

	return &EPuckHandle{
		// I2CDevices:options.I2CDevices,
		// I2CAddress:options.I2CAddress,
		Device:device,
		GyroDevice:gyro,
		GyroEnabled:gyroEnabled,
		Sensors:ePuckSensors{},
		Gyro:ePuckGyro{},
	}
}

//UpdateSensors update the handle sensor structure
func (e *EPuckHandle) UpdateSensors() bool {
	var sensorData = make([]byte, SensorDataSize)
	if err := e.Device.Read(sensorData); err != nil {
		fmt.Println(err)
		return false
	}
	for i := 0; i < 8; i++ {
		e.Sensors.prox[i] = uint16(sensorData[i*2+1])*256 + uint16(sensorData[2*i])
		e.Sensors.ambient[i] = uint16(sensorData[16+i*2+1])*256 + uint16(sensorData[16+2*i])
	}
	for i := 0; i < 4; i++ {
		e.Sensors.mic[i] = uint16(sensorData[32+i*2+1])*256 + uint16(sensorData[32+2*i])
	}
	e.Sensors.sel = sensorData[40] & 0x0F
	e.Sensors.button = sensorData[40] >> 4
	for i := 0; i < 2; i++ {
		e.Sensors.motorStep[i] = int16(uint16(sensorData[41+i*2+1])*256 + uint16(sensorData[41+2*i]))
	}
	e.Sensors.tv = sensorData[45]
	if e.GyroEnabled {
		return e.UpdateGyro()
	}
	return true
}

//UpdateGyro update the handle sensor structure
func (e *EPuckHandle) UpdateGyro() bool {
	if !e.GyroEnabled {
		fmt.Println("GYRO NOT ENABLED: enable Gyro first!")
		return false
	}
	var gyroData = make([]byte, 6)
	if err := e.GyroDevice.ReadReg(0x43,gyroData); err != nil {
		fmt.Println(err)
		return false
	}
	// for i := 0; i < 6; i++ {
	// 	fmt.Println("[",i,"]",biu.ToBinaryString(gyroData[i]))
	// }
	for i := 0; i < 3; i++ {
		e.Gyro.Values[i] = int16(uint16(gyroData[i*2+1])+uint16(gyroData[2*i])*256) - e.Gyro.offsets[i]
	}
	return true
}

//CalibrateGyro CalibrateGyro
func (e *EPuckHandle) CalibrateGyro() bool {
	if !e.GyroEnabled {
		fmt.Println("GYRO NOT ENABLED: enable Gyro first!")
		return false
	}
	var gyroSum [3]int
	var gyroData = make([]byte, 6)
	time.Sleep(500 * time.Millisecond)
	for i := 0; i < NumSamplesCalibrationCONST; i++ {
		if err := e.GyroDevice.ReadReg(0x43,gyroData); err != nil {
			fmt.Println(err)
			return false
		}
		for i := 0; i < 3; i++ {
			gyroSum[i] = gyroSum[i] + int(int16(uint16(gyroData[i*2+1])+uint16(gyroData[2*i])*256))
		}
	}
	for i := 0; i < 3; i++ {
		e.Gyro.offsets[i] = int16(gyroSum[i] / NumSamplesCalibrationCONST)
	}
	e.Gyro.calibrated = true
	fmt.Println("Gyro Calibrated!")
	return true

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
			fmt.Println(err)
			continue
		}
		// for i := 0; i < ActuatorSize; i++ {
		// 	fmt.Println("[",i,"]",biu.ToBinaryString(Actuator[i]))
		// }
		return e.UpdateSensors()
	}
	return false
}

//Stop epuck stop
func (e *EPuckHandle) Stop() bool {
	var Actuator = make([]uint8, ActuatorSize)
	return e.SendCmd(Actuator)
}

//FreeForward free forward
func (e *EPuckHandle) FreeForward(lspeed int, rspeed int) bool {
	var Actuator = make([]uint8, ActuatorSize)
	RSL := uint8(rspeed)
	RSH := uint8(rspeed>>8)
	LSL := uint8(lspeed)
	LSH := uint8(lspeed>>8)
	Actuator[0] = LSL
	Actuator[1] = LSH
	Actuator[2] = RSL
	Actuator[3] = RSH
	return e.SendCmd(Actuator)
}

//Forward go forward 
func (e *EPuckHandle) Forward(speed int) bool {
	return e.FreeForward(speed,speed)
}

//FreeSpin freespin
func (e *EPuckHandle) FreeSpin(speed int) bool {

	return e.FreeForward(speed,-speed)
}

//Spin epuck spin around
func (e *EPuckHandle) Spin(degree float32) bool {
	if !e.GyroEnabled {
		fmt.Println("GYRO NOT ENABLED: use FreeSpin() instead or enable Gyro!")
		return false
	}
	speed := 512
	sign := 1
	var count float32
	if degree < 0 {
		sign = -1
	}
	if !e.Gyro.calibrated {
		e.CalibrateGyro()
	}
	finalPosition := degree / GyroDPSCoefficientCONST * 2
	fmt.Println("final position:",finalPosition)
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
	for {
		e.FreeSpin(sign * speed)
		<- t.C
		if e.UpdateGyro() {
			count += float32(e.Gyro.Values[2])
		}
		// 角度时间系数是250.0/32768.0 dps
		deviation := (finalPosition - count)/finalPosition
		fmt.Println("count and dev",count,"and",deviation)
		switch {
		case deviation < 0:
			sign = - sign
			speed /= 2
		case deviation<0.1 || speed == 0:
			return e.Stop()
		case deviation < 0.5:
			speed /= 2
		}
	}
}
