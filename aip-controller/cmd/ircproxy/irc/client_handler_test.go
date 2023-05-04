package irc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	irc2 "github.com/greenboxal/aip/aip-controller/cmd/ircproxy/irc"
)

func TestHandlerFunc(t *testing.T) {
	t.Parallel()

	hit := false
	var f irc2.HandlerFunc = func(c *irc2.Client, m *irc2.Message) {
		hit = true
	}

	f.Handle(nil, nil)
	assert.True(t, hit, "HandlerFunc doesn't work correctly as Handler")
}
