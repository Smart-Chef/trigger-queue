package sensors

import (
	"net"
	"os"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
)

var scaleInstance *Scale
var scaleOnce sync.Once

// Scale should be treated as a singleton
type Scale struct {
	name string
	addr *net.UDPAddr
}

// setupScale connects to the physical sensor
func (Scale) setupScale() (*Scale, error) {
	remoteAddr, err := net.ResolveUDPAddr("udp", os.Getenv("SCALE_ADDR"))
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

	return &Scale{
		name: "testScale",
		addr: remoteAddr,
	}, nil
}

// GetWeight gets the current weight value from teh scale sensor
func (s *Scale) GetWeight() (float64, error) {
	ln, err := net.ListenUDP("udp", s.addr)
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

	value, err := strconv.Atoi(string(buffer[:length]))
	if err != nil {
		log.Error("Non-int value received")
		return 0, err
	}

	return float64(value), nil
}

// Implement Singleton GetInstance
func (*Scale) GetInstance() (*Scale, error) {
	var err error
	err = nil
	scaleOnce.Do(func() {
		scaleInstance, err = new(Scale).setupScale()
	})
	return scaleInstance, err
}
