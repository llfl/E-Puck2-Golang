package main

import(
	driver "github.com/llfl/E-Puck2-Golang/driverUtils"
	// "fmt"
	"time"
)

func main()  {
	epuck := driver.NewEPuckHandle()

	actuatorState := 0
	for true{
		switch actuatorState {
		case 0:
			epuck.FreeSpin(true)
			actuatorState = 1
		case 1:
			epuck.Stop()
			actuatorState = 2
		case 2:
			epuck.FreeSpin(false)
			actuatorState = 3
		case 3:
			epuck.Stop()
			actuatorState = 0
		}

		time.Sleep(1 * time.Second)
	}
	// fmt.Println("hello world")
}