package main

import "github.com.br/sk8sta13/temperatures/internal/webserver"

func main() {
	webserver := webserver.NewWebServer("b")
	webserver.Start("b")
}
