package sensors

import (
	"fmt"
	"net"
	"os"
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
	conn *net.UDPConn
}

// setupThermometer connects to the physical thermometer
func (Thermometer) setupThermometer() (*Thermometer, error) {
	remoteAddr, err := net.ResolveUDPAddr("udp", os.Getenv("THERMOMETER_ADDR"))
	ln, err := net.ListenUDP("udp", remoteAddr)
	if err != nil {
		return nil, err
	}
	ln.SetReadBuffer(16)
	ln.SetWriteBuffer(16)
	log.Infof("Established connection to %s \n", remoteAddr)
	log.Infof("Local UDP client address : %s \n", ln.LocalAddr().String())
	// Keep this open all the time?
	defer ln.Close()

	return &Thermometer{
		name: "testThermometer",
		addr: remoteAddr,
		conn: ln,
	}, err
}

// GetTemp gets the current temperature value from the thermometer
func (t *Thermometer) GetTemp() (float64, error) {
	remoteAddr, err := net.ResolveUDPAddr("udp", os.Getenv("THERMOMETER_ADDR"))
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}
	ln, err := net.ListenUDP("udp", remoteAddr)
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}
	defer ln.Close()

	var buffer []byte
	var length = 0
	temp := make([]byte, 128)

	tempLength, _ := ln.Read(temp)
	buffer = temp
	length = tempLength
	// }

	// fmt.Println("UDP Server : ", addr)
	values := strings.Split(string(buffer[:length]), ",")
	probe1, err := strconv.Atoi(values[0])
	if err != nil {
		log.Error("Non-int value received")
		return 0, err
	}
	value := float64(probe1) / 100.0

	if value == 330.36 {
		log.Warn("Received Null temperature, disregarding")
		return 0, nil
	}

	return value, nil
}

// Implement Singleton GetInstance
func (*Thermometer) GetInstance() (*Thermometer, error) {
	var err error
	err = nil
	thermometerOnce.Do(func() {
		thermometerInstance, err = new(Thermometer).setupThermometer()
	})
	return thermometerInstance, err
}
