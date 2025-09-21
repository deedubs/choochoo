package main

import "github.com/deedubs/choochoo/internal/server"

func main() {
	srv := server.NewWebhookServer()
	srv.Start()
}

