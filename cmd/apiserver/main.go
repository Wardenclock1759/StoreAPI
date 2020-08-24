package main

import (
	"github.com/Wardenclock1759/StoreAPI/internal/apiserver"
	"log"
)

func main() {
	s := apiserver.New()
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
