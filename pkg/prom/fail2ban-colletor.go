package prom

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/fail2ban-exporter/pkg/core"
)

const USED = "used"
const PRODUCED = "produced"

// Fail2banCollector object to describe and collect metrics
type Fail2banCollector struct {
	connector              core.Fail2banConnector
	jailFailedCurrent      *prometheus.Desc
	jailFailedTotal        *prometheus.Desc
	jailBannedCurrent      *prometheus.Desc
	jailBannedTotal        *prometheus.Desc
	jailBannedCountryTotal *prometheus.Desc
}

// NewFail2banCollector method to construct Fail2banCollector
func NewFail2banCollector(connector core.Fail2banConnector) *Fail2banCollector {
	return &Fail2banCollector{
		connector: connector,
		jailFailedCurrent: prometheus.NewDesc("fail2ban_jail_failed_current",
			"Number of currently failed.", []string{"jail"}, nil,
		),
		jailFailedTotal: prometheus.NewDesc("fail2ban_jail_failed_total",
			"Number of total failed.", []string{"jail"}, nil,
		),
		jailBannedCurrent: prometheus.NewDesc("fail2ban_jail_banned_current",
			"Number of currently banned.", []string{"jail"}, nil,
		),
		jailBannedTotal: prometheus.NewDesc("fail2ban_jail_banned_total",
			"Number of total banned.", []string{"jail"}, nil,
		),
		jailBannedCountryTotal: prometheus.NewDesc("fail2ban_jail_country_banned_total",
			"Number of total banned by country.", []string{"jail", "country"}, nil,
		),
	}
}

// Describe implements required describe function for all prometheus collectors
func (collector *Fail2banCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.jailFailedCurrent
	ch <- collector.jailFailedTotal
	ch <- collector.jailBannedCurrent
	ch <- collector.jailBannedTotal
	ch <- collector.jailBannedCountryTotal
}

// Collect implements required collect function for all prometheus collectors
func (collector *Fail2banCollector) Collect(ch chan<- prometheus.Metric) {
	jailsState, err := collector.connector.GetJailsState()

	if err == nil {
		collector.fillCurrentlyFailed(ch, jailsState.CurrentlyFailed)
		collector.fillTotalFailed(ch, jailsState.TotalFailed)
		collector.fillCurrentlyBanned(ch, jailsState.CurrentlyBanned)
		collector.fillTotalBanned(ch, jailsState.TotalBanned)
		collector.fillCountriesBanned(ch, jailsState.CountriesBanned)
	} else {
		log.Errorf("Unable to read telemetry information : %s", err)
	}
}

func (collector *Fail2banCollector) fillCurrentlyFailed(ch chan<- prometheus.Metric, failed map[string]uint) {
	for jail := range failed {
		ch <- prometheus.MustNewConstMetric(collector.jailFailedCurrent, prometheus.GaugeValue, float64(failed[jail]), jail)
	}
}

func (collector *Fail2banCollector) fillTotalFailed(ch chan<- prometheus.Metric, failed map[string]uint) {
	for jail := range failed {
		ch <- prometheus.MustNewConstMetric(collector.jailFailedTotal, prometheus.CounterValue, float64(failed[jail]), jail)
	}
}

func (collector *Fail2banCollector) fillCurrentlyBanned(ch chan<- prometheus.Metric, banned map[string]uint) {
	for jail := range banned {
		ch <- prometheus.MustNewConstMetric(collector.jailBannedCurrent, prometheus.GaugeValue, float64(banned[jail]), jail)
	}
}

func (collector *Fail2banCollector) fillTotalBanned(ch chan<- prometheus.Metric, banned map[string]uint) {
	for jail := range banned {
		ch <- prometheus.MustNewConstMetric(collector.jailBannedTotal, prometheus.CounterValue, float64(banned[jail]), jail)
	}
}

func (collector *Fail2banCollector) fillCountriesBanned(ch chan<- prometheus.Metric, countries map[string]map[string]uint) {
	for jail := range countries {
		for country := range countries[jail] {
			ch <- prometheus.MustNewConstMetric(collector.jailBannedCountryTotal, prometheus.CounterValue, float64(countries[jail][country]), jail, country)
		}
	}
}
