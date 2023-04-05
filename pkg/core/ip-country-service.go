package core

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetCountryFromIp(url string, ip string) (*string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", url, ip), nil)
	if err != nil {
		return nil, err
	}
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
	if err := json.NewDecoder(res.Body).Decode(&ipCountry); err != nil {
		return nil, err
	}
	return &ipCountry.Country, nil
}
