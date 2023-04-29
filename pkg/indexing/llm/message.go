package llm

type Message struct {
	Entries []MessageEntry
}

func (m *Message) CalculateTokenCount(tokenizer BasicTokenizer) (int, error) {
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

type MessageEntry struct {
	Role    string
	Content string
}

func (m *MessageEntry) CalculateTokenCount(tokenizer BasicTokenizer) (int, error) {
	return tokenizer.Count(m.Content)
}

type MessageComposer struct {
	Tokenizer BasicTokenizer

	Constraints
	Statistics

	entries []MessageEntry
}

func (mc *MessageComposer) Append(entry MessageEntry) error {
	count, err := entry.CalculateTokenCount(mc.Tokenizer)

	if err != nil {
		return err
	}

	mc.entries = append(mc.entries, entry)
	mc.TokenCount += count

	return nil
}

func (mc *MessageComposer) RemainingTokens() int {
	return mc.MaxTokens - mc.TokenCount
}

func (mc *MessageComposer) Validate() error {
	return nil
}

func (mc *MessageComposer) Build() (Message, error) {
	if err := mc.Validate(); err != nil {
		return Message{}, err
	}

	m := Message{
		Entries: mc.entries,
	}

	return m, nil
}
