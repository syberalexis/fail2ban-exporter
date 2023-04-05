package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/syberalexis/fail2ban-exporter/pkg/database"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/fail2ban-exporter/pkg/core"
	"github.com/syberalexis/fail2ban-exporter/pkg/prom"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// Default variables
	version                  = "dev"
	defaultPort              = 9902
	defaultAddress           = "0.0.0.0"
	defaultLibFolder         = "/var/lib/fail2ban-exporter"
	defaultDbPath            = defaultLibFolder + "/fail2ban-exporter.sqlite3"
	defaultFail2banDb        = "/var/lib/fail2ban/fail2ban.sqlite3"
	defaultCountryServiceURL = "https://api.country.is"

	app        = kingpin.New(filepath.Base(os.Args[0]), "")
	appVersion = app.Version(version)
	help       = app.HelpFlag.Short('h')
	debug      = app.Flag("debug", "Enable debug mode.").Bool()

	address = app.Flag("address", "Listen address.").Default(fmt.Sprintf("%s", defaultAddress)).Short('a').String()
	port    = app.Flag("port", "Listen port.").Default(fmt.Sprintf("%d", defaultPort)).Short('p').Int()

	// socket       = app.Flag("socket", "Fail2ban socket path.").Default(defaultSocket).Short('s').String()
	localisation      = app.Flag("localisation", "Locate IPs by countries.").Default("false").Short('l').Bool()
	countryServiceUrl = app.Flag("country-service-url", "Service URL to get country from IP.").Default(defaultCountryServiceURL).String()

	fail2banDbPath = app.Flag("fail2ban-db-path", "Fail2ban Database path.").Default(fmt.Sprintf("%s", defaultFail2banDb)).String()
)

// Linky-exporter command main
func main() {
	// Main action
	app.Action(func(c *kingpin.ParseContext) error { run(); return nil })

	// Parsing
	args, err := app.Parse(os.Args[1:])

	if err != nil {
		log.Error(errors.Wrapf(err, "Error parsing commandline arguments"))
		app.Usage(os.Args[1:])
		os.Exit(2)
	} else {
		kingpin.MustParse(args, err)
	}
}

// Main run function
func run() {
	if debug != nil && *debug {
		log.SetLevel(log.DebugLevel)
		log.Info("Debug mode enabled !")
	}

	// Checks before running
	// TODO check fail2ban db file
	if _, err := os.Stat(defaultLibFolder); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(defaultLibFolder, os.ModePerm)
		log.Fatal(err)
	}

	// Initializze Database
	database.InitFail2banDB(*fail2banDbPath)
	if *localisation {
		database.InitIpDB(defaultDbPath)
	}

	// Parse parameters
	connector := core.Fail2banConnector{CountryEnabled: *localisation, CountryServiceUrl: *countryServiceUrl}

	// Run exporter
	exporter := prom.Fail2banExporter{Address: *address, Port: *port}
	exporter.Run(connector)
}
