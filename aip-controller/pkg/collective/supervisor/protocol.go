package supervisor

import (
	"encoding/json"

	"github.com/greenboxal/aip/aip-controller/pkg/collective"
)

func DecodeMessage(data []byte) (collective.Message, error) {
	var msg collective.Message

	if err := json.Unmarshal(data, &msg); err != nil {
		return collective.Message{}, err
	}

	return msg, nil
}

func EncodeMessage(msg collective.Message) ([]byte, error) {
	data, err := json.Marshal(msg)

	if err != nil {
		return nil, err
	}

	return data, nil
}
