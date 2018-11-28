package sensors

import (
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var thermometerInstance *Thermometer
var thermometerOnce sync.Once
var thermVal *TempValue
var quitTemp chan int

// Scale should be treated as a singleton
type Thermometer struct {
	name string
	addr *net.UDPAddr
	conn *net.UDPConn
	quit chan int
}

type TempValue struct {
	val float64
}

// setupThermometer connects to the physical thermometer
func setupThermometer() (*Thermometer, error) {
	var err error
	t := Thermometer{name: "range-driver"}
	thermVal = &TempValue{0}
	quitTemp = make(chan int)

	t.addr, err = net.ResolveUDPAddr("udp", os.Getenv("THERMOMETER_ADDR"))
	if err != nil {
		return nil, err
	}

	go t.fetchFromSocket()

	log.Info("Done setting up thermometer")
	return &t, err
}

// Cleanup the thermometer
func (t *Thermometer) Cleanup() {
	quitTemp <- 0
}

func (t *Thermometer) fetchFromSocket() {
	var err error
	t.conn, err = net.ListenUDP("udp", t.addr)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer t.conn.Close()

	for {
		select {
		case <-quitTemp:
			log.Info("Closing Thermometer connection")
			return
		default:
			var buffer []byte
			var length = 0
			temp := make([]byte, 128)

			tempLength, _ := t.conn.Read(temp)
			buffer = temp
			length = tempLength

			values := strings.Split(string(buffer[:length]), ",")
			probe1, err := strconv.Atoi(values[0])
			if err != nil {
				log.Error("Non-int value received")
				break
			}
			value := float64(probe1) / 100.0

			if value == 330.36 {
				log.Warn("Received Null temperature, disregarding")
				break
			}
			thermVal.val = value
			time.Sleep(time.Second / 2)
		}
	}
}

// GetTemp gets the current temperature value from the thermometer
func (t *Thermometer) GetTemp() (float64, error) {
	v := thermVal.val
	log.Infof("Thermomter Value: %d", v)
	return v, nil
}

// Implement Singleton GetInstance
func (*Thermometer) GetInstance() (*Thermometer, error) {
	var err error
	err = nil
	thermometerOnce.Do(func() {
		thermometerInstance, err = setupThermometer()
	})
	return thermometerInstance, err
}
