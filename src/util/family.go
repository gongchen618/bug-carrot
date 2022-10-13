package util

import (
	"bug-carrot/model"
	"bug-carrot/param"
	"errors"
)

func BuildFamilyMemberListFromNameList(name []string) []param.FamilyMember {
	m := model.GetModel()
	defer m.Close()

	members, err := m.GetAllFamilyMember()
	if err != nil {
		ErrorPrint(errors.New("get all family member failed"), nil, "")
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
