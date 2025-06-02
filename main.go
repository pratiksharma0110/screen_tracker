package main

import (
	"encoding/json"
	"fmt"
	"log"

	"os"
	"os/signal"
	"path/filepath"
	"screen_tracker/internal/hypr"
	"screen_tracker/internal/socket"
	"screen_tracker/internal/tracker"
	"screen_tracker/pkg/utils"
	"syscall"
)

func main() {
	HIS := utils.GetEnv("HYPRLAND_INSTANCE_SIGNATURE")
	XDG := utils.GetEnv("XDG_RUNTIME_DIR")
	sockPath := filepath.Join(XDG, "hypr", HIS, ".socket2.sock") //socket2 as its event based

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, os.Interrupt)

	data := &tracker.AppData{}

	done := make(chan struct{})
	windowChan := make(chan map[string]float64)

	go func() {
		<-sigs
		fmt.Println("\nReceived interrupt, stopping...")
		close(done)
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				buf, err := hypr.HyprIPC(sockPath)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					continue
				}
				app := hypr.GetActiveWindow(buf)
				appDetails := data.TrackTimer(app)

				windowChan <- appDetails //window channel recevis map in format of {app:time in sec}}, handle time formatting from client side
			}
		}
	}()

	conn := internal.CreateSockt(internal.SOCKET_TYPE, internal.SOCKET_PATH)

	for {
		select {
		case <-done:
			fmt.Println("Exiting..")
			return
		case app := <-windowChan:
			jsonData, err := json.Marshal(app)
			if err != nil {
				log.Println("Failed to marshal app data:", err)
				continue
			}

			//pass info to client connected
			_, err = conn.Write(append(jsonData, '\n'))
			if err != nil {
				log.Println("Error writing to client:", err)
				return
			}

		}
	}

}
