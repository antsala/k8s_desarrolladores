package main

import (
	"fmt"
	"log"
	"net/http"
	"net"
	"os"
)

func main() {
	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	server := http.NewServeMux()
	server.HandleFunc("/", helloContainer)

	log.Printf("Servidor escuchando en puerto %s", port)
	err := http.ListenAndServe(":"+port, server)
	log.Fatal(err)
}

func helloContainer(w http.ResponseWriter, r *http.Request) {
	log.Printf("Sirviendo request: %s", r.URL.Path)
	host, _ := os.Hostname()
	fmt.Fprintf(w, "Hola Mundo!\n")
	fmt.Fprintf(w, "Version: 1.0.0\n")
	fmt.Fprintf(w, "Hostname: %s\n", host)

	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Fprintf(w, "Dirección IP: %s\n", ipnet.IP.String())
			}
		}
	}
}
