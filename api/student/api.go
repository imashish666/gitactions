package student

import (
	"fmt"
	"net/http"
	"www-api/config"
	"www-api/internal/constants"
	"www-api/internal/datatypes"
	"www-api/internal/logger"
	service "www-api/service/student"
	"www-api/utils"

	"github.com/gin-gonic/gin"
)

type InfoAPI struct {
	config         config.Config
	log            logger.ZapLogger
	getStudentInfo func(fid, email string) (datatypes.StudentInfo, error)
}

func NewInfoAPI(conf config.Config, log logger.ZapLogger, connections *datatypes.Connections) InfoAPI {
	serv := service.NewStudentService(log, connections)
	return InfoAPI{
		config:         conf,
		log:            log,
		getStudentInfo: serv.StudentInfo,
	}
}

// @Summary      Get student info
// @Description  fetches info of a student
// @Tags         Student
// @Produce      json
// @Success      200 {object} datatypes.StudentInfo
// @Failure      400 {object} string
// @Failure      404 {object} string
// @Failure      500 {object} string
// @Router       /api/user [get]
func (r InfoAPI) GetInfo(c *gin.Context) {
	var request datatypes.StudentInfoRequest
	err := c.BindJSON(&request)
	if err != nil {
		r.log.Error("error binding request body", map[string]interface{}{"error": err})
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = utils.ValidateEmail(request.Email, r.log)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	studentInfo, err := r.getStudentInfo("", request.Email)
	if err != nil {
		if err == constants.ResourceNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("user %s doesn't exists", request.Email)})
			return
		}
		r.log.Error("error occured while fecthing student info", map[string]interface{}{"error": err})
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return

	}

	r.log.Info("student info", map[string]interface{}{"student_info": studentInfo})
	c.JSON(http.StatusOK, studentInfo)
}
