package http

import "encoding/json"

type MessageEnvelope struct {
	ID      int             `json:"id"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func WrapMessage(message WebsocketMessage) (MessageEnvelope, error) {
	payload, err := json.Marshal(message)
	if err != nil {
		return MessageEnvelope{}, err
	}

	return MessageEnvelope{
		Type:    message.Type(),
		Payload: payload,
	}, nil
}
func WrapReply(id int, message WebsocketMessage) (MessageEnvelope, error) {
	payload, err := json.Marshal(message)
	if err != nil {
		return MessageEnvelope{}, err
	}

	return MessageEnvelope{
		ID:      id,
		Type:    message.Type(),
		Payload: payload,
	}, nil
}

type WebsocketMessage interface {
	Type() string
}

type WelcomeMessage struct {
	Message string `json:"message"`
	Date    string `json:"date"`
	Version string `json:"version"`
}

func (m WelcomeMessage) Type() string {
	return "welcome"
}
