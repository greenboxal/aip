package collective

type Message struct {
	ID        string `json:"id"`
	ThreadID  string `json:"thread_id"`
	ReplyToID string `json:"reply_to_id"`

	Channel string `json:"channel"`
	From    string `json:"from"`
	Text    string `json:"text"`
}
