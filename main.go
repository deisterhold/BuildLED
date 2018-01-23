package main

import (
	"fmt"
	"os"
	"time"
)

var unicorn = Unicorn{8, 4}

func init() {
	for i := 0; i < 4; i++ {
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
		}
	}
}

func dance() {
	unicorn.SetPixel(0, 0, 16, 16, 16)
	unicorn.SetPixel(0, 3, 0, 0, 16)
	unicorn.SetPixel(7, 0, 0, 16, 0)
	unicorn.SetPixel(7, 3, 16, 0, 0)
	unicorn.Show()
	time.Sleep(time.Second * 2)

	for y := 0; y < 4; y++ {
		for x := 0; x < 8; x++ {
			unicorn.SetPixel(x, y, 64, 0, 0)
			unicorn.Show()
			time.Sleep(time.Millisecond * 50)
		}
	}
	for y := 0; y < 4; y++ {
		for x := 0; x < 8; x++ {
			unicorn.SetPixel(x, y, 0, 64, 0)
			unicorn.Show()
			time.Sleep(time.Millisecond * 50)
		}
	}
	for y := 0; y < 4; y++ {
		for x := 0; x < 8; x++ {
			unicorn.SetPixel(x, y, 0, 0, 64)
			unicorn.Show()
			time.Sleep(time.Millisecond * 50)
		}
	}
	unicorn.Clear()
	unicorn.Show()
}

func updateStatus() {
	for y := 0; y < 4; y++ {
		server := servers[y]
		credential := credentials[y]
		status, _ := FetchBuild(server, credential)
		fmt.Println(y, status)
		for x := 0; x < 8; x++ {
			if x < len(status) {
				updateLED(8-x, y, status[x])
			} else {
				unicorn.SetPixel(8-x, y, 0, 0, 0)
			}
		}
	}
	unicorn.Show()
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
		break
	}

	fmt.Println(x, y, status.Status)
}
