package model

import (
	"bug-carrot/param"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const familyCollectionName = "family"

//userCollection returns the mongo.collection of users
func (m *model) familyCollection() *mongo.Collection {
	return m.database.Collection(familyCollectionName)
}

type FamilyInterface interface {
	AddOneFamilyMember(member param.FamilyMember) error
	DeleteOneFamilyMemberByStudentID(studentID string) error
	GetAllFamilyMember() ([]param.FamilyMember, error)
	//GetOneFamilyMemberByParam(member param.FamilyMember) error
	GetOneFamilyMemberByStudentID(studentID string) (param.FamilyMember, error)
	GetOneFamilyMemberByQQ(qq int64) (param.FamilyMember, error)
	GetOneFamilyMemberByName(name string) (param.FamilyMember, error)
	UpdateFamilyMemberByStudentID(studentID string, member param.FamilyMember) (param.FamilyMember, error)
}

func (m *model) AddOneFamilyMember(member param.FamilyMember) error {
	filter := bson.M{"student_id": member.StudentID}
	update := bson.M{"$setOnInsert": member}

	boolTrue := true
	opt := options.UpdateOptions{
		Upsert: &boolTrue,
	}

	res, err := m.familyCollection().UpdateOne(m.context, filter, update, &opt)
	if err != nil {
		return err
	}
	if res.UpsertedCount == 0 {
		return errors.New("not create")
	}

	return nil
}

func (m *model) DeleteOneFamilyMemberByStudentID(studentID string) error {
	filter := bson.M{"student_id": studentID}
	res, err := m.familyCollection().DeleteOne(m.context, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("not found")
	}

	return nil
}

func (m *model) GetAllFamilyMember() ([]param.FamilyMember, error) {
	cursor, err := m.familyCollection().Find(m.context, bson.M{})
	if err != nil {
		return nil, err
	}

	var familyMembers []param.FamilyMember
	if err = cursor.All(m.context, &familyMembers); err != nil {
		return nil, err
	}

	return familyMembers, nil
}

// GetOneFamilyMemberByParam: find one member with the param that not nil
//func (m *model) GetOneFamilyMemberByParam(member param.FamilyMember) {
//var familyMember param.FamilyMember
//v := reflect.ValueOf(member)
//for i := 0; i < v.NumField(); i++ {
//	fieldInfo := v.Type().Field(i) // a reflect.StructField
//	tag := fieldInfo.Tag           // a reflect.StructTag
//	filter := bson.M{tag.Get("bson"): bson.M{"$eq": v.Field(i)}}
//	cursor, err := m.foodCollection().FindOne(m.context, filter).Decode(&familyMember)
//	if err != nil {
//		continue
//	}
//}
//}

func (m *model) GetOneFamilyMemberByStudentID(studentID string) (param.FamilyMember, error) {
	var familyMember param.FamilyMember
	filter := bson.M{"student_id": bson.M{"$eq": studentID}}
	err := m.familyCollection().FindOne(m.context, filter).Decode(&familyMember)
	if err != nil {
		return familyMember, err
	}
	return familyMember, nil
}
func (m *model) GetOneFamilyMemberByQQ(qq int64) (param.FamilyMember, error) {
	var familyMember param.FamilyMember
	filter := bson.M{"qq": bson.M{"$eq": qq}}
	err := m.familyCollection().FindOne(m.context, filter).Decode(&familyMember)
	if err != nil {
		return familyMember, err
	}
	return familyMember, nil
}
func (m *model) GetOneFamilyMemberByName(name string) (param.FamilyMember, error) {
	var familyMember param.FamilyMember
	filter := bson.M{"name": bson.M{"$eq": name}}
	err := m.familyCollection().FindOne(m.context, filter).Decode(&familyMember)
	if err != nil {
		return familyMember, err
	}
	return familyMember, nil
}

func (m *model) UpdateFamilyMemberByStudentID(studentID string, member param.FamilyMember) (param.FamilyMember, error) {
	filter := bson.M{"student_id": studentID}

	replacedMember, err := m.GetOneFamilyMemberByStudentID(studentID)
	if err != nil {
		return param.FamilyMember{}, err
	}

	member.StudentID = replacedMember.StudentID
	err = m.familyCollection().FindOneAndReplace(m.context, filter, member).Decode(&replacedMember)
	if err != nil {
		return param.FamilyMember{}, err
	}

	return replacedMember, err
}
