package customer

import (
	"fmt"
	"net/http"
	"www-api/config"
	"www-api/internal/constants"
	"www-api/internal/datatypes"
	"www-api/internal/logger"
	service "www-api/service/customer"
	"www-api/utils"

	"github.com/gin-gonic/gin"
)

type CustomerAPI struct {
	config               config.Config
	log                  logger.ZapLogger
	getPrivacyStatus     func(fid string) (map[string]int, error)
	getTimezone          func(fid string) (datatypes.TimezoneResponse, error)
	getNotificationEmail func(fid string) (datatypes.Notification, error)
	getFilterType        func(fid string) (string, error)
}

func NewCustomerAPI(conf config.Config, log logger.ZapLogger, connections *datatypes.Connections) CustomerAPI {
	serv := service.NewCustomerService(log, connections)
	return CustomerAPI{
		config:               conf,
		log:                  log,
		getPrivacyStatus:     serv.ProuctPrivacyStatus,
		getTimezone:          serv.Timezone,
		getNotificationEmail: serv.Notification,
		getFilterType:        serv.GetFilterType,
	}
}

// @Summary      Get Privacy Status
// @Description  fetches info of a student
// @Tags         Customer
// @Produce      json
// @Success      200 {object} map[string]int
// @Failure      400 {object} string
// @Failure      404 {object} string
// @Failure      500 {object} string
// @Router       /api/customer/privacy/status [get]
func (r CustomerAPI) PrivacyStatus(c *gin.Context) {
	var request datatypes.CustomerRequest
	err := c.BindJSON(&request)
	if err != nil {
		r.log.Error("error binding request body", map[string]interface{}{"error": err})
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = utils.ValidateFid(request.Fid, r.log)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	privacyStatus, err := r.getPrivacyStatus(request.Fid)
	if err != nil {
		if err == constants.ResourceNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("user %s doesn't exists", request.Fid)})
			return
		}
		r.log.Error("error occured while fecthing student info", map[string]interface{}{"error": err})
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return

	}

	r.log.Info("privacy_status", map[string]interface{}{"privacyStatus": privacyStatus})
	c.JSON(http.StatusOK, privacyStatus)
}

// @Summary      Get Timezone
// @Description  fetches timezone for a fid
// @Tags         Customer
// @Produce      json
// @Success      200 {object} datatypes.TimezoneResponse
// @Failure      400 {object} string
// @Failure      404 {object} string
// @Failure      500 {object} string
// @Router       /api/customer/timezone [get]
func (r CustomerAPI) Timezone(c *gin.Context) {
	var request datatypes.CustomerRequest
	err := c.BindJSON(&request)
	if err != nil {
		r.log.Error("error binding request body", map[string]interface{}{"error": err})
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = utils.ValidateFid(request.Fid, r.log)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	timezone, err := r.getTimezone(request.Fid)
	if err != nil {
		if err == constants.ResourceNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("user %s doesn't exists", request.Fid)})
			return
		}
		r.log.Error("error occured while fecthing student info", map[string]interface{}{"error": err})
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, timezone)
}

// @Summary      Get Notification
// @Description  fetches notification email
// @Tags         Customer
// @Produce      json
// @Success      200 {object} datatypes.Notification
// @Failure      400 {object} string
// @Failure      404 {object} string
// @Failure      500 {object} string
// @Router       /api/customer/notification/config/aware [get]
func (r CustomerAPI) Notification(c *gin.Context) {
	var request datatypes.CustomerRequest
	err := c.BindJSON(&request)
	if err != nil {
		r.log.Error("error binding request body", map[string]interface{}{"error": err})
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = utils.ValidateFid(request.Fid, r.log)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	notification, err := r.getNotificationEmail(request.Fid)
	if err != nil {
		if err == constants.ResourceNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("user %s doesn't exists", request.Fid)})
			return
		}
		r.log.Error("error occured while fecthing student info", map[string]interface{}{"error": err})
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return

	}

	c.JSON(http.StatusOK, notification)
}

// @Summary      Get Filter
// @Description  fetches filters for an fid
// @Tags         Customer
// @Produce      json
// @Success      200 {object} datatypes.FilterType
// @Failure      400 {object} string
// @Failure      404 {object} string
// @Failure      500 {object} string
// @Router       /api/customer/filter-type [get]
func (r CustomerAPI) FilterType(c *gin.Context) {
	var request datatypes.CustomerRequest
	err := c.BindJSON(&request)
	if err != nil {
		r.log.Error("error binding request body", map[string]interface{}{"error": err})
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = utils.ValidateFid(request.Fid, r.log)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	filterType, err := r.getFilterType(request.Fid)
	if err != nil {
		r.log.Error("error occured while fecthing student info", map[string]interface{}{"error": err})
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"filteringType": filterType})
}
