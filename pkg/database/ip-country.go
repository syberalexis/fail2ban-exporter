package database

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitIpDB(path string) error {
	var err error

	sqlitedb := sqlite.Open(path)
	db, err = gorm.Open(sqlitedb, &gorm.Config{})

	if err != nil {
		return err
	}

	db.Migrator().CreateTable(&IpCountry{})
	db.Migrator().AutoMigrate(&IpCountry{})

	return nil
}

type IpCountry struct {
	Ip      string `gorm:"primaryKey"`
	Country string
}

func UpdateIpCountryInfo(ipCountry IpCountry) error {
	result := db.Save(&ipCountry)
	return result.Error
}

func GetIpCountryByIp(ip string) *IpCountry {
	log.Debugf("GetIpCountryByIp : %s", ip)

	var ipCountry *IpCountry
	result := db.Find(&ipCountry, ip)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return ipCountry
}

func GetCountryCountFromIps(ips []string) map[string]uint {
	log.Debugf("GetCountryCountFromIps : %s", ips)

	var countryCount map[string]uint
	result := db.Raw("SELECT count(ip) FROM ipcountry WHERE ip IN ? GROUP BY country", ips).Scan(&countryCount)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return countryCount
}
