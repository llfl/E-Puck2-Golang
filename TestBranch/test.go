package main

import(
	driver "github.com/llfl/E-Puck2-Golang/driverUtils"
	// "fmt"
	"time"
)

func main()  {
	epuck := driver.NewEPuckHandle(driver.EnableGyro())
	// epuck := driver.NewEPuckHandle()

	actuatorState := 0
	epuck.CalibrateGyro()
	for true{
		switch actuatorState {
		case 0:
			epuck.Spin(90)
			
			actuatorState = 1
		case 1:
			epuck.Stop()
			
			actuatorState = 2
		case 2:
			epuck.FreeSpin(-100)
			
			actuatorState = 3
		case 3:
			epuck.Stop()
			
			actuatorState = 0
		}
		// epuck.Stop()
		// // fmt.Println(res)
		// epuck.UpdateGyro()
		// fmt.Println(epuck.Gyro)
		time.Sleep(1 * time.Second)
	}
	// fmt.Println("hello world")
}