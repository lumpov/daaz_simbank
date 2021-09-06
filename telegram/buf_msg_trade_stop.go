package telegram

import (
	"fmt"
	"strings"
)

// BufMsgTradeStop message on new trade stopped
type BufMsgTradeStop struct {
	Ports []portInfo
}

// NewBufMsgTradeStop instance
func NewBufMsgTradeStop(port, phone string) BufMsgTradeStop {
	return BufMsgTradeStop{
		Ports: []portInfo{{
			phone: phone,
			port:  port,
		}},
	}
}

// IsEqual check type of msg
func (m BufMsgTradeStop) IsEqual(msg interface{}) bool {
	switch msg.(type) {
	case BufMsgTradeStop:
		return true
	default:
		return false
	}
}

// Combine all data from old msg to new
func (m BufMsgTradeStop) Combine(msg interface{}) BufMsg {
	return BufMsgTradeStop{
		Ports: append(m.Ports, msg.(BufMsgTradeStop).Ports...),
	}
}

// Build all data to string
func (m BufMsgTradeStop) Build() string {
	str := []string{
		fmt.Sprintf("<b>Отключено портов: <code>%d</code></b>", len(m.Ports)),
	}

	for i := range m.Ports {
		str = append(str, fmt.Sprintf("<b>- <code>%s</code> <code>(%s)</code></b>", m.Ports[i].port, m.Ports[i].phone))
	}

	str = append(str, "\n", "<b>Внимание!</b> Порядок доставки сообщений может быть искажен.")

	return strings.Join(str, "\n")
}
