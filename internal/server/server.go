package server

import (
	"flag"
	"github.com/caarlos0/env/v7"
	"log"
	"os"
	"strconv"
	"time"
)

type Server struct {
	Address       string        `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

var (
	address       = flag.String("a", "localhost:8080", "server address")
	restore       = flag.Bool("r", true, "restore value")
	storeInterval = flag.Duration("i", 300*time.Second, "store interval")
	storeFile     = flag.String("f", "/tmp/devops-metrics-db.json", "store file")
)

var Config Server

func init() {
	Config = Server{}

	flag.Parse()
	err := env.Parse(&Config)
	if err != nil {
		log.Fatal(err)
	}

	if Config.Address = os.Getenv("ADDRESS"); Config.Address == "" {
		Config.Address = *address
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		if Config.StoreInterval, err = time.ParseDuration(envStoreInterval); err != nil {
			Config.StoreInterval = *storeInterval
		}
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		if Config.Restore, err = strconv.ParseBool(envRestore); err != nil {
			Config.Restore = *restore
		}
	}
	if Config.StoreFile = os.Getenv("STORE_FILE"); Config.StoreFile == "" {
		Config.StoreFile = *storeFile
	}
}
