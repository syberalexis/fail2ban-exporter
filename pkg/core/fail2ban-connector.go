package core

import (
	"encoding/json"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/syberalexis/fail2ban-exporter/pkg/database"

	log "github.com/sirupsen/logrus"
)

// ==========================================================
// Regex to extract informations from F2B command line return
var jailListRegexp = regexp.MustCompile(".+Jail list:[\t ]+(.+)")
var currentlyFailedRegexp = regexp.MustCompile(".+Currently failed:[\t ]+(.+)")
var totalFailedRegexp = regexp.MustCompile(".+Total failed:[\t ]+(.+)")
var currentlyBannedRegexp = regexp.MustCompile(".+Currently banned:[\t ]+(.+)")
var totalBannedRegexp = regexp.MustCompile(".+Total banned:[\t ]+(.+)")

// ==========================================================

type Fail2banConnector struct {
	CountryEnabled    bool
	CountryServiceUrl string
}

/*
Get all Jails' state
*/
func (connector *Fail2banConnector) GetJailsState() (*JailsState, error) {
	log.Debug("GetJailState")

	jailState := JailsState{CurrentlyFailed: make(map[string]uint), TotalFailed: make(map[string]uint), CurrentlyBanned: make(map[string]uint), TotalBanned: make(map[string]uint), CountriesBanned: make(map[string]map[string]uint)}
	jails, _ := connector.getJailList()
	for _, jail := range jails {
		cf, tf, cb, tb, _ := connector.getJailInfo(jail)
		jailState.CurrentlyFailed[jail] = cf
		jailState.TotalFailed[jail] = tf
		jailState.CurrentlyBanned[jail] = cb
		jailState.TotalBanned[jail] = tb

		if connector.CountryEnabled {
			bannedIps := connector.getJailBannedIps(jail)
			connector.importCountries(bannedIps)
			jailState.CountriesBanned[jail] = database.GetCountryCountFromIps(bannedIps)
		}
	}

	return &jailState, nil
}

/*
Get the list of enabled jails
*/
func (connector *Fail2banConnector) getJailList() ([]string, error) {
	cmd := exec.Command("/usr/bin/fail2ban-client", "status")
	stdout, err := cmd.Output()
	if err != nil {
		log.Errorf("Impossible to get Jail list from command : %s", err)
		return nil, err
	}
	log.Debugf("Jails : %s", stdout)

	jails := jailListRegexp.FindStringSubmatch(string(stdout))[1]
	log.Debugf("Finded jails : %s", jails)
	return strings.Split(string(jails), ", "), nil
}

/*
Return the info of one jail with it's name
Return Current Failed, Total Failed, Current Banned, Total Banned
*/
func (connector *Fail2banConnector) getJailInfo(jail string) (uint, uint, uint, uint, error) {
	cmd := exec.Command("/usr/bin/fail2ban-client", "status", jail)
	stdout, err := cmd.Output()
	if err != nil {
		return 0, 0, 0, 0, err
	}

	currentlyFailed, _ := strconv.Atoi(string(currentlyFailedRegexp.FindStringSubmatch(string(stdout))[0]))
	totalFailed, _ := strconv.Atoi(string(totalFailedRegexp.FindStringSubmatch(string(stdout))[0]))
	currentlyBanned, _ := strconv.Atoi(string(currentlyBannedRegexp.FindStringSubmatch(string(stdout))[0]))
	totalBanned, _ := strconv.Atoi(string(totalBannedRegexp.FindStringSubmatch(string(stdout))[0]))
	return uint(currentlyFailed), uint(totalFailed), uint(currentlyBanned), uint(totalBanned), nil
}

/*
Return the Banned IPs of one Jail
*/
func (connector *Fail2banConnector) getJailBannedIps(jail string) []string {
	cmd := exec.Command("/usr/bin/fail2ban-client", "get", jail, "banned")
	stdout, err := cmd.Output()

	if err != nil {
		log.Errorf("Error on fail2ban-client command : %s", err)
		return nil
	}

	var ips []string
	if err2 := json.Unmarshal([]byte(strings.ReplaceAll(string(stdout), "'", "\"")), &ips); err2 != nil {
		log.Errorf("Impossible to unmarshal ips : %s", err2)
		return nil
	}

	log.Debugf("Banned ip for %s is %s", jail, stdout)
	return ips
}

/*
Import associated country for each IPs
*/
func (connector *Fail2banConnector) importCountries(ips []string) {
	for _, ip := range ips {
		log.Debugf("Is in DB IP %s", ip)
		if database.GetIpCountryByIp(ip) == nil {
			log.Debugf("Search country for IP %s", ip)
			country, err := GetCountryFromIp(connector.CountryServiceUrl, ip)
			if err == nil {
				log.Debugf("Add country %s for IP %s", *country, ip)
				database.UpdateIpCountryInfo(database.IpCountry{Ip: ip, Country: *country})
			}
		}
	}
}
