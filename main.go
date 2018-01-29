package main

import (
	"fmt"
	"os"
	"time"
)

var unicorn = Unicorn{8, 4}

func init() {
	for i := 0; i < unicorn.Height; i++ {
		servers[i] = TFSHostedServer{os.Getenv("VSTS_ACCOUNT"), TFSBuildDefinition{"", os.Getenv("VSTS_PROJECT"), os.Getenv("VSTS_BUILD_ID")}}
		credentials[i] = TFSCredentials{os.Getenv("VSTS_USERNAME"), os.Getenv("VSTS_PASSWORD")}
	}
}

func main() {
	unicorn.Init()
	defer unicorn.CleanUp()
	unicorn.Clear()
	unicorn.Show()

	exit := make(chan bool, 1)
	go startServer(exit)

	buildTimer := time.NewTicker(time.Second * 10)
	danceTimer := time.NewTicker(time.Hour)

	dance()
	updateStatus()

	for {
		select {
		case <-buildTimer.C:
			fmt.Println("Build")
			updateStatus()
			break
		case <-danceTimer.C:
			fmt.Println("Dance")
			dance()
			break
		default:
			time.Sleep(time.Second)
		}
	}
}

func dance() {
	for y := 0; y < unicorn.Height; y++ {
		for x := 0; x < unicorn.Width; x++ {
			unicorn.SetPixel(x, y, 64, 0, 0)
			unicorn.Show()
			time.Sleep(time.Millisecond * 25)
		}
	}
	for y := 0; y < unicorn.Height; y++ {
		for x := 0; x < unicorn.Width; x++ {
			unicorn.SetPixel(x, y, 0, 64, 0)
			unicorn.Show()
			time.Sleep(time.Millisecond * 25)
		}
	}
	for y := 0; y < unicorn.Height; y++ {
		for x := 0; x < unicorn.Width; x++ {
			unicorn.SetPixel(x, y, 0, 0, 64)
			unicorn.Show()
			time.Sleep(time.Millisecond * 25)
		}
	}
	unicorn.Clear()
	unicorn.Show()
}

func updateStatus() {
	for y := 0; y < unicorn.Height; y++ {
		status, _ := FetchBuild(servers[y], credentials[y])
		for x := 0; x < unicorn.Width; x++ {
			if x < len(status) {
				updateLED(unicorn.Width-(x+1), y, status[x])
			} else {
				updateLED(unicorn.Width-(x+1), y, TFSBuildStatus{})
			}
		}
		unicorn.Show()
	}
}

func updateLED(x, y int, status TFSBuildStatus) {
	switch status.Result {
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
		fmt.Println(status)
		break
	}
}
