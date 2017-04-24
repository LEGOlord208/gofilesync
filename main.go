package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getlantern/systray"
	"github.com/legolord208/gofilesync/api"
	"github.com/legolord208/stdutil"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	var src string
	var dst string
	var lazy bool

	flag.StringVar(&src, "src", "", "The source folder to copy.")
	flag.StringVar(&dst, "dst", "", "The destination folder to paste.")
	flag.BoolVar(&lazy, "lazy", false, "Whether or not to only sync necessary files.")

	flag.Parse()

	if (src != "" && dst == "") || (dst != "" && src == "") {
		stdutil.PrintErr("'src' or 'dst' set, but not both.", nil)
		return
	}

	if src != "" && dst != "" {
		var err error
		if lazy {
			err = gofilesync.LazySync(src, dst)
		} else {
			err = gofilesync.ForceSync(src, dst)
		}

		if err != nil {
			stdutil.PrintErr("Error while syncing", err)
		}
		return
	}

	loadData()
	schedule(data.Schedule)
	go initWebserver()
	systray.Run(onReady)
}

var itemStatus *systray.MenuItem

func onReady() {
	systray.SetTooltip("gofilesync")

	chanStop := make(chan os.Signal, 2)
	signal.Notify(chanStop, os.Interrupt, syscall.SIGTERM)

	itemStatus = systray.AddMenuItem("", "Status")
	itemConfig := systray.AddMenuItem("Configure", "Open settings webpage")

	itemExit := systray.AddMenuItem("Exit", "Exit the application")

	status(false, "No issues")
	go func() {
		for {
			select {
			case <-itemStatus.ClickedCh:
				status(false, "No issues")
			case <-itemConfig.ClickedCh:
				open.Run("http://localhost" + port + "/")
			case <-chanStop:
				stop()
			case <-itemExit.ClickedCh:
				stop()
			}
		}
	}()
}

func status(err bool, message string) {
	prefix := "Status"
	if err {
		systray.SetIcon(iconErr)
		prefix = "Error"
	} else {
		systray.SetIcon(icon)
	}
	itemStatus.SetTitle(prefix + ": " + message)
}

func stop() {
	systray.Quit() // Stays awake for some reason? Let's force stop it ¯\_(ツ)_/¯
	scheduleStop()

	time.Sleep(time.Second)
	os.Exit(0)
}
