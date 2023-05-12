package chat

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
)

func TestMessage_String(t *testing.T) {
	msg := Compose(
		Entry(msn.RoleUser, "Hello"),
		Entry(msn.RoleAI, "Hi"),
		Entry(msn.RoleSystem, "Hello"),
		Entry(msn.RoleAI, "Hi"),
		Entry(msn.RoleAI, "Hi"),
	)

	str := msg.String()

	require.Equal(t, str, "User: Hello\nAI: Hi\nSystem: Hello\nAI: Hi\nAI: Hi")
}
