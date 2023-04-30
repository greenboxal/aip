package irc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/greenboxal/aip/cmd/ircproxy/irc"
)

func TestHandlerFunc(t *testing.T) {
	t.Parallel()

	hit := false
	var f irc.HandlerFunc = func(c *irc.Client, m *irc.Message) {
		hit = true
	}

	f.Handle(nil, nil)
	assert.True(t, hit, "HandlerFunc doesn't work correctly as Handler")
}
