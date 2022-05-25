package param

import (
	"time"
)

type Schedule struct {
	Date         time.Time `bson:"date"`
	Title        string    `bson:"title"`
	Description  string    `bson:"description"`
	ScheduleId   string    `bson:"schedule_id"`
	ExistFlag    bool      `bson:"exist_flag"`
	Alarm24hFlag bool      `bson:"alarm_24h_flag"`
	Alarm1hFlag  bool      `bson:"alarm_1h_flag"`
}
