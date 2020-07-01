package main

import(
	driver "github.com/llfl/E-Puck2-Golang/driverUtils"
	"fmt"
)

func main()  {
	epuck := driver.NewEPuckHandle()
	count := 0
	actuatorState := 0
	for true{
		count = (count + 1)%20
		if count == 0{
			switch actuatorState {
			case 0:
				epuck.spin(90.0)
				actuatorState = 1
			case 1:
				epuck.spin(90.0)
				actuatorState = 2
			case 2:
				epuck.spin(-90.0)
				actuatorState = 3
			case 3:
				epuck.spin(-90.0)
				actuatorState = 0
			}
		}
	}
}