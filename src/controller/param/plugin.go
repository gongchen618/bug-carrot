package param

type PluginInterface interface { // 群聊消息插件
	GetPluginName() string

	IsTime() bool
	DoTime() error

	IsMatched(msg GroupMessage) bool
	DoMatched(msg GroupMessage) error

	Listen(msg GroupMessage)
	Close()
}

type GroupMessage struct {
	RequestGroupMessage
	WordsMap map[WordSplit]bool
}
