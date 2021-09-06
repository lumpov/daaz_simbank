package trade

import (
	"fmt"
	"milliard-easy/daaz_simbank/context"
	"milliard-easy/daaz_simbank/daaz"
	"milliard-easy/daaz_simbank/internal/email"
	"milliard-easy/daaz_simbank/log"
	"milliard-easy/daaz_simbank/telegram"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tarm/serial"
	"github.com/warthog618/modem/at"
)

const (
	baud = 115200
)

// Controller for trade system
// Automatically get phone number, create/close wallet, close payments etc..
type Controller struct {
	ports            []string
	workers          map[string]chan struct{}
	bot              *telegram.Bot
	api              *daaz.API
	cfg              *context.Config
	bufMsgController *telegram.BufMsgController
}

// NewController instance
func NewController(cfg *context.Config, ports []string, bot *telegram.Bot, api *daaz.API, bufMsgController *telegram.BufMsgController) *Controller {
	return &Controller{
		cfg:              cfg,
		ports:            ports,
		bot:              bot,
		api:              api,
		workers:          make(map[string]chan struct{}),
		bufMsgController: bufMsgController,
	}
}

// Start trading
func (c *Controller) Start() {
	wg := sync.WaitGroup{}
	wg.Add(len(c.ports))
	for i := range c.ports {
		go func(i int) {
			defer wg.Done()
			c.executePort(c.ports[i])
		}(i)
	}
	wg.Wait()
}

func (c *Controller) executePort(port string) {
	l := logrus.WithField("port", port)

	comPort, err := openPort(port, baud)
	if err != nil {
		l.WithError(err).Errorf(log.ErrorColor, "Cannot open port")
		c.bufMsgController.SendBufMsg(telegram.NewBufMsgCannotOpenPort(port))
		return
	}

	gsm := at.New(comPort)
	stateC := pingPort(gsm, port)
	portC := NewPortController(gsm, comPort, c.cfg, c.api, c.bufMsgController, port)

	for state := range stateC {
		if state {
			if _, ok := c.workers[port]; !ok {
				c.workers[port] = portC.StartTrade()
			}
		} else {
			if _, ok := c.workers[port]; ok {
				close(c.workers[port])
			}
		}
	}
}

func openPort(port string, baud int) (*serial.Port, error) {
	work := make(chan *serial.Port)

	go func() {
		for {
			comPort, err := serial.OpenPort(&serial.Config{
				Name: port,
				Baud: baud,
			})

			if err == nil {
				work <- comPort
			} else {
				time.Sleep(time.Second * 10)
			}
		}
	}()

	select {
	case <-time.After(time.Minute):
		return nil, fmt.Errorf("timeout")
	case p := <-work:
		return p, nil
	}
}

func (c *Controller) notify(port, str string) {
	email.SendEmail(c.cfg, fmt.Sprintf("[%s] New notification", port), str)
	c.bot.Send(str)
}

func pingPort(gsm *at.AT, port string) chan bool {
	stateC := make(chan bool)
	var pingM sync.Mutex

	go func() {
		lastState := false

		changeState := func(v bool) {
			pingM.Lock()

			if lastState != v {
				lastState = v
				stateC <- lastState
			}

			pingM.Unlock()
		}

		for {
			time.Sleep(time.Second * 3)
			_, err := gsm.Command("")
			if lastState && err == at.ErrDeadlineExceeded {
				changeState(false)
			} else if !lastState && err == nil {
				changeState(true)
			}
		}
	}()

	return stateC
}
