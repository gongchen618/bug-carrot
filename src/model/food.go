package model

import (
	"bug-carrot/param"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const foodCollectionName = "food"

//userCollection returns the mongo.collection of users
func (m *model) foodCollection() *mongo.Collection {
	return m.database.Collection(foodCollectionName)
}

type FoodInterface interface {
	AddFood(food param.Food) error
	DeleteFood(food param.Food) error
	GetFoodAll() ([]param.Food, error)
	GetFoodByAddress(address string) ([]param.Food, error)
	GetFoodByName(name string) ([]param.Food, error)
}

func (m *model) AddFood(food param.Food) error {
	filter := bson.M{"name": food.Name}
	update := bson.M{"$setOnInsert": food}

	boolTrue := true
	opt := options.UpdateOptions{
		Upsert: &boolTrue,
	}

	res, err := m.foodCollection().UpdateOne(m.context, filter, update, &opt)
	if err != nil {
		return err
	}
	if res.UpsertedCount == 0 {
		return errors.New("not create")
	}

	return nil
}

func (m *model) DeleteFood(food param.Food) error {
	filter := bson.M{"name": food.Name}
	res, err := m.foodCollection().DeleteOne(m.context, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("not found")
	}

	return nil
}

func (m *model) GetFoodAll() ([]param.Food, error) {
	cursor, err := m.foodCollection().Find(m.context, bson.M{})
	if err != nil {
		return nil, err
	}

	var foods []param.Food
	if err = cursor.All(m.context, &foods); err != nil {
		return nil, err
	}

	return foods, nil
}

func (m *model) GetFoodByAddress(address string) ([]param.Food, error) {
	filter := bson.M{"address": bson.M{"$eq": address}}
	cursor, err := m.foodCollection().Find(m.context, filter)
	if err != nil {
		return nil, err
	}

	var foods []param.Food
	if err = cursor.All(m.context, &foods); err != nil {
		return nil, err
	}

	return foods, nil
}

func (m *model) GetFoodByName(name string) ([]param.Food, error) {
	filter := bson.M{"name": bson.M{"$eq": name}}
	cursor, err := m.foodCollection().Find(m.context, filter)
	if err != nil {
		return nil, err
	}

	var foods []param.Food
	if err = cursor.All(m.context, &foods); err != nil {
		return nil, err
	}

	return foods, nil
}
