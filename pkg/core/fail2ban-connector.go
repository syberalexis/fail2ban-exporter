package core

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/syberalexis/fail2ban-exporter/pkg/database"

	log "github.com/sirupsen/logrus"
)

var jailListRegexp = regexp.MustCompile(".+Jail list:[\t ]+(.+)")
var currentlyFailedRegexp = regexp.MustCompile(".+Currently failed:[\t ]+(.+)")
var totalFailedRegexp = regexp.MustCompile(".+Total failed:[\t ]+(.+)")
var currentlyBannedRegexp = regexp.MustCompile(".+Currently banned:[\t ]+(.+)")
var totalBannedRegexp = regexp.MustCompile(".+Total banned:[\t ]+(.+)")

type Fail2banConnector struct {
	CountryEnabled    bool
	CountryServiceUrl string
}

func (connector *Fail2banConnector) GetJailState() (*JailsState, error) {
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

func (connector *Fail2banConnector) getJailList() ([]string, error) {
	cmd := exec.Command("/usr/bin/fail2ban-client", "status")
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	log.Debugf("Jails : %s", stdout)

	jails := jailListRegexp.FindStringSubmatch(string(stdout))[1]
	log.Debugf("Finded jails : %s", jails)
	return strings.Split(string(jails), ", "), nil
}

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

func (connector *Fail2banConnector) getJailBannedIps(jail string) []string {
	return database.GetIpsByJail(jail)
}

func (connector *Fail2banConnector) importCountries(ips []string) {
	for _, ip := range ips {
		country, err := GetCountryFromIp(connector.CountryServiceUrl, ip)
		if err != nil {
			log.Debugf("Add country %s for IP %s", *country, ip)
			database.UpdateIpCountryInfo(database.IpCountry{Ip: ip, Country: *country})
		}
	}
}
