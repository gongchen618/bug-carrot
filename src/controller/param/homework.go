package param

import (
	"time"
)

type SubjectType string

//const (
//	SubjectMath         SubjectType = "微积分"
//	SubjectPhysics      SubjectType = "大学物理"
//	SubjectDiscreteMath SubjectType = "离散数学"
//	SubjectEnglish      SubjectType = "英语"
//)

type Homework struct {
	Subject    string    `bson:"subject"`
	Context    string    `bson:"context"`
	CreateTime time.Time `bson:"create_time"`
}
