package atRisk

import (
	"net/http"
	"strconv"
	"www-api/config"
	"www-api/internal/constants"
	"www-api/internal/datatypes"
	"www-api/internal/logger"
	service "www-api/service/at-risk"
	"www-api/utils"

	"github.com/gin-gonic/gin"
)

type RiskAPI struct {
	config        config.Config
	log           logger.ZapLogger
	createCache   func(key, value string) (datatypes.AtRiskResponse, error)
	deleteCache   func(key string) (datatypes.AtRiskResponse, error)
	getScore      func(email string) ([]datatypes.RiskScore, error)
	extentTTL     func(email string, ttl int) error
	getEventScore func(email, timestamp, mid string) (datatypes.EventScoreResponse, error)
}

func NewRiskAPI(conf config.Config, log logger.ZapLogger, connections *datatypes.Connections) RiskAPI {
	connections.Redis[constants.AtRiskReadRedisKey].Options().DB = constants.RedisDB6
	connections.Redis[constants.AtRiskWriteRedisKey].Options().DB = constants.RedisDB6
	serv := service.NewRiskService(log, connections)
	return RiskAPI{
		config:        conf,
		log:           log,
		createCache:   serv.CreateCache,
		deleteCache:   serv.DeleteCache,
		getScore:      serv.GetScore,
		extentTTL:     serv.ExtendTTL,
		getEventScore: serv.GetEventScore,
	}
}

// @Summary      Create a cache
// @Description  add/update a value in cache
// @Tags         AtRisk
// @Produce      json
// @Success      200 {object} datatypes.AtRiskResponse
// @Failure      400 {object} string
// @Failure      404 {object} string
// @Failure      500 {object} string
// @Router       /at-risk/cache/create [post]
func (r RiskAPI) CreateCache(c *gin.Context) {
	var request datatypes.CacheRequest
	err := c.BindJSON(&request)
	if err != nil {
		r.log.Error("error binding request body", map[string]interface{}{"error": err})
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = utils.ValidateAtRiskKey(request.AtRiskKey, r.log)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = utils.ValidateAtRiskValue(request.AtRiskValue, r.log)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	score, err := r.createCache(request.AtRiskKey, request.AtRiskValue)
	if err != nil {
		r.log.Error("error occured while setting key to cache", map[string]interface{}{"error": err})
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	r.log.Info("cache created", map[string]interface{}{"key": request.AtRiskKey, "value": request.AtRiskValue, "totalAtRiskScore": score})
	c.JSON(http.StatusOK, gin.H{"totalAtRiskScore": score})
}

// @Summary      Delete a cache key
// @Description  removes a key from cache
// @Tags         AtRisk
// @Produce      json
// @Success      200 {object} datatypes.AtRiskResponse
// @Failure      400 {object} string
// @Failure      404 {object} string
// @Failure      500 {object} string
// @Router       /at-risk/cache/delete [delete]
func (r RiskAPI) DeleteCache(c *gin.Context) {
	var request datatypes.CacheRequest
	err := c.BindJSON(&request)
	if err != nil {
		r.log.Error("error binding request body", map[string]interface{}{"error": err})
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = utils.ValidateAtRiskKey(request.AtRiskKey, r.log)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	score, err := r.deleteCache(request.AtRiskKey)
	if err != nil {
		if err == constants.ResourceNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"message": "key doesn't exists"})
			return
		}
		r.log.Error("error occured while deleting cache", map[string]interface{}{"error": err})
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	r.log.Info("cache deleted", map[string]interface{}{"key": request.AtRiskKey, "totalAtRiskScore": score})
	c.JSON(http.StatusOK, score)
}

// @Summary      Get a score
// @Description  fetches score from database
// @Tags         AtRisk
// @Produce      json
// @Success      200 {array} datatypes.RiskScore
// @Failure      400 {object} string
// @Failure      404 {object} string
// @Failure      500 {object} string
// @Router       /at-risk/score [get]
func (r RiskAPI) Score(c *gin.Context) {
	var request datatypes.AtRiskRequest
	err := c.BindJSON(&request)
	if err != nil {
		r.log.Error("error binding request body", map[string]interface{}{"error": err})
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = utils.ValidateEmail(request.UserEmail, r.log)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	scores, err := r.getScore(request.UserEmail)
	if err != nil {
		r.log.Error("error occured while fetching scores from database", map[string]interface{}{"error": err})
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	r.log.Info("successfully fetched atRiskScore", map[string]interface{}{"email": request.UserEmail, "atRiskScores": scores})
	c.JSON(http.StatusOK, scores)
}

// @Summary      Extent TTL
// @Description  extends the expiry for a key in cache
// @Tags         AtRisk
// @Produce      json
// @Param        userEmail query string true "user email"
// @Param        timestamp query string true "timestamp"
// @Success      200 {object} string
// @Failure      400 {object} string
// @Failure      404 {object} string
// @Failure      500 {object} string
// @Router       /at-risk/extend-ttl [post]
func (r RiskAPI) ExtendTTL(c *gin.Context) {
	var request datatypes.AtRiskRequest
	err := c.BindJSON(&request)
	if err != nil {
		r.log.Error("error binding request body", map[string]interface{}{"error": err})
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = utils.ValidateEmail(request.UserEmail, r.log)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var ttl int

	if request.TTL == "" {
		ttl = 60 * 60 * 24 * 90
		r.log.Info("setting deafult ttl", map[string]interface{}{"ttl": ttl})
	} else {
		ttl, err = strconv.Atoi(request.TTL)
		if err != nil {
			r.log.Error("invalid ttl value received", map[string]interface{}{"error": err, "ttl": request.TTL})
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid ttl in request body, should be numeric"})
			return
		}
	}

	err = r.extentTTL(request.UserEmail, ttl)
	if err != nil {
		r.log.Error("error occured while extending ttl", map[string]interface{}{"email": request.UserEmail, "ttl": ttl})
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	r.log.Info("successfully set ttl", map[string]interface{}{"email": request.UserEmail, "ttl": ttl})
	c.JSON(http.StatusOK, "ttl extended")
}

// @Summary      Get event score details
// @Description  fetches score for a specific event
// @Tags         AtRisk
// @Produce      json
// @Success      200 {object} datatypes.EventScoreResponse
// @Failure      400 {object} string
// @Failure      404 {object} string
// @Failure      500 {object} string
// @Router       /at-risk/event-score-details [get]
func (r RiskAPI) EventScore(c *gin.Context) {
	var request datatypes.AtRiskRequest
	err := c.BindJSON(&request)
	if err != nil {
		r.log.Error("error binding request body", map[string]interface{}{"error": err})
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = utils.ValidateEmail(request.UserEmail, r.log)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = utils.ValidateTimestamp(request.Timestamp, r.log)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	score, err := r.getEventScore(request.UserEmail, request.Timestamp, request.Mid)
	if err != nil {
		if err == constants.ResourceNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"message": "key doesn't exists"})
			return
		}
		r.log.Error("error occured while fetching event score, error", map[string]interface{}{"error": err})
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	r.log.Info("successfully fetched EventScoreDetails", map[string]interface{}{"key": score.AtRiskKey, "value": score.AtRiskValue, "score": score.AtRiskScore})
	c.JSON(http.StatusOK, score)
}
