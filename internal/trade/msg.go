package trade

// Message from SIM
type Message struct {
	ID    string
	Phone string
	Body  string
}

// NewMessage instance
func NewMessage(id, phone, body string) Message {
	return Message{ID: id, Phone: phone, Body: body}
}
