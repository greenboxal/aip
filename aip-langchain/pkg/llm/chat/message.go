package chat

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	"github.com/greenboxal/aip/aip-langchain/pkg/tokenizers"
)

func Compose(entries ...MessageEntry) Message {
	return Message{
		Entries: entries,
	}
}

func Merge(messages ...Message) Message {
	entries := make([]MessageEntry, 0, len(messages))

	for _, msg := range messages {
		entries = append(entries, msg.Entries...)
	}

	return Message{
		Entries: entries,
	}
}

func Entry(role msn.Role, content string) MessageEntry {
	return MessageEntry{
		Role: role,
		Text: content,
	}
}

func Append(msg Message, entries ...MessageEntry) Message {
	newEntries := make([]MessageEntry, 0, len(msg.Entries)+len(entries))
	newEntries = append(newEntries, msg.Entries...)
	newEntries = append(newEntries, entries...)

	return Message{
		Entries: newEntries,
	}
}

type Message struct {
	Entries []MessageEntry
}

func (m Message) CalculateTokenCount(tokenizer tokenizers.BasicTokenizer) (int, error) {
	total := 0

	for _, entry := range m.Entries {
		count, err := entry.CalculateTokenCount(tokenizer)

		if err != nil {
			return 0, err
		}

		total += count
	}

	return total, nil
}

func (m Message) AsText() string {
	return m.String()
}

func (m Message) String() string {
	entries := lo.Map(m.Entries, func(entry MessageEntry, _index int) string {
		return entry.String()
	})

	return strings.Join(entries, "\n")
}

type MessageEntry msn.Message

func (m MessageEntry) CalculateTokenCount(tokenizer tokenizers.BasicTokenizer) (int, error) {
	return tokenizer.Count(m.Text)
}

func (m MessageEntry) String() string {
	if m.Fn != "" {
		return fmt.Sprintf("%s (%s): %s(%s)", m.Role, m.Name, m.Fn, m.FnArgs)
	}
	if m.Name != "" {
		return fmt.Sprintf("%s (%s): %s", m.Role, m.Name, m.Text)
	}

	return fmt.Sprintf("%s: %s", m.Role, m.Text)
}
