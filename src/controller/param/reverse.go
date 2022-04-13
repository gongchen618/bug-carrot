package param

type RequestQQFriendAdd struct {
	Time        int64  `json:"time"`
	SelfId      int64  `json:"self_id"`
	PostType    string `json:"post_type"`
	RequestType string `json:"request_type"`
	UserId      int64  `json:"user_id"`
	Comment     string `json:"comment"`
	Flag        string `json:"flag"`
}

type RequestPrivateMessage struct {
	SubType    string `json:"sub_type"`
	UserId     int64  `json:"user_id"`
	RawMessage string `json:"raw_message"`
}

type RequestGroupMessage struct {
	SubType    string    `json:"sub_type"`
	RawMessage string    `json:"raw_message"`
	UserId     int64     `json:"user_id"`
	GroupId    int64     `json:"group_id"`
	Anonymous  anonymous `json:"anonymous"`
}

type anonymous struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Flag string `json:"flag"`
}
