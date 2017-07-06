package api

import (
	"agaza/config"
	"agaza/logger"
	"agaza/models"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"agaza/connectionManager"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

const (
	apiBase       = "/api"
	rootURI       = apiBase + "/"
	leavesURI     = apiBase + "/leaves"
	leavesIDURI   = apiBase + "/leaves/<leave_id>"
	agazaTypesURI = apiBase + "/agazatypes"
)

var (
	configuration = config.LoadConfiguration()
)

//FasthttpAPIHandler is structure wrapper for the fast http objects
type FasthttpAPIHandler struct {
	router *routing.Router
}

//StartAndServeAPIs will start the server and attach the endpoints to the router
func (f *FasthttpAPIHandler) StartAndServeAPIs() {
	f.router = routing.New()
	f.router.Get(rootURI, f.SetHeaders, f.Logger, f.GetRoot)
	f.router.Get(leavesURI, f.SetHeaders, f.Logger, f.GetLeaves)
	f.router.Post(leavesURI, f.SetHeaders, f.Logger, f.PostLeaves)
	f.router.Put(leavesIDURI, f.SetHeaders, f.Logger, f.PutLeaves)

	f.router.Get(agazaTypesURI, f.SetHeaders, f.Logger, f.GetAgazaTypes)

	fasthttp.ListenAndServe(":"+configuration.APIserverport, f.router.HandleRequest)
}

//GetRoot the root page
func (f *FasthttpAPIHandler) GetRoot(c *routing.Context) error {
	agaza, _ := json.Marshal("agaza")
	c.Write(agaza)
	return nil
}

//GetAgazaTypes get all agaza types
func (f *FasthttpAPIHandler) GetAgazaTypes(c *routing.Context) error {
	agaza := new(models.AgazaType)
	agaza.ID = "1"
	agaza.Name = "annual"
	response, _ := json.Marshal(agaza)
	c.Write(response)
	c.SetStatusCode(http.StatusOK)
	return nil
}

//GetLeaves the root page
func (f *FasthttpAPIHandler) GetLeaves(c *routing.Context) error {
	const shortForm = "2006-Jan-02"

	from, err := time.Parse(shortForm, string(c.QueryArgs().Peek("from")))
	if err != nil && string(c.QueryArgs().Peek("from")) != "" {
		logger.Error.Println("date must be in format 2017-Jun-05")
		c.SetStatusCode(http.StatusBadRequest)
		return errors.New("date must be in format 2017-Jun-05")
	}
	to, err := time.Parse(shortForm, string(c.QueryArgs().Peek("to")))
	if err != nil && string(c.QueryArgs().Peek("to")) != "" {
		logger.Error.Println("date must be in format 2017-Jun-05")
		c.SetStatusCode(http.StatusBadRequest)
		return errors.New("date must be in format 2017-Jun-05")
	}

	//userName := string(c.QueryArgs().Peek("user"))
	//departmentID := string(c.QueryArgs().Peek("department"))

	if to.String() != "0001-01-01 00:00:00 +0000 UTC" && from.String() != "0001-01-01 00:00:00 +0000 UTC" {
		connection := connectionManager.GetRedisConnection()
		leaves, err := connection.GetLeavesInRange(from, to)
		if err != nil {
			logger.Error.Println("Failed to get leaves")
			c.SetStatusCode(http.StatusInternalServerError)
			return err
		}
		response, err := json.Marshal(leaves)
		if err != nil {
			logger.Error.Println("Failed to marshal leaves")
			c.SetStatusCode(http.StatusInternalServerError)
			return err
		}
		c.Write(response)
		c.SetStatusCode(http.StatusOK)
		return nil
	}

	return nil
}

//PostLeaves posts a new leve request
func (f *FasthttpAPIHandler) PostLeaves(c *routing.Context) error {
	const shortForm = "2006-Jan-02"
	userID := string(c.PostArgs().Peek("user_id"))
	reason := string(c.PostArgs().Peek("reason"))
	leaveType := string(c.PostArgs().Peek("type"))
	departmentID := string(c.PostArgs().Peek("departmentID"))

	from, err := time.Parse(shortForm, string(c.QueryArgs().Peek("from")))
	if err != nil && string(c.QueryArgs().Peek("from")) != "" {
		logger.Error.Println("date must be in format 2017-Jun-05")
		c.SetStatusCode(http.StatusBadRequest)
		return errors.New("date must be in format 2017-Jun-05")
	}
	to, err := time.Parse(shortForm, string(c.QueryArgs().Peek("to")))
	if err != nil && string(c.QueryArgs().Peek("to")) != "" {
		logger.Error.Println("date must be in format 2017-Jun-05")
		c.SetStatusCode(http.StatusBadRequest)
		return errors.New("date must be in format 2017-Jun-05")
	}

	connection := connectionManager.GetRedisConnection()
	leave := new(models.Leave)
	leave.DepartmentID = departmentID
	leave.From = from
	leave.To = to
	leave.UserID = userID
	leave.Type, _ = strconv.Atoi(leaveType)
	leave.Status = "approved"
	user, err := connection.GetUserByID(userID)
	leave.UserName = user.Name
	leave.Reason = reason

	if err != nil {
		logger.Error.Println("Failed to get user")
		c.SetStatusCode(http.StatusInternalServerError)
		return err
	}

	resID, err := connection.NewLeave(leave)
	if err != nil {
		logger.Error.Println("Failed to create leave")
		c.SetStatusCode(http.StatusInternalServerError)
		return err
	}

	leave.ID = resID

	exposedLeave := new(models.ExposedLeave)
	exposedLeave.DepartmentID = departmentID
	exposedLeave.From = from
	exposedLeave.To = to
	exposedLeave.UserName = user.Name
	exposedLeave.ID = resID
	exposedLeave.UserID = userID

	response, err := json.Marshal(exposedLeave)
	if err != nil {
		logger.Error.Println("Failed to marshal leaves")
		c.SetStatusCode(http.StatusInternalServerError)
		return err
	}
	c.Write(response)
	c.SetStatusCode(http.StatusCreated)

	return nil
}

//PutLeaves edit an existing leave
func (f *FasthttpAPIHandler) PutLeaves(c *routing.Context) error {
	const shortForm = "2006-Jan-02"
	leaveID := string(c.Param("leave_id"))

	userID := string(c.PostArgs().Peek("user_id"))
	reason := string(c.PostArgs().Peek("reason"))
	leaveType := string(c.PostArgs().Peek("type"))
	departmentID := string(c.PostArgs().Peek("departmentID"))

	from, err := time.Parse(shortForm, string(c.QueryArgs().Peek("from")))
	if err != nil && string(c.QueryArgs().Peek("from")) != "" {
		logger.Error.Println("date must be in format 2017-Jun-05")
		c.SetStatusCode(http.StatusBadRequest)
		return errors.New("date must be in format 2017-Jun-05")
	}
	to, err := time.Parse(shortForm, string(c.QueryArgs().Peek("to")))
	if err != nil && string(c.QueryArgs().Peek("to")) != "" {
		logger.Error.Println("date must be in format 2017-Jun-05")
		c.SetStatusCode(http.StatusBadRequest)
		return errors.New("date must be in format 2017-Jun-05")
	}

	connection := connectionManager.GetRedisConnection()
	leave := new(models.Leave)
	leave.DepartmentID = departmentID
	leave.From = from
	leave.To = to
	leave.UserID = userID
	leave.Type, _ = strconv.Atoi(leaveType)
	leave.Status = "approved"
	user, err := connection.GetUserByID(userID)
	leave.UserName = user.Name
	leave.Reason = reason
	leave.ID = leaveID

	if err != nil {
		logger.Error.Println("Failed to get user")
		c.SetStatusCode(http.StatusInternalServerError)
		return err
	}

	resID, err := connection.EditLeave(leave)
	if err != nil {
		logger.Error.Println("Failed to create leave")
		c.SetStatusCode(http.StatusInternalServerError)
		return err
	}

	leave.ID = resID

	exposedLeave := new(models.ExposedLeave)
	exposedLeave.DepartmentID = departmentID
	exposedLeave.From = from
	exposedLeave.To = to
	exposedLeave.UserName = user.Name
	exposedLeave.ID = resID
	exposedLeave.UserID = userID

	response, err := json.Marshal(exposedLeave)
	if err != nil {
		logger.Error.Println("Failed to marshal leaves")
		c.SetStatusCode(http.StatusInternalServerError)
		return err
	}
	c.Write(response)
	c.SetStatusCode(http.StatusOK)

	return nil
}
