package main

import (
	"github.com/Dzhodddi/EcommerceAPI/gateway/config"
	"github.com/Dzhodddi/EcommerceAPI/gateway/graph"
	"github.com/Dzhodddi/EcommerceAPI/pkg/middleware"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	server, err := graph.NewGraphQLServer(config.AccountUrl, config.ProductUrl, config.OrderUrl, config.PaymentUrl, config.RecommenderUrl)
	if err != nil {
		log.Fatal(err)
	}

	srv := handler.New(server.ToExecutableSchema())
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	engine := gin.Default()

	engine.Use(middleware.GinContextToContextMiddleware())

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	engine.POST("/graphql",
		middleware.AuthorizeJWT(),
		gin.WrapH(srv),
	)
	engine.GET("/playground", gin.WrapH(playground.Handler("Playground", "/graphql")))
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	log.Fatal(engine.Run(":8080"))
}
