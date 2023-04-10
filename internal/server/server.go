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
	Restore       bool          `env:"RESTORE" envDefault:"false"`
}

var (
	address       = flag.String("a", "localhost:8080", "server address")
	restore       = flag.Bool("r", false, "restore value")
	storeInterval = flag.Duration("i", 0*time.Second, "store interval")
	storeFile     = flag.String("f", "/tmp/devops-metrics-db.json", "store file")
)

func InitConfigServer() *Server {
	cfg := &Server{}

	flag.Parse()
	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Address = os.Getenv("ADDRESS"); cfg.Address == "" {
		cfg.Address = *address
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval == "" {
		if cfg.StoreInterval, err = time.ParseDuration(envStoreInterval); err != nil {
			cfg.StoreInterval = *storeInterval
		}
	}

	if envRestore := os.Getenv("RESTORE"); envRestore == "" {
		if cfg.Restore, err = strconv.ParseBool(envRestore); err != nil {
			cfg.Restore = *restore
		}
	}

	if cfg.StoreFile = os.Getenv("STORE_FILE"); cfg.StoreFile == "" {
		cfg.StoreFile = *storeFile
	}

	return cfg
}
