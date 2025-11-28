package main

import (
	_ "embed"

	"log"

	"strings"

	"github.com/getlantern/systray"
	"github.com/ncruces/zenity"
)

//go:embed icon.ico
var iconData []byte

func setupTray() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("NFC Tools for Arduino")
	systray.SetTooltip("Read and Write to any NFC Tag from here using Arduino")

	systray.SetIcon(iconData)

	mRead := systray.AddMenuItem("Read NFC", "Read all Current Data on NFC Tag")

	sWrite := systray.AddMenuItem("Write to NFC", "Write a record to the NFC Tag")
	mWriteText := sWrite.AddSubMenuItem("Text", "Write plain text to NFC")

	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "Stop Running")

	go func() {
		for {
			select {
			case <-mRead.ClickedCh:
				if err := Device.Exec("nfcRead"); err != nil {
					log.Println(err)
				}
			case <-mWriteText.ClickedCh:
				text, err := zenity.Entry("Enter any text:")
				if err != nil && !strings.Contains(err.Error(), "dialog canceled") {
					log.Println(err)
					return
				}

				if err := Device.Execf("nfcWriteText %s", text); err != nil {
					log.Println(err)
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	if Device != nil {
		Device.Close()
	}
}
