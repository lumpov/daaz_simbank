package telegram

import (
	"fmt"
	"strings"
)

// portTariffInfo to send in message
type portTariffInfo struct {
	port          string
	currentTariff string
}

// BufMsgChangeTariff message on port tariff changing
type BufMsgChangeTariff struct {
	Ports []portTariffInfo
}

// NewBufMsgChangeTariff instance
func NewBufMsgChangeTariff(port, currentTariff string) BufMsgChangeTariff {
	return BufMsgChangeTariff{
		Ports: []portTariffInfo{{
			port:          port,
			currentTariff: currentTariff,
		}},
	}
}

// IsEqual check type of msg
func (m BufMsgChangeTariff) IsEqual(msg interface{}) bool {
	switch msg.(type) {
	case BufMsgChangeTariff:
		return true
	default:
		return false
	}
}

// Combine all data from old msg to new
func (m BufMsgChangeTariff) Combine(msg interface{}) BufMsg {
	return BufMsgChangeTariff{
		Ports: append(m.Ports, msg.(BufMsgChangeTariff).Ports...),
	}
}

// Build all data to string
func (m BufMsgChangeTariff) Build() string {
	str := []string{
		fmt.Sprintf(`<b>Происходит активация тарифа "СУПЕР - МТС" у <code>%d</code> портов</b>`, len(m.Ports)),
	}

	for i := range m.Ports {
		str = append(str, fmt.Sprintf("<b>- <code>%s</code> <code>(%s)</code></b>", m.Ports[i].port, m.Ports[i].currentTariff))
	}

	str = append(str, "\n", "<b>Внимание!</b> Порядок доставки сообщений может быть искажен.")

	return strings.Join(str, "\n")
}
