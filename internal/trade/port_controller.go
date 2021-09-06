package trade

import (
	"encoding/hex"
	"fmt"
	"milliard-easy/daaz_simbank/context"
	"milliard-easy/daaz_simbank/daaz"
	"milliard-easy/daaz_simbank/internal/email"
	"milliard-easy/daaz_simbank/log"
	"milliard-easy/daaz_simbank/telegram"
	"regexp"
	"strings"
	"time"

	"github.com/reactivex/rxgo/v2"
	"github.com/sirupsen/logrus"
	"github.com/tarm/serial"
	"github.com/warthog618/modem/at"
	"github.com/xlab/at/sms"
)

// PortController for trade on custom port
type PortController struct {
	gsm              *at.AT
	serial           *serial.Port
	api              *daaz.API
	cfg              *context.Config
	port             string
	log              *logrus.Entry
	bufMsgController *telegram.BufMsgController
}

// NewPortController instance
func NewPortController(gsm *at.AT, serial *serial.Port, cfg *context.Config, api *daaz.API, bufMsgController *telegram.BufMsgController, port string) *PortController {
	return &PortController{
		gsm:              gsm,
		serial:           serial,
		api:              api,
		port:             port,
		cfg:              cfg,
		log:              logrus.WithField("port", port),
		bufMsgController: bufMsgController,
	}
}

// StartTrade for port
func (c *PortController) StartTrade() chan struct{} {
	stopC := make(chan struct{})

	go func() {
		c.log.Infof(log.InfoColor, "Starting trade")
		defer c.log.Infof(log.InfoColor, "Stopping trade")

		c.log.Infof(log.InfoColor, "Autoconfig successfully")

		paySMSc, numSMSc, _, stopSMS := c.listenSMS()
		defer close(stopSMS)

		stopSMSlog := c.logSMS(rxgo.Merge([]rxgo.Observable{paySMSc, numSMSc}))
		defer close(stopSMSlog)

		stopSMSredirect := c.redirectSMS(rxgo.Merge([]rxgo.Observable{paySMSc, numSMSc}))
		defer close(stopSMSredirect)

		phoneC, forceStop := c.detectPhone(numSMSc, stopC)
		phone := <-phoneC
		close(forceStop)

		// c.setupTariff(tariffSMSc, phone)
		// c.makeCall()

		stopDisableCall := c.disableCall()
		defer close(stopDisableCall)

		stopSMSpay := c.autopay(paySMSc, phone)
		defer close(stopSMSpay)

		c.api.AddWallet(phone, c.port)

		c.api.ToggleWalletState(phone, true)
		defer c.api.ToggleWalletState(phone, false)

		c.bufMsgController.SendBufMsg(telegram.NewBufMsgTradeStart(c.port, phone))
		defer c.bufMsgController.SendBufMsg(telegram.NewBufMsgTradeStop(c.port, phone))
		<-stopC
	}()

	return stopC
}

func (c *PortController) setupTariff(tariffSMSc rxgo.Observable, phone string) error {
	// smsC := tariffSMSc.Observe()

	// TODO: Написать средство общения по CUSD в удобном виде, не только на COMMAND
	// TODO: По возможности уйти от библиотеки и перейти в ручной режим

	// TODO: получить тариф (*111*59#)
	// TODO: убрать тариф (*111*8888#)
	// TODO: убрать опцию "Все супер" (*111*249#)
	// TODO: убрать опцию "Вам звонили" (*111*38#)

	// cmd := "*111*59#"
	// response, err := c.sendATasync(c.gsm, fmt.Sprintf(`+CUSD=1,"%s"`, cmd), "+CUSD")
	// if err != nil {
	// 	return err
	// }

	// decode := utils.DecodeUCS2(response[0])
	// c.log.Print(strings.Join(response, "====="))
	// c.log.Print(decode)
	// c.log.WithFields(logrus.Fields{
	// 	"cmd":      cmd,
	// 	"response": decode,
	// })

	return nil
	// for {
	// c.sendATasync(c.gsm, `+CUSD=1,"*111*59#"`, "+CUSD") // TODO: handle error

	// 	for {
	// 		select {
	// 		case <-time.After(time.Minute):
	// 			break
	// 		case sms := <-smsC:
	// 			msg := sms.V.(Message)

	// 			part1 := "Ваш тариф: "
	// 			part2 := " . П"

	// 			firstI := strings.Index(msg.Body, part1)
	// 			secondI := strings.Index(msg.Body, part2)

	// 			if firstI > -1 && firstI < len(msg.Body) {
	// 				if secondI > -1 && secondI < len(msg.Body) {
	// 					tariff := msg.Body[firstI+len(part1) : secondI]
	// 					c.log.WithFields(logrus.Fields{
	// 						"tariff": tariff,
	// 					}).Infof("Tariff received")
	// 					if strings.TrimSpace(tariff) == "Барнаул - Супер МТС" {
	// 						return nil
	// 					} else {
	// 						break
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// c.bufMsgController.SendBufMsg(telegram.NewBufMsgChangeTariff(c.port, phone))
	// c.log.Infof("Changing tariff to SUPER - MTC")

	// c.sendATasync(c.gsm, `+CUSD=1,"*111*8888*1#"`, `+CUSD`)

	// for {
	// 	select {
	// 	case <-time.After(time.Minute):
	// 		return fmt.Errorf("timeout")
	// 	case sms := <-smsC:
	// 		msg := sms.V.(Message)

	// 		if strings.Contains(strings.ToLower(msg.Body), "тариф изменен") {

	// 		}
	// 	}
	// }

	// for {
	// 	c.sendATasync(c.gsm, `+CUSD=1,"*111*59#"`, "+CUSD") // TODO: handle error

	// 	for {
	// 		select {
	// 		case <-time.After(time.Minute):
	// 			break
	// 		case sms := <-smsC:
	// 			msg := sms.V.(Message)

	// 			part1 := "Ваш тариф: "
	// 			part2 := " . П"

	// 			firstI := strings.Index(msg.Body, part1)
	// 			secondI := strings.Index(msg.Body, part2)

	// 			if firstI > -1 && firstI < len(msg.Body) {
	// 				if secondI > -1 && secondI < len(msg.Body) {
	// 					tariff := msg.Body[firstI+len(part1) : secondI]
	// 					if strings.TrimSpace(tariff) != "Барнаул - Супер МТС" {
	// 						c.bufMsgController.SendBufMsg(telegram.NewBufMsgChangeTariff(c.port, phone))
	// 						c.log.WithFields(logrus.Fields{
	// 							"tariff":  tariff,
	// 							"msg":     msg.Body,
	// 							"firstI":  firstI,
	// 							"secondI": secondI,
	// 						}).Infof("Changing tariff to SUPER - MTC")
	// 						c.sendATasync(c.gsm, `+CUSD=1,"*111*8888*1#"`, `+CUSD`)
	// 						// TODO: Change tariff to SUPER - MTC
	// 						for {
	// 							select {
	// 							case <-time.After(time.Minute * 5):
	// 								break
	// 							case sms = <-smsC:
	// 								msg := sms.V.(Message)

	// 								if strings.Contains(strings.ToLower(msg.Body), "тариф изменен") {

	// 								}
	// 							}
	// 						}
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }
}

func (c *PortController) makeCall() {
	// TODO
}

func (c *PortController) disableCall() chan struct{} {
	stopC := make(chan struct{})

	go func() {
		c.log.Infof(log.InfoColor, "SMS call will dismiss")
		defer c.log.Infof(log.InfoColor, "SMS call dismiss disabled")

		c.gsm.AddIndication("RING", func(s []string) {
			c.gsm.Command("H") // TODO: handle error
		})
		defer c.gsm.CancelIndication("RING")

		<-stopC
	}()

	return stopC
}

func (c *PortController) listenSMS() (rxgo.Observable, rxgo.Observable, rxgo.Observable, chan struct{}) {
	paySMSc, numSMSc, tariffSMSc, stopSMS := make(chan rxgo.Item), make(chan rxgo.Item), make(chan rxgo.Item), make(chan struct{})

	go func() {
		c.log.Infof(log.InfoColor, "SMS listener activated")
		defer c.log.Infof(log.InfoColor, "SMS listener deactivated")

		// TODO: SMS обрывается
		c.gsm.AddIndication("+CMGL", func(s []string) {
			id := s[0][strings.Index(s[0], ":")+2 : strings.Index(s[0], ",")]

			bs, err := hex.DecodeString(s[1])
			if err != nil {
				c.log.Errorf(log.ErrorColor, "CMGL hex decode error: %v", err)
				return
			}

			rawMsg := new(sms.Message)
			_, err = rawMsg.ReadFrom(bs)
			if err != nil {
				c.log.Errorf(log.ErrorColor, "CMGL pdu read error: %v", err)
				return
			}

			msg := NewMessage(id, string(rawMsg.Address), rawMsg.Text)

			if strings.Contains(strings.ToLower(msg.Body), "ваш номер") {
				numSMSc <- rxgo.Of(msg)
			} else if strings.Contains(strings.ToLower(msg.Body), "поступил платеж") {
				paySMSc <- rxgo.Of(msg)
			} else if strings.Contains(strings.ToLower(msg.Body), "ваш тариф") {
				tariffSMSc <- rxgo.Of(msg)
			}
			c.log.WithFields(logrus.Fields{
				"id":   msg.ID,
				"body": msg.Body,
			}).Infof(log.InfoColor, "Deleting message")
			c.gsm.Command(fmt.Sprintf("+CMGD=%s", msg.ID))
		}, at.WithTrailingLine) // TODO: handle error
		defer c.gsm.CancelIndication("+CMGL")

		for {
			select {
			case <-stopSMS:
				return
			case <-time.After(time.Second * 10):
				// прочитать все смс
				c.gsm.Command(`+CMGL=4`) // TODO: handle error
			}
		}
	}()

	return rxgo.FromEventSource(paySMSc), rxgo.FromEventSource(numSMSc), rxgo.FromEventSource(tariffSMSc), stopSMS
}

func (c *PortController) logSMS(obs rxgo.Observable) chan struct{} {
	stopC := make(chan struct{})

	go func() {
		c.log.Infof(log.InfoColor, "SMS logger activated")
		defer c.log.Infof(log.InfoColor, "SMS logger deactivated")

		ch := obs.Observe()
		for {
			select {
			case u := <-ch:
				msg := u.V.(Message)
				c.log.WithFields(logrus.Fields{
					"from": msg.Phone,
					"text": msg.Body,
				}).Infof(log.InfoColor, "New SMS received")
			case <-stopC:
				return
			}
		}
	}()

	return stopC
}

func (c *PortController) redirectSMS(obs rxgo.Observable) chan struct{} {
	stopC := make(chan struct{})

	go func() {
		c.log.Infof(log.InfoColor, "SMS redirect activated")
		defer c.log.Infof(log.InfoColor, "SMS redirect deactivated")

		ch := obs.Observe()
		for {
			select {
			case u := <-ch:
				msg := u.V.(Message)
				email.SendEmail(c.cfg, fmt.Sprintf("[%s] New SMS", c.port), strings.Join([]string{
					fmt.Sprintf("Port: %s", c.port),
					fmt.Sprintf("From: %s", msg.Phone),
					fmt.Sprintf("Body: %s", msg.Body),
				}, "\n"))
			case <-stopC:
				return
			}
		}
	}()

	return stopC
}

func (c *PortController) autopay(obs rxgo.Observable, phone string) chan struct{} {
	stopC := make(chan struct{})

	go func() {
		c.log.Infof(log.InfoColor, "SMS autopay activated")
		defer c.log.Infof(log.InfoColor, "SMS autopay deactivated")
		ch := obs.Observe()
		for {
			select {
			case u := <-ch:
				msg := u.V.(Message)
				if err := c.api.AddPayment(phone[len(phone)-4:], msg.Body); err != nil {
					c.log.WithError(err).Errorf(log.ErrorColor, "Cannot add payment")
				}
			case <-stopC:
				return
			}
		}
	}()

	return stopC
}

func (c *PortController) detectPhone(numSMSc rxgo.Observable, stopC chan struct{}) (chan string, chan struct{}) {
	smsC := numSMSc.Observe()
	forceStop := make(chan struct{})
	result := make(chan string)
	refetch := make(chan struct{})

	go func() {
		refetch <- struct{}{}
		for {
			select {
			case <-forceStop:
				return
			case <-stopC:
				return
			case <-time.After(time.Minute):
				refetch <- struct{}{}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-forceStop:
				return
			case <-stopC:
				return
			case <-refetch:
				c.sendATasync(c.gsm, `+CUSD=1,"*111*0887#"`, "+CUSD") // TODO: handle error // TODO: check and log answer
			case sms := <-smsC:
				msg := sms.V.(Message)
				if strings.Contains(strings.ToLower(msg.Body), strings.ToLower("ваш номер")) {
					re := regexp.MustCompile("[0-9]+")
					arr := re.FindAllString(msg.Body, -1)
					for i := range arr {
						if strings.HasPrefix(arr[i], "7") || strings.HasPrefix(arr[i], "+7") || strings.HasPrefix(arr[i], "8") {
							result <- arr[i]
						}
					}
				}

			}
		}
	}()

	return result, forceStop
}

func (c *PortController) autoconfig() {
	// PDU режим
	c.gsm.Command("+CMGF=0") // TODO: handle error

	// ввод пин кода
	// c.gsm.Command(`+CPIN="0000"`) // TODO: handle error

	// расширенный вывод информации об ошибке
	c.gsm.Command(`+CMEE=2`) // TODO: handle error

	// память SIM карты для хранения СМС
	c.gsm.Command(`+CPMS="SM","SM","SM"`) // TODO: handle error

	// разрешить индикацию СМС без содержимого
	c.gsm.Command(`+CNMI=2,1`) // TODO: handle error

	// текущий оператор
	c.gsm.Command("+COPS?") // TODO: handle error

	waitSeconds := 15
	time.Sleep(time.Second * time.Duration(waitSeconds))

	// статус работы
	i, e := c.gsm.Command("+CREG?") // TODO: handle error
	if e != nil {
		waitSeconds := 60
		time.Sleep(time.Second * time.Duration(waitSeconds))
	} else {
		cregStatus := strings.Split(i[0], ",")[1]
		if cregStatus == "2" {
			waitSeconds := 80
			time.Sleep(time.Second * time.Duration(waitSeconds))
		} else if cregStatus == "0" || cregStatus == "3" || cregStatus == "4" {
			waitSeconds := 40
			time.Sleep(time.Second * time.Duration(waitSeconds))
			c.gsm.Command("+CFUN=1,1")
			time.Sleep(time.Second * 60)
		}
	}
}

func (c *PortController) sendATasync(gsm *at.AT, cmd, prefix string) ([]string, error) {
	res := make(chan []string)

	if err := gsm.AddIndication(prefix, func(s []string) {
		res <- s
	}); err != nil {
		return nil, err
	}
	defer gsm.CancelIndication(prefix)

	if _, err := gsm.Command(cmd); err != nil {
		return nil, err
	}

	select {
	case v := <-res:
		return v, nil
	case <-time.After(time.Minute):
		return nil, fmt.Errorf("timeout exception")
	}

}
