package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"screen_tracker/internal/hypr"
	"screen_tracker/internal/tracker"
	"screen_tracker/pkg/utils"
	"syscall"
	"time"
)

func main() {
	HIS := utils.GetEnv("HYPRLAND_INSTANCE_SIGNATURE")
	XDG := utils.GetEnv("XDG_RUNTIME_DIR")
	sockPath := filepath.Join(XDG, "hypr", HIS, ".socket2.sock")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, os.Interrupt)

	data := &tracker.AppData{}

	done := make(chan struct{})
	windowChan := make(chan map[string]time.Duration)

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

				windowChan <- appDetails
			}
		}
	}()

	for {
		select {
		case <-done:
			fmt.Println("Exiting..")
			return
		case app := <-windowChan:
			data.Print(app)
		}
	}
}
