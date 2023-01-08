package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"remainder_app_2/controller"
	"strconv"

	"github.com/gorilla/mux"
	consulapi "github.com/hashicorp/consul/api"
)

func main() {
	serviceRegistryWithConsul()
	fmt.Println("Starting the application...")
	router := mux.NewRouter()
	router.HandleFunc("/tmrevent", controller.Findtmrevents).Methods("GET")
	router.HandleFunc("/check", check)
	http.ListenAndServe(":12345", router)
}

func serviceRegistryWithConsul() {
	config := consulapi.DefaultConfig()
	fmt.Print(config)
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Println(err)
	}
	serviceID := "go_micro_2"
	port, _ := strconv.Atoi(getPort()[1:len(getPort())])
	address := getHostname()

	registration := &consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    "go_micro_2",
		Port:    port,
		Address: address,
		Check: &consulapi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%v/check", address, port),
			Interval: "10s",
			Timeout:  "30s",
		},
	}

	regiErr := consul.Agent().ServiceRegister(registration)

	if regiErr != nil {
		log.Printf("Failed to register service: %s:%v ", address, port)
	} else {
		log.Printf("successfully register service: %s:%v", address, port)
	}
}

func getPort() (port string) {
	port = os.Getenv("PORT")
	if len(port) == 0 {
		port = "12345"
	}
	port = ":" + port
	return
}

func getHostname() (hostname string) {
	hostname, _ = os.Hostname()
	return
}

func check(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Consul check")
}
