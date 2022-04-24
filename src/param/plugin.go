package param

type PluginInterface interface { // 群聊消息插件
	GetPluginName() string

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
	PluginName string

	FlagCanTime           bool
	FlagCanMatchedGroup   bool
	FlagCanMatchedPrivate bool
	FlagCanListen         bool
}
