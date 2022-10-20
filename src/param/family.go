package param

import "time"

type FamilyMember struct {
	StudentID string    `bson:"student_id" json:"student_id"`
	Name      string    `bson:"name" json:"name"`
	QQ        int64     `bson:"qq" json:"qq"`
	Phone     string    `bson:"phone" json:"phone"`
	Mail      string    `bson:"mail" json:"mail"`
	Address   string    `bson:"address" json:"address"`
	Birthday  time.Time `bson:"birthday" json:"birthday"`
}

type PersonWithQQ struct {
	Name string `bson:"name" json:"name"`
	QQ   int64  `bson:"qq" json:"qq"`
}
