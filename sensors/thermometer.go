package sensors

import (
	"errors"
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
	//defer ln.Close()

	return &Thermometer{
		name: "testThermometer",
		addr: remoteAddr,
	}, err
}

// GetTemp gets the current temperature value from the thermometer
func (t *Thermometer) GetTemp() (float64, error) {
	// Ensure we have a connection
	if t.conn == nil {
		return 0, errors.New("connection not setup")
	}

	var buffer []byte
	var length = 0
	temp := make([]byte, 128)

	tempLength, _ := t.conn.Read(temp)
	buffer = temp
	length = tempLength

	values := strings.Split(string(buffer[:length]), ",")
	probe1, err := strconv.Atoi(values[0])
	if err != nil {
		fmt.Println("Non-int value received")
	}
	value := float64(probe1) / 100.0
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
