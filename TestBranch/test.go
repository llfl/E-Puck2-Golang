package main

import(
	driver "github.com/llfl/E-Puck2-Golang/driverUtils"
	"fmt"
	"time"
)

func main()  {
	epuck := driver.NewEPuckHandle()

	actuatorState := 0
	for true{
		switch actuatorState {
		case 0:
			epuck.Spin(90)
			fmt.Println("ok0")
			actuatorState = 1
		case 1:
			epuck.Stop()
			fmt.Println("ok1")
			actuatorState = 2
		case 2:
			epuck.Spin(-90)
			fmt.Println("ok2")
			actuatorState = 3
		case 3:
			epuck.Stop()
			fmt.Println("ok3")
			actuatorState = 0
		}

		time.Sleep(3 * time.Second)
	}
	// fmt.Println("hello world")
}