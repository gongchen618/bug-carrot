package model

import (
	"bug-carrot/param"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const homeworkCollectionName = "homework"

//userCollection returns the mongo.collection of users
func (m *model) homeworkCollection() *mongo.Collection {
	return m.database.Collection(homeworkCollectionName)
}

type HomeworkInterface interface {
	AddHomework(homework param.Homework) error
	DeleteHomework(subject string, context string) error
	ClearAllHomework() error
	GetHomeworkFromNow() ([]param.Homework, error)
	GetHomeWorkByWeekDay(weekday time.Weekday) ([]param.Homework, error)
}

func (m *model) AddHomework(homework param.Homework) error {
	filter := bson.M{"context": homework.Context}
	update := bson.M{"$setOnInsert": homework}

	boolTrue := true
	opt := options.UpdateOptions{
		Upsert: &boolTrue,
	}

	res, err := m.homeworkCollection().UpdateOne(m.context, filter, update, &opt)
	if err != nil {
		return err
	}
	if res.UpsertedCount == 0 {
		return errors.New("not create")
	}

	return nil
}

func (m *model) DeleteHomework(subject string, context string) error {
	filter := bson.M{"subject": subject, "context": context}
	res, err := m.homeworkCollection().DeleteOne(m.context, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("not found")
	}

	return nil
}

func (m *model) ClearAllHomework() error {
	filter := bson.M{}
	res, err := m.homeworkCollection().DeleteMany(m.context, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("already clear")
	}

	return nil
}

func (m *model) GetHomeworkFromNow() ([]param.Homework, error) {
	cursor, err := m.homeworkCollection().Find(m.context, bson.M{})
	if err != nil {
		return nil, err
	}

	var homeworks []param.Homework
	if err = cursor.All(m.context, &homeworks); err != nil {
		return nil, err
	}

	return homeworks, nil
}

func (m *model) GetHomeWorkByWeekDay(weekday time.Weekday) ([]param.Homework, error) {
	filter := bson.M{"weekday": weekday}

	cursor, err := m.homeworkCollection().Find(m.context, filter)
	if err != nil {
		return nil, err
	}

	var homeworks []param.Homework
	if err = cursor.All(m.context, &homeworks); err != nil {
		return nil, err
	}

	return homeworks, nil
}
