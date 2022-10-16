package controller

import (
	"bug-carrot/model"
	"bug-carrot/util"
	"bug-carrot/util/context"
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
)

func GetAllBallotRequestHandler(c echo.Context) error {
	m := model.GetModel()
	defer m.Close()

	token := c.QueryParam("token")
	if token != util.Token {
		return context.Error(c, http.StatusUnauthorized, "wrong token", nil)
	}

	ballots, err := m.GetAllBallot()
	if err != nil {
		util.ErrorPrint(err, nil, "get all ballots failed")
		return context.Error(c, http.StatusInternalServerError, "", err)
	}

	return context.Success(c, ballots)
}

func CreateOneBallotByTitleRequestHandler(c echo.Context) error {
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
	musterTitle := c.QueryParam("muster")
	ms, err := m.GetOneMusterByTitle(musterTitle)
	if musterTitle == "" || err != nil {
		return context.Error(c, http.StatusBadRequest, "muster cannot be empty", nil)
	}
	defaultOption := c.QueryParam("default_option")
	if defaultOption == "" {
		defaultOption = "未填写"
	}

	if err := m.CreateOneBallotByTitle(title, ms, defaultOption); err != nil {
		util.ErrorPrint(err, nil, "create new muster failed")
		return context.Error(c, http.StatusInternalServerError, "insert in db failed", err)
	}

	return context.Success(c, nil)
}

func DeleteOneBallotByTitleRequestHandler(c echo.Context) error {
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

	if err := m.DeleteOneBallotByTitle(title); err != nil {
		util.ErrorPrint(err, nil, "delete muster failed")
		return context.Error(c, http.StatusInternalServerError, "delete in db failed", err)
	}

	return context.Success(c, nil)
}

func AddAnOptionToOneBallotRequestHandler(c echo.Context) error {
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
	option := c.QueryParam("option")
	if option == "" {
		return context.Error(c, http.StatusBadRequest, "option cannot be empty", nil)
	}

	bt, err := m.AddAnOptionToOneBallot(title, option)
	if err != nil {
		util.ErrorPrint(err, nil, "add option failed")
		return context.Error(c, http.StatusInternalServerError, "insert in db failed", err)
	}

	return context.Success(c, bt)
}

func DeleteAnOptionOnOneBallotRequestHandler(c echo.Context) error {
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
	option := c.QueryParam("option")
	if option == "" {
		return context.Error(c, http.StatusBadRequest, "option cannot be empty", nil)
	}

	bt, err := m.DeleteAnOptionOnOneBallot(title, option)
	if err != nil {
		util.ErrorPrint(err, nil, "delete option failed")
		return context.Error(c, http.StatusInternalServerError, "delete in db failed", err)
	}

	return context.Success(c, bt)
}

func UpdateOptionsOnOneBallotForMembersRequestHandler(c echo.Context) error {
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
	option := c.QueryParam("option")
	if option == "" {
		return context.Error(c, http.StatusBadRequest, "option cannot be empty", nil)
	}

	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	p := nameString{}
	if err := json.Unmarshal(bodyBytes, &p); err != nil {
		util.ErrorPrint(err, nil, "unmarshal failed")
		return context.Error(c, http.StatusBadRequest, "unmarshal failed", err)
	}

	bt, err := m.UpdateOptionsOnOneBallotForMembers(title, option, p.Name)
	if err != nil {
		util.ErrorPrint(err, nil, "update options failed")
		return context.Error(c, http.StatusInternalServerError, "update in db failed", err)
	}

	return context.Success(c, bt)
}
