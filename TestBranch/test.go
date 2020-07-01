package main

import(
	driver "github.com/llfl/E-Puck2-Golang/driverUtils"
	"fmt"
	"time"
)

func main()  {
	epuck := driver.NewEPuckHandle(driver.EnableGyro())

	actuatorState := 0
	epuck.CalibrateGyro()
	for true{
		switch actuatorState {
		case 0:
			epuck.Spin(100)
			epuck.UpdateSensors()
			fmt.Println("ok0",epuck.Sensors)
			fmt.Println(epuck.Gyro.Values)
			actuatorState = 1
		case 1:
			epuck.Stop()
			fmt.Println("ok1",epuck.Sensors)
			fmt.Println(epuck.Gyro.Values)
			actuatorState = 2
		case 2:
			epuck.Spin(-100)
			fmt.Println("ok2",epuck.Sensors)
			fmt.Println(epuck.Gyro.Values)
			actuatorState = 3
		case 3:
			epuck.Stop()
			fmt.Println("ok3",epuck.Sensors)
			fmt.Println(epuck.Gyro.Values)
			actuatorState = 0
		}
		// epuck.Stop()
		// // fmt.Println(res)
		// fmt.Println(epuck.Gyro)
		time.Sleep(1 * time.Second)
	}
	// fmt.Println("hello world")
}