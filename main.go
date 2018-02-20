package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var unicorn = Unicorn{8, 4}

func init() {
	servers[0] = TFSOnPremServer{os.Getenv("TFS_INSTANCE"), TFSBuildDefinition{os.Getenv("TFS_COLLECTION"), os.Getenv("TFS_PROJECT_1"), os.Getenv("TFS_BUILD_ID_1")}}
	servers[1] = TFSOnPremServer{os.Getenv("TFS_INSTANCE"), TFSBuildDefinition{os.Getenv("TFS_COLLECTION"), os.Getenv("TFS_PROJECT_2"), os.Getenv("TFS_BUILD_ID_2")}}
	servers[2] = TFSOnPremServer{os.Getenv("TFS_INSTANCE"), TFSBuildDefinition{os.Getenv("TFS_COLLECTION"), os.Getenv("TFS_PROJECT_3"), os.Getenv("TFS_BUILD_ID_3")}}
	servers[3] = TFSOnPremServer{os.Getenv("TFS_INSTANCE"), TFSBuildDefinition{os.Getenv("TFS_COLLECTION"), os.Getenv("TFS_PROJECT_4"), os.Getenv("TFS_BUILD_ID_4")}}
	credentials[0] = TFSCredentials{os.Getenv("TFS_USERNAME"), os.Getenv("TFS_PASSWORD"), os.Getenv("TFS_DOMAIN")}
	credentials[1] = TFSCredentials{os.Getenv("TFS_USERNAME"), os.Getenv("TFS_PASSWORD"), os.Getenv("TFS_DOMAIN")}
	credentials[2] = TFSCredentials{os.Getenv("TFS_USERNAME"), os.Getenv("TFS_PASSWORD"), os.Getenv("TFS_DOMAIN")}
	credentials[3] = TFSCredentials{os.Getenv("TFS_USERNAME"), os.Getenv("TFS_PASSWORD"), os.Getenv("TFS_DOMAIN")}
}

func main() {
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM)
		<-sigs
		unicorn.CleanUp()
	}()

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
		status, err := FetchBuild(servers[y], credentials[y])

		if err != nil {
			fmt.Println("Error:", err.Error())
		}

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
	switch status.Status {
	case "completed":
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
		default:
			unicorn.SetPixel(x, y, 0, 0, 0)
			fmt.Println(status)
			break
		}
	case "inProgress":
		unicorn.SetPixel(x, y, 0, 0, 64)
		break
	default:
		unicorn.SetPixel(x, y, 0, 0, 0)
		fmt.Println(status)
		break
	}
}
