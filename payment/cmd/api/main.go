package main

import (
	"github.com/Dzhodddi/EcommerceAPI/payment/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"

	"github.com/Dzhodddi/EcommerceAPI/payment/internal"
	"github.com/IBM/sarama"
	"github.com/tinrab/retry"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var repository internal.Repository

	db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	repository, err = internal.NewPostgresRepository(db)
	if err != nil {
		panic(err)
	}

	// Setup Kafka consumer
	var consumer sarama.Consumer
	if config.KafkaBrokers != "" {
		retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
			kafkaConfig := sarama.NewConfig()
			kafkaConfig.Consumer.Return.Errors = true

			consumer, err = sarama.NewConsumer([]string{config.KafkaBrokers}, kafkaConfig)
			if err != nil {
				log.Printf("Failed to create Kafka consumer: %v", err)
			}
			return
		})
	}

	dodoClient := internal.NewDodoClient(config.DodoAPIKEY, config.DodoTestMode)
	service := internal.NewPaymentService(dodoClient, repository)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err = http.ListenAndServe(":9090", nil)
		if err != nil {
			panic(err)
		}
	}()
	log.Fatal(internal.StartServers(service, consumer, config.OrderServiceURL, config.GrpcPort, config.WebhookPort))
}
