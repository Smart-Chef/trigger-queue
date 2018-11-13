package sensors

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var thermometerInstance *Thermometer
var thermometerOnce sync.Once

func init() {
	// Load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// Scale should be treated as a singleton
type Thermometer struct {
	name string
	addr *net.UDPAddr

	// TODO: add other metadata
}

// setupThermometer connects to the physical thermometer
func (Thermometer) setupThermometer() *Thermometer {
	remoteAddr, err := net.ResolveUDPAddr("udp", os.Getenv("THERMOMETER_ADDR"))
	ln, err := net.ListenUDP("udp", remoteAddr)
	if err != nil {
		log.Fatal(err)
	}

	ln.SetReadBuffer(16)
	ln.SetWriteBuffer(16)
	log.Printf("Established connection to %s \n", remoteAddr)
	log.Printf("Local UDP client address : %s \n", ln.LocalAddr().String())
	defer ln.Close()

	// TODO: set up socket
	return &Thermometer{
		name: "testThermometer",
		addr: remoteAddr,
	}
}

// GetTemp gets the current temperature value from the thermometer
func (t *Thermometer) GetTemp() float64 {
	ln, err := net.ListenUDP("udp", t.addr)
	temp := make([]byte, 128)
	var buffer []byte
	var length = 0

	tempLength, _ := ln.Read(temp)
	buffer = temp
	length = tempLength

	values := strings.Split(string(buffer[:length]), ",")
	fmt.Println("Received from UDP server : ", values[0]+","+values[1])
	probe1, err := strconv.Atoi(values[0])
	if err != nil {
		fmt.Println("Non-int value received")
	}
	value := float64(probe1) / 100.0
	fmt.Println(value)
	defer ln.Close()

	return value
}

// Implement Singleton GetInstance
func (*Thermometer) GetInstance() *Thermometer {
	thermometerOnce.Do(func() {
		thermometerInstance = new(Thermometer).setupThermometer()
	})
	return thermometerInstance
}
