package database

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var fail2banDB *gorm.DB

func InitFail2banDB(dbPath string) error {
	var err error

	sqlitedb := sqlite.Open(dbPath)
	fail2banDB, err = gorm.Open(sqlitedb, &gorm.Config{})

	if err != nil {
		return err
	}

	return nil
}

func GetIpsByJail(jail string) []string {
	var ips []string
	result := fail2banDB.Table("bans").Where("jail = ?", jail).Take(&ips)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return ips
}

func GetIpsByJailLimit(jail string, limit int) []string {
	var ips []string
	result := fail2banDB.Table("bans").Where("jail = ?", jail).Limit(limit).Last(&ips)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return ips
}
