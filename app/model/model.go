package model

type ResponseData struct {
	Code         int
	Data         string
	OriginalData interface{}
}

type Sender struct {
	Age      int    `json:"age"`
	NickName string `json:"nickname"`
	Sex      string `json:"sex"`
	UserId   int64  `json:"user_id"`
}

type FriendPrivateMessage struct {
	Font        int    `json:"font"`
	Message     string `json:"message"`
	MessageId   int64  `json:"message_id"`
	MessageType string `json:"message_type"`
	PostType    string `json:"post_type"`
	RawMessage  string `json:"raw_message"`
	SelfId      int64  `json:"self_id"`
	Sender      Sender `json:"sender"`
	SubType     string `json:"sub_type"`
	TargetId    int64  `json:"target_id"`
	Time        int64  `json:"time"`
	UserId      int64  `json:"user_id"`
}
