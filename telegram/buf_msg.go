package telegram

import (
	"sync"
	"time"
)

// BufMsg interface for message to send buffered
type BufMsg interface {
	Build() string
	Combine(msg interface{}) BufMsg
	IsEqual(msg interface{}) bool
}

// BufMsgController struct to control buf msgs
type BufMsgController struct {
	msgs []BufMsg
	m    sync.Mutex
	tg   *Bot
}

// NewBufMsgController instance
func NewBufMsgController(tg *Bot) *BufMsgController {
	return &BufMsgController{
		msgs: []BufMsg{},
		tg:   tg,
	}
}

// Start thread of sending
func (c *BufMsgController) Start() {
	go func() {
		for {
			time.Sleep(time.Second * 6)
			c.m.Lock()
			for i := range c.msgs {
				c.tg.Send(c.msgs[i].Build())
			}

			c.msgs = []BufMsg{}
			c.m.Unlock()
		}
	}()
}

// SendBufMsg add message to queue
func (c *BufMsgController) SendBufMsg(msg BufMsg) {
	c.m.Lock()
	defer c.m.Unlock()

	for i := range c.msgs {
		if c.msgs[i].IsEqual(msg) {
			c.msgs[i] = c.msgs[i].Combine(msg)
			return
		}
	}

	c.msgs = append(c.msgs, msg)
}
