package chat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessage_String(t *testing.T) {
	msg := Compose(
		Entry(RoleUser, "Hello"),
		Entry(RoleAI, "Hi"),
		Entry(RoleSystem, "Hello"),
		Entry(RoleAI, "Hi"),
		Entry(RoleAI, "Hi"),
	)

	str := msg.String()

	require.Equal(t, str, "User: Hello\nAI: Hi\nSystem: Hello\nAI: Hi\nAI: Hi")
}
