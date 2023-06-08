package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gurgendav/verve-storage/models"
	"github.com/gurgendav/verve-storage/pkg/gredis"
	"log"
	"net/http"
)

func getPromotionById(c *gin.Context) {
	id := c.Param("id")

	promotionString, err := gredis.Get(id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	var promotion models.Promotion
	err = json.Unmarshal([]byte(promotionString), &promotion)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, promotion)
	}
}

func main() {
	r := gin.Default()

	r.GET("/promotions/:id", getPromotionById)

	err := r.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = r.Run()
	if err != nil {
		log.Fatal(err)
		return
	}
}
