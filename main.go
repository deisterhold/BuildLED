package main

import (
	"fmt"
	"time"
)

var unicorn = Unicorn{8, 4}

func init() {
	fmt.Println("Test")
}

func main() {
	unicorn.Init()
	defer unicorn.CleanUp()

	exit := make(chan bool, 1)
	go startServer(exit)

	buildTimer := time.NewTicker(time.Second * 10)
	danceTimer := time.NewTicker(time.Hour)

	for {
		select {
		case <-buildTimer.C:
			fmt.Println("Build")
			break
		case <-danceTimer.C:
			fmt.Println("Dance")
			break
		}
	}

	// status, err := FetchBuild(server1, credentials1)
}

func updateStatus() {
	for y := 0; y < 4; y++ {
		server := servers[y]
		credential := credentials[y]
		status, _ := FetchBuild(server, credential)
		for x := 0; x < 8; x++ {
			updateLED(x, y, status[x])
		}
	}
}

func updateLED(x, y int, status TFSBuildStatus) {
	switch status.Status {
	case "succeeded":
		unicorn.SetPixel(x, y, 0, 64, 0)
		break
	case "partiallySucceeded":
		unicorn.SetPixel(x, y, 64, 32, 0)
		break
	case "failed":
		unicorn.SetPixel(x, y, 64, 0, 0)
		break
	case "canceled":
		unicorn.SetPixel(x, y, 64, 0, 64)
		break
	case "inProgress":
		unicorn.SetPixel(x, y, 0, 0, 64)
		break
	default:
		unicorn.SetPixel(x, y, 0, 0, 0)
		break
	}
}
