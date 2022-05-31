package model

import (
	"bug-carrot/param"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const keyWordCollectionName = "keyword"

func (m *model) KeyWordCollection() *mongo.Collection {
	return m.database.Collection(keyWordCollectionName)
}

type KeyWordInterface interface {
	AddKeyWord(KeyWord param.KeyWord) error
	DeleteKeyWord(keyWord string) error
	GetKeyWord(keyWord string) (param.KeyWord, error)
}

func (m *model) AddKeyWord(KeyWord param.KeyWord) error {
	filter := bson.M{"keyword": KeyWord.KeyWord}
	update := bson.M{"$setOnInsert": KeyWord}

	boolTrue := true
	opt := options.UpdateOptions{
		Upsert: &boolTrue,
	}

	res, err := m.KeyWordCollection().UpdateOne(m.context, filter, update, &opt)
	if err != nil {
		return err
	}
	if res.UpsertedCount == 0 {
		return errors.New("not create")
	}

	return nil
}

func (m *model) DeleteKeyWord(keyWord string) error {
	filter := bson.M{"keyword": keyWord}
	res, err := m.KeyWordCollection().DeleteOne(m.context, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("not found")
	}

	return nil
}

func (m *model) GetKeyWord(keyWord string) (param.KeyWord, error) {
	filter := bson.M{"keyword": keyWord}

	var kw param.KeyWord
	err := m.KeyWordCollection().FindOne(m.context, filter).Decode(&kw)
	if err != nil {
		return param.KeyWord{}, err
	}

	return kw, nil
}
