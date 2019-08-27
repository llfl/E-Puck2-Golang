package main

import (
	"golang.org/x/exp/io/i2c"
	// "bytes"
	// "binary"
	"fmt"
)

/***************常量*****************/

// I2CDevices I2C设备
const I2CDevices string = "/dev/i2c-4"

// I2CAddress I2C设备的地址
const I2CAddress = 0x39

// ActuatorSize 定义控制器的长度
const ActuatorSize int = (19 + 1)

/***************变量*****************/
var actuator = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

/**************主函数*****************/

func main() {
	var counter, actuatorState int
	device, err := i2c.Open(&i2c.Devfs{Dev: I2CDevices}, I2CAddress)
	if err != nil {
		panic(err)
	}
	for true {
		counter = (counter + 1) % 2000
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
				fmt.Printf("checksum: %d\n", checksum)
			}
			actuator[ActuatorSize-1] = checksum
		}
		device.Write(actuator)
		fmt.Println("actuator updated!")
	}

}

// // Uint8ToByte 讲uint8转换成字节串
// func Uint8ToByte(n []uint8) []byte {
// 	bytesBuffer := bytes.NewBuffer([]byte{})

// }
