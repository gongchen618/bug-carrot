package param

type PluginInterface interface { // 群聊消息插件
	GetPluginName() string
	GetPluginAuthor() string

	CanTime() bool
	IsTime() bool
	DoTime() error

	CanMatchedGroup() bool
	IsMatchedGroup(msg GroupMessage) bool
	DoMatchedGroup(msg GroupMessage) error

	CanMatchedPrivate() bool
	IsMatchedPrivate(msg PrivateMessage) bool
	DoMatchedPrivate(msg PrivateMessage) error

	CanListen() bool
	Listen(msg GroupMessage)

	NeedDatabase() bool
	DoIgnoreRiskControl() bool
	Close()
}

type GroupMessage struct {
	RequestGroupMessage
	WordsMap
}

type PrivateMessage struct {
	RequestPrivateMessage
	WordsMap
}

type PluginIndex struct {
	PluginName   string
	PluginAuthor string

	FlagCanTime           bool
	FlagCanMatchedGroup   bool
	FlagCanMatchedPrivate bool
	FlagCanListen         bool
	FlagUseDatabase       bool
	FlagIgnoreRiskControl bool
}
