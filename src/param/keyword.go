package param

type KeyWord struct {
	KeyWord string `bson:"keyword"`
	Content string `bson:"content"`
	Author  int64  `bson:"author"`
}
