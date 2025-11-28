package main

import (
	"fmt"
	"log"

	"time"

	"github.com/ncruces/zenity"
	"github.com/tarm/serial"
	arduinocomm "github.com/z-sk1/arduino-comm"
)

var Device *arduinocomm.Device

func main() {
	portName := askForPort()
	openPort(portName)

	go readSerialLoop()
	setupTray()
}

func askForPort() string {
	port, err := zenity.Entry("Enter a COM Port: (e.g COM3, COM5)")
	if err != nil {
		log.Fatal(err)
	}
	return port
}

func openPort(comPort string) {
	c := &serial.Config{Name: comPort, Baud: 9600}
	port, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	Device = arduinocomm.New(port)
}

func readSerialLoop() {
	fmt.Println("Started reading serial...")
	buf := make([]byte, 128)
	var lineBuf []byte

	for {
		n, err := Device.Port.Read(buf)
		if err != nil {
			fmt.Println("Read error:", err)
			continue
		}

		if n > 0 {
			for _, b := range buf[:n] {
				if b == '\n' {
					fmt.Println(string(lineBuf))
					lineBuf = nil
				} else if b != '\r' {
					lineBuf = append(lineBuf, b)
				}
			}
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}
