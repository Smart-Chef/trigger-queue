package sensors

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

var thermometerInstance *Thermometer
var thermometerOnce sync.Once

func init() {
	// todo start python driver
	fmt.Println("")

}

// Scale should be treated as a singleton
type Thermometer struct {
	name string
	addr *net.UDPAddr

	// TODO: add other metadata
}

// setupThermometer connects to the physical thermometer
func (Thermometer) setupThermometer() *Thermometer {
	// TODO: add code to connect to the actual sensor
	remoteAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10000")
	ln, err := net.ListenUDP("udp", remoteAddr)
	if err != nil {
		log.Fatal(err)
	}
	ln.SetReadBuffer(16)
	ln.SetWriteBuffer(16)
	fmt.Printf("Established connection to %s \n", remoteAddr)
	// log.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
	log.Printf("Local UDP client address : %s \n", ln.LocalAddr().String())
	defer ln.Close()
	//todo set up socket
	return &Thermometer{
		name: "testThermometer",
		addr: remoteAddr,
	}
}

// GetTemp gets the current temperature value from the thermometer
func (t *Thermometer) GetTemp() float64 {
	remoteAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10000")
	ln, err := net.ListenUDP("udp", remoteAddr)

	if err != nil {
		log.Error(err.Error())
	}
	temp := make([]byte, 128)
	var buffer []byte
	var length = 0

	// for {
	tempLength, err := ln.Read(temp)

	if err != nil {
		log.Error(err.Error())
	}
	// if err != nil {
	// 	break
	// }
	buffer = temp
	length = tempLength
	// }

	// fmt.Println("UDP Server : ", addr)
	values := strings.Split(string(buffer[:length]), ",")
	fmt.Println("Received from UDP server : ", values[0]+","+values[1])
	probe1, err := strconv.Atoi(values[0])
	if err != nil {
		log.Error("Non-int value received")
	}
	value := float64(probe1) / 100.0
	defer ln.Close()

	if value == 330.36 {
		log.Error("Received Null temperature, disregarding")
		return 0
	}

	// fmt.Println(value)

	return value
}

// Implement Singleton GetInstance
func (*Thermometer) GetInstance() *Thermometer {
	thermometerOnce.Do(func() {
		thermometerInstance = new(Thermometer).setupThermometer()
	})
	return thermometerInstance
}
