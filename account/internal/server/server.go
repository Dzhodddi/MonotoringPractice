package server

import (
	"github.com/Dzhodddi/EcommerceAPI/account/internal/accounts"
	"github.com/Dzhodddi/EcommerceAPI/account/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
	"log"
	"net/http"
)

func Run() error {
	var repository accounts.Repository
	db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "accounts",
		RefreshInterval: 15,
		StartServer:     false,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.Postgres{},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}
	repository, err = accounts.NewPostgresRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	defer repository.Close()
	log.Println("Listening on port 8080...")
	service := accounts.NewService(repository)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err = http.ListenAndServe(":9090", nil); err != nil {
			panic(err)
		}
	}()
	return accounts.ListenGRPC(service, 8080)
}
