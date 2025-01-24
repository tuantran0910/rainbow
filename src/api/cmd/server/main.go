package main

import (
	"github.com/tuantran0910/rainbow/internal/routes"
)

func main() {
	// Create a new router
	r := routes.NewRouter()

	// Run the server
	r.Run(":5000")
}
