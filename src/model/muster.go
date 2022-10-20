package model

import (
	"bug-carrot/param"
	"bug-carrot/util"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const musterCollectionName = "muster"

func (m *model) musterCollection() *mongo.Collection {
	return m.database.Collection(musterCollectionName)
}

type MusterInterface interface {
	GetAllMuster() ([]param.Muster, error)
	GetOneMusterByTitle(title string) (param.Muster, error)
	CreateOneMusterByTitle(title string) error
	DeleteOneMusterByTitle(title string) error
	AddPersonsToOneMuster(title string, name []string) (param.Muster, error)
	DeletePersonsOnOneMuster(title string, name []string) (param.Muster, error)
}

func (m *model) GetAllMuster() ([]param.Muster, error) {
	cursor, err := m.musterCollection().Find(m.context, bson.M{})
	if err != nil {
		return nil, err
	}

	var musters []param.Muster
	if err = cursor.All(m.context, &musters); err != nil {
		return nil, err
	}

	return musters, nil
}

func (m *model) GetOneMusterByTitle(title string) (param.Muster, error) {
	var ms param.Muster
	filter := bson.M{"title": bson.M{"$eq": title}}
	err := m.musterCollection().FindOne(m.context, filter).Decode(&ms)
	if err != nil {
		return ms, err
	}
	return ms, nil
}

func (m *model) CreateOneMusterByTitle(title string) error {
	ms := param.Muster{
		Title: title,
	}

	filter := bson.M{"title": title}
	update := bson.M{"$setOnInsert": ms}

	boolTrue := true
	opt := options.UpdateOptions{
		Upsert: &boolTrue,
	}

	res, err := m.musterCollection().UpdateOne(m.context, filter, update, &opt)
	if err != nil {
		return err
	}
	if res.UpsertedCount == 0 {
		return errors.New("not create")
	}

	return nil
}

func (m *model) DeleteOneMusterByTitle(title string) error {
	filter := bson.M{"title": title}
	res, err := m.musterCollection().DeleteOne(m.context, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("not found")
	}

	return nil
}

func (m *model) AddPersonsToOneMuster(title string, name []string) (param.Muster, error) {
	var ms param.Muster
	filter := bson.M{"title": bson.M{"$eq": title}}
	err := m.musterCollection().FindOne(m.context, filter).Decode(&ms)
	if err != nil {
		return ms, err
	}

	members := buildFamilyMemberListFromNameList(name)
	var vis map[string]bool
	for _, n := range ms.People {
		vis[n.Name] = true
	}
	for _, member := range members {
		_, ok := vis[member.Name]
		if !ok {
			ms.People = append(ms.People, param.PersonWithQQ{
				Name: member.Name,
				QQ:   member.QQ,
			})
		}
	}

	err = m.musterCollection().FindOneAndReplace(m.context, filter, ms).Decode(&ms)
	if err != nil {
		return param.Muster{}, err
	}

	return ms, nil
}

func (m *model) DeletePersonsOnOneMuster(title string, name []string) (param.Muster, error) {
	var ms param.Muster
	filter := bson.M{"title": bson.M{"$eq": title}}
	err := m.musterCollection().FindOne(m.context, filter).Decode(&ms)
	if err != nil {
		return ms, err
	}

	members := buildFamilyMemberListFromNameList(name)
	var vis map[string]bool
	for _, n := range members {
		vis[n.Name] = true
	}

	var newPeople []param.PersonWithQQ
	for _, member := range ms.People {
		_, ok := vis[member.Name]
		if !ok {
			newPeople = append(newPeople, param.PersonWithQQ{
				Name: member.Name,
				QQ:   member.QQ,
			})
		}
	}

	ms.People = newPeople
	err = m.musterCollection().FindOneAndReplace(m.context, filter, ms).Decode(&ms)
	if err != nil {
		return param.Muster{}, err
	}

	return ms, nil
}

func buildFamilyMemberListFromNameList(name []string) []param.FamilyMember {
	m := GetModel()
	defer m.Close()

	members, err := m.GetAllFamilyMember()
	if err != nil {
		util.ErrorPrint(errors.New("get all family member failed"), nil, "")
		return nil
	}

	var vis map[string]bool
	for _, n := range name {
		vis[n] = true
	}

	var responseMembers []param.FamilyMember
	for _, member := range members {
		_, ok := vis[member.Name]
		if ok {
			responseMembers = append(responseMembers, member)
		}
	}
	return responseMembers
}
