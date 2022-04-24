package param

import "time"

type Homework struct {
	Subject    string    `bson:"subject"`
	Context    string    `bson:"context"`
	CreateTime time.Time `bson:"create_time"`
}
