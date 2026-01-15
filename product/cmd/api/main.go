package main

import "github.com/Dzhodddi/EcommerceAPI/product/internal/server"

func main() {
	err := server.Run()
	if err != nil {
		panic(err)
	}
}
