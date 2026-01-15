package server

import (
	"github.com/Dzhodddi/EcommerceAPI/product/internal/config"
	"github.com/Dzhodddi/EcommerceAPI/product/internal/product"
	"github.com/IBM/sarama"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func Run() error {
	var repository product.Repository

	producer, err := sarama.NewAsyncProducer([]string{config.BootstrapServers}, nil)
	if err != nil {
		return err
	}
	defer func(producer sarama.AsyncProducer) {
		err = producer.Close()
		if err != nil {
			log.Println(err)
		}
	}(producer)

	repository, err = product.NewElasticRepository(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer repository.Close()
	log.Println("Listening on port 8080...")
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err = http.ListenAndServe(":9090", nil)
		if err != nil {
			panic(err)
		}
	}()
	service := product.NewProductService(repository, producer)
	return product.ListenGRPC(service, 8080)
}
