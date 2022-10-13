package controller

import (
	"bug-carrot/model"
	"bug-carrot/param"
	"bug-carrot/util"
	"bug-carrot/util/context"
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
)

func GetAllFamilyMembersRequestHandler(c echo.Context) error {
	m := model.GetModel()
	defer m.Close()

	members, err := m.GetAllFamilyMember()
	if err != nil {
		util.ErrorPrint(err, nil, "get all family member failed")
		return context.Error(c, http.StatusInternalServerError, "", err)
	}

	return context.Success(c, members)
}

func CreateOneFamilyMemberRequestHandler(c echo.Context) error {
	m := model.GetModel()
	defer m.Close()

	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	p := param.FamilyMember{}
	if err := json.Unmarshal(bodyBytes, &p); err != nil {
		util.ErrorPrint(err, nil, "unmarshal failed on creating family member")
		return context.Error(c, http.StatusBadRequest, "unmarshal failed", err)
	}

	if p.StudentID == "" {
		util.ErrorPrint(nil, nil, "missing param student_id")
		return context.Error(c, http.StatusBadRequest, "missing param student_id", nil)
	}

	if p.StudentID == "" {
		util.ErrorPrint(nil, nil, "missing param student_id")
		return context.Error(c, http.StatusBadRequest, "missing param student_id", nil)
	}
	if p.Name == "" {
		util.ErrorPrint(nil, nil, "missing param name")
		return context.Error(c, http.StatusBadRequest, "missing param name", nil)
	}
	if p.QQ == 0 {
		util.ErrorPrint(nil, nil, "missing param qq")
		return context.Error(c, http.StatusBadRequest, "missing param qq", nil)
	}
	member, err := m.GetOneFamilyMemberByStudentID(p.StudentID)
	if err == nil && member.StudentID != "" {
		return context.Error(c, http.StatusBadRequest, "student_id already used", nil)
	}

	if err := m.AddOneFamilyMember(p); err != nil {
		util.ErrorPrint(err, nil, "create new family member failed")
		return context.Error(c, http.StatusInternalServerError, "insert in db failed", err)
	}

	return context.Success(c, p)
}

func DeleteOneFamilyMemberRequestHandler(c echo.Context) error {
	m := model.GetModel()
	defer m.Close()

	studentID := c.QueryParam("student_id")

	if err := m.DeleteOneFamilyMemberByStudentID(studentID); err != nil {
		util.ErrorPrint(err, nil, "create new family member failed")
		return context.Error(c, http.StatusInternalServerError, "insert in db failed", err)
	}

	return context.Success(c, nil)
}

func UpdateOneFamilyMemberRequestHandler(c echo.Context) error {
	m := model.GetModel()
	defer m.Close()

	studentID := c.QueryParam("student_id")

	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	p := param.FamilyMember{}
	if err := json.Unmarshal(bodyBytes, &p); err != nil {
		util.ErrorPrint(err, nil, "unmarshal failed on updating family member")
		return context.Error(c, http.StatusBadRequest, "unmarshal failed", err)
	}

	if p.StudentID == "" {
		util.ErrorPrint(nil, nil, "missing param student_id")
		return context.Error(c, http.StatusBadRequest, "missing param student_id", nil)
	}
	if p.Name == "" {
		util.ErrorPrint(nil, nil, "missing param name")
		return context.Error(c, http.StatusBadRequest, "missing param name", nil)
	}
	if p.QQ == 0 {
		util.ErrorPrint(nil, nil, "missing param qq")
		return context.Error(c, http.StatusBadRequest, "missing param qq", nil)
	}

	member, err := m.GetOneFamilyMemberByStudentID(p.StudentID)
	if err == nil && member.StudentID != "" {
		return context.Error(c, http.StatusBadRequest, "student_id already used", nil)
	}

	if _, err := m.UpdateFamilyMemberByStudentID(studentID, p); err != nil {
		util.ErrorPrint(err, nil, "update family member failed")
		return context.Error(c, http.StatusInternalServerError, "insert in db failed", err)
	}

	return context.Success(c, p)
}
