package supervisor

import (
	"encoding/json"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
)

func DecodeMessage(data []byte) (msn.Message, error) {
	var msg msn.Message

	if err := json.Unmarshal(data, &msg); err != nil {
		return msn.Message{}, err
	}

	return msg, nil
}

func EncodeMessage(msg msn.Message) ([]byte, error) {
	data, err := json.Marshal(msg)

	if err != nil {
		return nil, err
	}

	return data, nil
}
