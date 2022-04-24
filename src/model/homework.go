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
	DeleteHomework(homework param.Homework) error
	ClearAllHomework() error
	GetHomeworkByTimeRange(timeL time.Time, timeR time.Time) ([]param.Homework, error)
	//GetHomeworkBySubject(context string) ([]param.Homework, error)
}

func (m *model) AddHomework(homework param.Homework) error {
	homework.CreateTime = time.Now()

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

func (m *model) DeleteHomework(homework param.Homework) error {
	//filter := bson.M{"subject": homework.Subject, "context": homework.Context}
	//update := bson.M{"$set": homework}
	//_, err := m.homeworkCollection().UpdateOne(m.context, filter, update)
	//if err != nil {
	//	return err
	//}

	filter := bson.M{"subject": homework.Subject, "context": homework.Context}
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

func (m *model) GetHomeworkByTimeRange(timeL time.Time, timeR time.Time) ([]param.Homework, error) {
	filter := bson.M{"create_time": bson.M{
		"$gt": timeL,
		"$lt": timeR,
	}}

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
