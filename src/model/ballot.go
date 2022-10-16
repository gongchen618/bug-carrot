package model

import (
	"bug-carrot/param"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const ballotCollectionName = "ballot"

func (m *model) ballotCollection() *mongo.Collection {
	return m.database.Collection(ballotCollectionName)
}

type BallotInterface interface {
	GetAllBallot() ([]param.Ballot, error)
	GetOneBallotByTitle(title string) (param.Ballot, error)
	CreateOneBallotByTitle(title string, muster param.Muster, defaultOption string) error
	DeleteOneBallotByTitle(title string) error
	AddAnOptionToOneBallot(title string, option string) (param.Ballot, error)
	DeleteAnOptionOnOneBallot(title string, option string) (param.Ballot, error)
	UpdateOptionsOnOneBallotForMembers(title string, option string, name []string) (param.Ballot, error)
}

func (m *model) GetAllBallot() ([]param.Ballot, error) {
	cursor, err := m.ballotCollection().Find(m.context, bson.M{})
	if err != nil {
		return nil, err
	}

	var ballots []param.Ballot
	if err = cursor.All(m.context, &ballots); err != nil {
		return nil, err
	}

	return ballots, nil
}

func (m *model) GetOneBallotByTitle(title string) (param.Ballot, error) {
	var ms param.Ballot
	filter := bson.M{"title": bson.M{"$eq": title}}
	err := m.ballotCollection().FindOne(m.context, filter).Decode(&ms)
	if err != nil {
		return ms, err
	}
	return ms, nil
}

func (m *model) CreateOneBallotByTitle(title string, muster param.Muster, defaultOption string) error {
	bt := param.Ballot{
		Title: title,
	}
	for _, member := range muster.People {
		bt.TargetMember = append(bt.TargetMember, param.BallotMember{
			Info:   member,
			Option: defaultOption,
		})
	}

	filter := bson.M{"title": title}
	update := bson.M{"$setOnInsert": bt}

	boolTrue := true
	opt := options.UpdateOptions{
		Upsert: &boolTrue,
	}

	res, err := m.ballotCollection().UpdateOne(m.context, filter, update, &opt)
	if err != nil {
		return err
	}
	if res.UpsertedCount == 0 {
		return errors.New("not create")
	}

	return nil
}

func (m *model) DeleteOneBallotByTitle(title string) error {
	filter := bson.M{"title": title}
	res, err := m.ballotCollection().DeleteOne(m.context, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("not found")
	}

	return nil
}

func (m *model) AddAnOptionToOneBallot(title string, option string) (param.Ballot, error) {
	var bt param.Ballot
	filter := bson.M{"title": bson.M{"$eq": title}}
	err := m.ballotCollection().FindOne(m.context, filter).Decode(&bt)
	if err != nil {
		return bt, err
	}

	bt.OfferedOptions = append(bt.OfferedOptions, option)
	err = m.ballotCollection().FindOneAndReplace(m.context, filter, bt).Decode(&bt)
	if err != nil {
		return param.Ballot{}, err
	}

	return bt, nil
}

func (m *model) DeleteAnOptionOnOneBallot(title string, option string) (param.Ballot, error) {
	var bt param.Ballot
	filter := bson.M{"title": bson.M{"$eq": title}}
	err := m.ballotCollection().FindOne(m.context, filter).Decode(&bt)
	if err != nil {
		return bt, err
	}

	var newOptions []string
	for _, opt := range bt.OfferedOptions {
		if opt != option {
			newOptions = append(newOptions, opt)
		}
	}
	bt.OfferedOptions = newOptions
	err = m.ballotCollection().FindOneAndReplace(m.context, filter, bt).Decode(&bt)
	if err != nil {
		return param.Ballot{}, err
	}

	return bt, nil
}

func (m *model) UpdateOptionsOnOneBallotForMembers(title string, option string, name []string) (param.Ballot, error) {
	var bt param.Ballot
	filter := bson.M{"title": bson.M{"$eq": title}}
	err := m.ballotCollection().FindOne(m.context, filter).Decode(&bt)
	if err != nil {
		return bt, err
	}

	var vis map[string]bool
	for _, p := range name {
		vis[p] = true
	}

	for _, member := range bt.TargetMember {
		_, ok := vis[member.Info.Name]
		if ok {
			member.Option = option
		}
	}

	err = m.ballotCollection().FindOneAndReplace(m.context, filter, bt).Decode(&bt)
	if err != nil {
		return param.Ballot{}, err
	}

	return bt, nil
}