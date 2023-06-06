package main

import (
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Promotion struct {
	ID             string  `json:"id"`
	Price          float64 `json:"price"`
	ExpirationDate string  `json:"expiration_date"`
}

const PromotionsFileName = "promotions.csv"

var promotions map[string]Promotion

func getPromotionById(c *gin.Context) {
	id := c.Param("id")

	promotion, ok := promotions[id]

	if !ok {
		c.Status(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, promotion)
	}
}

func main() {
	startFileWatcher()

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

func startFileWatcher() {
	loadData(PromotionsFileName)
	ticker := time.NewTicker(30 * time.Minute)

	go func() {
		for range ticker.C {
			loadData(PromotionsFileName)
		}
	}()
}

func loadData(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	data := make(map[string]Promotion)

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		price, err := strconv.ParseFloat(record[1], 64)

		r := Promotion{
			ID:             record[0],
			Price:          price,
			ExpirationDate: record[2],
		}

		data[r.ID] = r
	}

	promotions = data
}
