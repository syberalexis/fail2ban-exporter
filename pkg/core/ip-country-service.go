package core

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func GetCountryFromIp(url string, ip string) (*string, error) {
	// Init request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", url, ip), nil)
	if err != nil {
		return nil, err
	}
	// Do request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// body, err := ioutil.ReadAll(res.Body)
	var ipCountry struct {
		Ip      string
		Country string
	}

	// Decode data
	if err := json.NewDecoder(res.Body).Decode(&ipCountry); err != nil {
		return nil, err
	}

	log.Debugf("Country for ip %s is : %s", ipCountry.Ip, ipCountry.Country)
	return &ipCountry.Country, nil
}
