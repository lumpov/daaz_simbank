package telegram

import (
	"fmt"
	"strings"
)

// PortInfo to send in message
type portInfo struct {
	port  string
	phone string
}

// BufMsgTradeStart message on new trade started
type BufMsgTradeStart struct {
	Ports []portInfo
}

// NewBufMsgTradeStart instance
func NewBufMsgTradeStart(port, phone string) BufMsgTradeStart {
	return BufMsgTradeStart{
		Ports: []portInfo{{
			port:  port,
			phone: phone,
		}},
	}
}

// IsEqual check type of msg
func (m BufMsgTradeStart) IsEqual(msg interface{}) bool {
	switch msg.(type) {
	case BufMsgTradeStart:
		return true
	default:
		return false
	}
}

// Combine all data from old msg to new
func (m BufMsgTradeStart) Combine(msg interface{}) BufMsg {
	return BufMsgTradeStart{
		Ports: append(m.Ports, msg.(BufMsgTradeStart).Ports...),
	}
}

// Build all data to string
func (m BufMsgTradeStart) Build() string {
	str := []string{
		fmt.Sprintf("<b>Запущено портов: <code>%d</code></b>", len(m.Ports)),
	}

	for i := range m.Ports {
		str = append(str, fmt.Sprintf("<b>- <code>%s</code> <code>(%s)</code></b>", m.Ports[i].port, m.Ports[i].phone))
	}

	str = append(str, "\n", "<b>Внимание!</b> Порядок доставки сообщений может быть искажен.")

	return strings.Join(str, "\n")
}
