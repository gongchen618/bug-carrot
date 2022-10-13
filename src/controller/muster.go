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

type muster struct {
	People []param.FamilyMember
}

func GetAllMusterRequestHandler(c echo.Context) error {
	m := model.GetModel()
	defer m.Close()

	token := c.QueryParam("token")
	if token != util.Token {
		return context.Error(c, http.StatusUnauthorized, "wrong token", nil)
	}

	members, err := m.GetAllFamilyMember()
	if err != nil {
		util.ErrorPrint(err, nil, "get all muster member failed")
		return context.Error(c, http.StatusInternalServerError, "", err)
	}

	return context.Success(c, members)
}

func CreateOneMusterByNameRequestHandler(c echo.Context) error {
	m := model.GetModel()
	defer m.Close()

	token := c.QueryParam("token")
	if token != util.Token {
		return context.Error(c, http.StatusUnauthorized, "wrong token", nil)
	}

	title := c.QueryParam("title")
	if title == "" {
		return context.Error(c, http.StatusBadRequest, "title cannot be empty", nil)
	}

	if err := m.CreateOneMusterByTitle(title); err != nil {
		util.ErrorPrint(err, nil, "create new muster failed")
		return context.Error(c, http.StatusInternalServerError, "insert in db failed", err)
	}

	return context.Success(c, nil)
}

func DeleteOneMusterByNameRequestHandler(c echo.Context) error {
	m := model.GetModel()
	defer m.Close()

	token := c.QueryParam("token")
	if token != util.Token {
		return context.Error(c, http.StatusUnauthorized, "wrong token", nil)
	}

	title := c.QueryParam("title")
	if title == "" {
		return context.Error(c, http.StatusBadRequest, "title cannot be empty", nil)
	}

	if err := m.DeleteOneMusterByTitle(title); err != nil {
		util.ErrorPrint(err, nil, "delete muster failed")
		return context.Error(c, http.StatusInternalServerError, "delete in db failed", err)
	}

	return context.Success(c, nil)
}

type nameString struct {
	Name []string
}

func AddPersonsToOneMusterRequestHandler(c echo.Context) error {
	m := model.GetModel()
	defer m.Close()

	token := c.QueryParam("token")
	if token != util.Token {
		return context.Error(c, http.StatusUnauthorized, "wrong token", nil)
	}

	title := c.QueryParam("title")
	if title == "" {
		return context.Error(c, http.StatusBadRequest, "title cannot be empty", nil)
	}

	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	p := nameString{}
	if err := json.Unmarshal(bodyBytes, &p); err != nil {
		util.ErrorPrint(err, nil, "unmarshal failed")
		return context.Error(c, http.StatusBadRequest, "unmarshal failed", err)
	}

	ms, err := m.AddPersonsToOneMuster(title, p.Name)
	if err != nil {
		util.ErrorPrint(err, nil, "create muster failed")
		return context.Error(c, http.StatusInternalServerError, "create in db failed", err)
	}

	return context.Success(c, ms)
}

func DeletePersonsOnOneMusterRequestHandler(c echo.Context) error {
	m := model.GetModel()
	defer m.Close()

	token := c.QueryParam("token")
	if token != util.Token {
		return context.Error(c, http.StatusUnauthorized, "wrong token", nil)
	}

	title := c.QueryParam("title")
	if title == "" {
		return context.Error(c, http.StatusBadRequest, "title cannot be empty", nil)
	}

	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	p := nameString{}
	if err := json.Unmarshal(bodyBytes, &p); err != nil {
		util.ErrorPrint(err, nil, "unmarshal failed")
		return context.Error(c, http.StatusBadRequest, "unmarshal failed", err)
	}

	ms, err := m.DeletePersonsOnOneMuster(title, p.Name)
	if err != nil {
		util.ErrorPrint(err, nil, "delete muster failed")
		return context.Error(c, http.StatusInternalServerError, "delete in db failed", err)
	}

	return context.Success(c, ms)
}
