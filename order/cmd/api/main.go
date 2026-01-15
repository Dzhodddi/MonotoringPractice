package main

import (
	"github.com/Dzhodddi/EcommerceAPI/order/internal"
	"github.com/Dzhodddi/EcommerceAPI/order/internal/config"
	"github.com/IBM/sarama"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	var repository internal.Repository

	producer, err := sarama.NewAsyncProducer([]string{config.BootstrapServers}, nil)
	if err != nil {
		log.Println(err)
	}
	defer func(producer sarama.AsyncProducer) {
		err := producer.Close()
		if err != nil {
			log.Println(err)
		}
	}(producer)

	db, err := gorm.Open(postgres.Open(config.DatabaseUrl), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	repository, err = internal.NewPostgresRepository(db)
	if err != nil {
		panic(err)
	}
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err = http.ListenAndServe(":9090", nil)
		if err != nil {
			panic(err)
		}
	}()
	defer repository.Close()
	log.Println("Listening on port 8080...")
	service := internal.NewOrderService(repository, producer)
	log.Fatal(internal.ListenGRPC(service, config.AccountUrl, config.ProductUrl, 8080))
}
