package telegram

import (
	"fmt"
	"strings"
)

// BufMsgCannotOpenPort message on port opening error
type BufMsgCannotOpenPort struct {
	Ports []string
}

// NewBufMsgCannotOpenPort instance
func NewBufMsgCannotOpenPort(port string) BufMsgCannotOpenPort {
	return BufMsgCannotOpenPort{
		Ports: []string{},
	}
}

// IsEqual check type of msg
func (m BufMsgCannotOpenPort) IsEqual(msg interface{}) bool {
	switch msg.(type) {
	case BufMsgCannotOpenPort:
		return true
	default:
		return false
	}
}

// Combine all data from old msg to new
func (m BufMsgCannotOpenPort) Combine(msg interface{}) BufMsg {
	return BufMsgCannotOpenPort{
		Ports: append(m.Ports, msg.(BufMsgCannotOpenPort).Ports...),
	}
}

// Build all data to string
func (m BufMsgCannotOpenPort) Build() string {
	str := []string{
		fmt.Sprintf("<b>Не удалось открыть порты: <code>%d</code></b>", len(m.Ports)),
	}

	for i := range m.Ports {
		str = append(str, fmt.Sprintf("<b>- <code>%s</code></b>", m.Ports[i]))
	}

	str = append(str, "\n", "<b>Внимание!</b> Порядок доставки сообщений может быть искажен.")

	return strings.Join(str, "\n")
}
