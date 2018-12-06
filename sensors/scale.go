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
var scaleVal *ScaleValue
var quitScale chan int

// Scale should be treated as a singleton
type Scale struct {
	name string
	addr *net.UDPAddr
	conn *net.UDPConn
}

type ScaleValue struct {
	val float64
}

// setupScale connects to the physical sensor
func (Scale) setupScale() (*Scale, error) {
	var err error
	s := Scale{name: "smart-chef-scale"}
	scaleVal = &ScaleValue{0}
	quitScale = make(chan int)

	s.addr, err = net.ResolveUDPAddr("udp", os.Getenv("SCALE_ADDR"))
	if err != nil {
		return nil, err
	}

	go s.fetchFromSocket()
	log.Info("Done setting up scale")
	return &s, nil
}

func (*Scale) Cleanup() {
	quitScale <- 0
}

func (s *Scale) fetchFromSocket() {
	var err error
	s.conn, err = net.ListenUDP("udp", s.addr)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer s.conn.Close()

	for {
		select {
		case <-quitScale:
			log.Info("Closing Scale connection")
			return
		default:
			var buffer []byte
			var length = 0
			weight := make([]byte, 128)

			weightength, err := s.conn.Read(weight)
			if err != nil {
				log.Error(err.Error())
				break
			}

			buffer = weight
			length = weightength

			value, err := strconv.Atoi(string(buffer[:length]))
			if err != nil {
				log.Error("Non-int value received")
				break
			}
			scaleVal.val = float64(value)
		}
	}
}

// GetWeight gets the current weight value from teh scale sensor
func (s *Scale) GetWeight() (float64, error) {
	v := scaleVal.val
	log.Info("Scale Value: %d", v)
	return v, nil
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
