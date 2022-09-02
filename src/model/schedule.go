package model

import (
	"bug-carrot/src/param"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const scheduleCollectionName = "schedule"

//userCollection returns the mongo.collection of users
func (m *model) scheduleCollection() *mongo.Collection {
	return m.database.Collection(scheduleCollectionName)
}

type ScheduleInterface interface {
	AddSchedule(schedule param.Schedule) error
	GetScheduleById(id string) (param.Schedule, error)
	DeleteScheduleById(id string) (param.Schedule, error)
	UpdateScheduleById(id string, schedule param.Schedule) (param.Schedule, error)
	GetScheduleAllFromNow() ([]param.Schedule, error)
	GetScheduleByTitleFromNow(title string, page int64, limit int64) ([]param.Schedule, error)
	GetScheduleCount() (int64, error)
}

func (m *model) AddSchedule(schedule param.Schedule) error {
	filter := bson.M{"schedule_id": schedule.ScheduleId}
	update := bson.M{"$setOnInsert": schedule}

	boolTrue := true
	opt := options.UpdateOptions{
		Upsert: &boolTrue,
	}

	res, err := m.scheduleCollection().UpdateOne(m.context, filter, update, &opt)
	if err != nil {
		return err
	}
	if res.UpsertedCount == 0 {
		return errors.New("not create")
	}

	return nil
}

func (m *model) GetScheduleById(id string) (param.Schedule, error) {
	filter := bson.M{"schedule_id": id}

	var schedule param.Schedule
	err := m.scheduleCollection().FindOne(m.context, filter).Decode(&schedule)
	if err != nil {
		return param.Schedule{}, err
	}

	return schedule, nil
}

func (m *model) DeleteScheduleById(id string) (param.Schedule, error) {
	filter := bson.M{"schedule_id": id}

	var schedule param.Schedule
	replacedSchedule, err := m.GetScheduleById(id)
	if err != nil {
		return param.Schedule{}, err
	}

	replacedSchedule.ExistFlag = false
	err = m.scheduleCollection().FindOneAndReplace(m.context, filter, replacedSchedule).Decode(&schedule)
	if err != nil {
		return param.Schedule{}, err
	}

	return replacedSchedule, err
}

func (m *model) UpdateScheduleById(id string, schedule param.Schedule) (param.Schedule, error) {
	filter := bson.M{"schedule_id": id}

	replacedSchedule, err := m.GetScheduleById(id)
	if err != nil {
		return param.Schedule{}, err
	}

	schedule.ScheduleId = replacedSchedule.ScheduleId
	err = m.scheduleCollection().FindOneAndReplace(m.context, filter, schedule).Decode(&replacedSchedule)
	if err != nil {
		return param.Schedule{}, err
	}

	return replacedSchedule, err
}

func (m *model) GetScheduleByTitleFromNow(title string, page int64, limit int64) ([]param.Schedule, error) {
	filter := bson.M{"title": bson.M{"$regex": primitive.Regex{Pattern: title, Options: "i"}},
		"date":       bson.M{"$gt": time.Now()},
		"exist_flag": true,
	}

	sort := bson.M{"date": 1}

	cursor, err := m.scheduleCollection().Find(m.context, filter, options.Find().SetSort(sort))
	if err != nil {
		return nil, err
	}

	var schedules []param.Schedule
	if err = cursor.All(m.context, &schedules); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (m *model) GetScheduleAllFromNow() ([]param.Schedule, error) {
	filter := bson.M{"date": bson.M{"$gt": time.Now()}, "exist_flag": true}

	sort := bson.M{"date": 1}
	cursor, err := m.scheduleCollection().Find(m.context, filter, options.Find().SetSort(sort))
	if err != nil {
		return nil, err
	}

	var schedules []param.Schedule
	if err = cursor.All(m.context, &schedules); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (m *model) GetScheduleCount() (int64, error) {
	cnt, err := m.scheduleCollection().CountDocuments(m.context, bson.M{})
	if err != nil {
		return cnt, err
	}

	return cnt, nil
}
