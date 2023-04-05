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
	Ip      string `gorm:"type:inet;primaryKey"`
	Country string
}

func UpdateIpCountryInfo(ipCountry IpCountry) error {
	result := db.Save(&ipCountry)
	return result.Error
}

func GetIpCountryByIp(ip string) *IpCountry {
	log.Debugf("GetIpCountryByIp : %s", ip)

	var ipCountry *IpCountry
	// result := db.Raw("SELECT * FROM ip_countries WHERE ip = '?'", ip).Scan(&ipCountry)
	result := db.Table("ip_countries").Select("country").Where("ip = ?", ip).First(&ipCountry)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return ipCountry
}

func GetCountryCountFromIps(ips []string) map[string]uint {
	log.Debugf("GetCountryCountFromIps : %s", ips)

	type CountryCount struct {
		Country string
		Count   uint64
	}
	var countryCounts []CountryCount
	// result := db.Raw("SELECT country, count(ip) FROM ip_countries WHERE ip IN (?) GROUP BY country", ips).Scan(&countryCounts)
	result := db.Table("ip_countries").Select("country, count(ip)").Where("ip IN ?", ips).Group("country").Scan(&countryCounts)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	countriesCount := make(map[string]uint)
	for _, countryCount := range countryCounts {
		log.Debugf("CountryCount : %s - %d", countryCount.Country, countryCount.Count)
		countriesCount[countryCount.Country] = uint(countryCount.Count)
	}

	return countriesCount
}
