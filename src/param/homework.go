package param

import "time"

type Homework struct {
	Subject string       `bson:"subject"`
	Context string       `bson:"context"`
	Weekday time.Weekday `bson:"weekday"`
}
