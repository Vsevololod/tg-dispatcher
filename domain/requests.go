package domain

type DestinationType int

const (
	VideoDownload DestinationType = iota
	VideoMessageSand
)

func (s DestinationType) String() string {
	return [...]string{"VideoDownload", "VideoMessageSand"}[s]
}

type VideoDownloadReq struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type MessageSendReq struct {
	UserId int64  `json:"user_id"`
	HashId string `json:"hash_id"`
}

type MessageReq struct {
	UUID        string          `json:"UUID"`
	Destination DestinationType `json:"destination"`
	Message     interface{}     `json:"message"`
}
