package interp_a

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ServerParams struct {
	Address string
}

func HostServer(
	connParams ServerParams,
	evaluator CanEvaluate,
) error {
	router := gin.Default()
	router.POST("/call", func(c *gin.Context) {
		listStr := c.PostForm("list")
		list := []interface{}{}

		err := json.Unmarshal([]byte(listStr), &list)
		if err != nil {
			logrus.Error(err)
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		result, err := evaluator.OpEvaluate(list)
		if err != nil {
			logrus.Error(err)
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, result)

	})

	err := router.Run(connParams.Address)
	if err != nil {
		return err
	}

	return nil
}
