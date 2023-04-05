package prom

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/fail2ban-exporter/pkg/core"
)

// Fail2banExporter object to run exporter server and expose metrics
type Fail2banExporter struct {
	Address string
	Port    int
}

// Run method to run http exporter server
func (exporter *Fail2banExporter) Run(connector core.Fail2banConnector) {
	log.Info(fmt.Sprintf("Beginning to serve on port : %d", exporter.Port))

	prometheus.MustRegister(NewFail2banCollector(connector))
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", exporter.Address, exporter.Port), nil))
}
