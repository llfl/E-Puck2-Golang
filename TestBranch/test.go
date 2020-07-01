package main

import (
	"golang.org/x/exp/io/i2c"
	// "bytes"
	// "binary"
	"fmt"
	"time"
	"reflect"
)

/***************常量*****************/

// I2CDevices I2C设备
const I2CDevices string = "/dev/i2c-4"

// I2CAddress I2C设备的地址
const I2CAddress = 0x1F

// ActuatorSize 定义控制器的长度 with checksum
const ActuatorSize int = (19 + 1)

// SensorDataSize 定义传感器数据的长度 with checksum
const SensorDataSize int = (46 + 1)

/***************变量*****************/

// var actuator = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var actuator = make([]byte, ActuatorSize)

var sensorData = make([]byte, SensorDataSize)

/*********epuck传感器结构体************/

type epuckSensor struct {
	prox      [8]uint16
	ambient   [8]uint16
	mic       [4]uint16
	sel       uint8
	button    uint8
	motorStep [2]int16
	tv        uint8
}

/**************主函数*****************/

func main() {
	var counter, actuatorState int
	device, err := i2c.Open(&i2c.Devfs{Dev: I2CDevices}, I2CAddress)
	if err != nil {
		panic(err)
	}
	for true {
		counter = (counter + 1) % 20
		if counter == 0 {
			switch actuatorState {
			case 0:
				actuator[0] = 0 // Left speed: 512
				actuator[1] = 2
				actuator[2] = 0 // Right speed: -512
				actuator[3] = 0xFE
				actuator[4] = 0    // Speaker sound
				actuator[5] = 0x0F // LED1, LED3, LED5, LED7 on/off flag
				actuator[6] = 100  // LED2 red
				actuator[7] = 0    // LED2 green
				actuator[8] = 0    // LED2 blue
				actuator[9] = 100  // LED4 red
				actuator[10] = 0   // LED4 green
				actuator[11] = 0   // LED4 blue
				actuator[12] = 100 // LED6 red
				actuator[13] = 0   // LED6 green
				actuator[14] = 0   // LED6 blue
				actuator[15] = 100 // LED8 red
				actuator[16] = 0   // LED8 green
				actuator[17] = 0   // LED8 blue
				actuator[18] = 0   // Settings.
				actuatorState = 1
			case 1:
				actuator[0] = 0 // Left speed: 512
				actuator[1] = 0
				actuator[2] = 0 // Right speed: -512
				actuator[3] = 0x00
				actuator[4] = 0    // Speaker sound
				actuator[5] = 0x00 // LED1, LED3, LED5, LED7 on/off flag
				actuator[6] = 0    // LED2 red
				actuator[7] = 100  // LED2 green
				actuator[8] = 0    // LED2 blue
				actuator[9] = 0    // LED4 red
				actuator[10] = 100 // LED4 green
				actuator[11] = 0   // LED4 blue
				actuator[12] = 0   // LED6 red
				actuator[13] = 100 // LED6 green
				actuator[14] = 0   // LED6 blue
				actuator[15] = 0   // LED8 red
				actuator[16] = 100 // LED8 green
				actuator[17] = 0   // LED8 blue
				actuator[18] = 0   // Settings.
				actuatorState = 2
			case 2:
				actuator[0] = 0 // Left speed: 512
				actuator[1] = 2
				actuator[2] = 0 // Right speed: -512
				actuator[3] = 0xFE
				actuator[4] = 0    // Speaker sound
				actuator[5] = 0x0F // LED1, LED3, LED5, LED7 on/off flag
				actuator[6] = 0    // LED2 red
				actuator[7] = 0    // LED2 green
				actuator[8] = 100  // LED2 blue
				actuator[9] = 0    // LED4 red
				actuator[10] = 0   // LED4 green
				actuator[11] = 100 // LED4 blue
				actuator[12] = 0   // LED6 red
				actuator[13] = 0   // LED6 green
				actuator[14] = 100 // LED6 blue
				actuator[15] = 0   // LED8 red
				actuator[16] = 0   // LED8 green
				actuator[17] = 100 // LED8 blue
				actuator[18] = 0   // Settings.
				actuatorState = 3
			case 3:
				actuator[0] = 0 // Left speed: 0
				actuator[1] = 0
				actuator[2] = 0 // Right speed: 0
				actuator[3] = 0
				actuator[4] = 0    // Speaker sound
				actuator[5] = 0x0  // LED1, LED3, LED5, LED7 on/off flag
				actuator[6] = 100  // LED2 red
				actuator[7] = 100  // LED2 green
				actuator[8] = 0    // LED2 blue
				actuator[9] = 100  // LED4 red
				actuator[10] = 100 // LED4 green
				actuator[11] = 0   // LED4 blue
				actuator[12] = 100 // LED6 red
				actuator[13] = 100 // LED6 green
				actuator[14] = 0   // LED6 blue
				actuator[15] = 100 // LED8 red
				actuator[16] = 100 // LED8 green
				actuator[17] = 0   // LED8 blue
				actuator[18] = 0   // Settings.
				actuatorState = 0
			}
			var checksum uint8
			for i := 0; i < (ActuatorSize - 1); i++ {
				checksum ^= actuator[i]
			}
			actuator[ActuatorSize-1] = checksum

			for trails := 0; trails < 3; trails++ {
				if err := device.Write(actuator); err != nil {
					continue
				}
				break
			}
			if err != nil {
				panic(err)
			}

			fmt.Println("actuator updated!")

			if err := device.Read(sensorData); err != nil {
				panic(err)
			}
			epuckData := sensorDataParser()
			fmt.Println(epuckData)
		}
		time.Sleep(50 * time.Millisecond)
	}

}

func sensorDataParser() epuckSensor {
	var epuckData epuckSensor
	for i := 0; i < 8; i++ {
		epuckData.prox[i] = uint16(sensorData[i*2+1])*256 + uint16(sensorData[2*i])
		epuckData.ambient[i] = uint16(sensorData[16+i*2+1])*256 + uint16(sensorData[16+2*i])
	}
	for i := 0; i < 4; i++ {
		epuckData.mic[i] = uint16(sensorData[32+i*2+1])*256 + uint16(sensorData[32+2*i])
	}
	epuckData.sel = sensorData[40] & 0x0F
	epuckData.button = sensorData[40] >> 4
	for i := 0; i < 2; i++ {
		epuckData.motorStep[i] = int16(uint16(sensorData[41+i*2+1])*256 + uint16(sensorData[41+2*i]))
	}
	epuckData.tv = sensorData[45]
	return epuckData
}
