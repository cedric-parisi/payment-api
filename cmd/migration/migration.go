package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cedric-parisi/payment-api/internal/config"
	"github.com/jinzhu/gorm"

	"github.com/cedric-parisi/payment-api/internal/models"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	connectionString = "host=%s port=%s user=%s dbname=%s password=%s sslmode=disable"
)

func main() {
	// Marshal the mock.json file into payments
	mockFile, err := os.Open("mock.json")
	if err != nil {
		log.Fatal("could not open mock.json: ", err)
	}

	var dest struct {
		Data []*models.Payment `json:"data"`
	}
	if err := json.NewDecoder(mockFile).Decode(&dest); err != nil {
		log.Fatal("could not decode mockFile: ", err)
	}

	for _, p := range dest.Data {
		p.CreatedAt = time.Now().UTC()
	}

	// setup config
	cfg := config.SetConfiguration()

	// open connection to db
	db, err := gorm.Open("postgres", fmt.Sprintf(connectionString, cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbName, cfg.DbPassword))
	if err != nil {
		log.Fatal("could not open db connection: ", err)
	}
	defer db.Close()

	db.LogMode(true)

	// create/update schemas according to struct defintions
	db.AutoMigrate(&models.Payment{}, &models.Attribute{}, &models.BeneficiaryParty{}, &models.ChargesInformation{}, &models.DebtorParty{}, &models.Fx{}, &models.SenderCharge{}, &models.SponsorParty{})

	// insert mock data
	for _, p := range dest.Data {
		db.Create(p)
	}
}
