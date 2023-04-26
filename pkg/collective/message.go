package collective

type Message struct {
	ID      string `json:"id"`
	Channel string `json:"channel"`
	From    string `json:"from"`
	Text    string `json:"text"`
}
