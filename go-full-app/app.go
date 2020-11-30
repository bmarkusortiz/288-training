package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

var lang = flag.String("lang", "en", "run app with language support - default is english")
var alive = flag.Bool("alive", true, "Condition for the app to return a healthy or un healthy response")
var started = time.Now()

func main() {
	var port = "8080"
	var ip = "default"
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthz", healtzHandler)
	http.HandleFunc("/ready", readyHandler)
	http.HandleFunc("/flip", flipHandler)

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			ip = ipnet.IP.String()
		}
	}
	fmt.Printf("Starting server on PORT: %[1]v and IP: %[2]v \n", port, ip)
	http.ListenAndServe(":"+port, nil)
}

func rootHandler(response http.ResponseWriter, request *http.Request) {

	flag.Parse()

	switch *lang {
	case "en":
		fmt.Fprintf(response, "Hello %s!. Welcome!\n", request.URL.Path[1:])
	case "es":
		fmt.Fprintf(response, "Hola %s!. Bienvenido!\n", request.URL.Path[1:])
	default:
		fmt.Fprintf(response, "Error! unknown lang option -> %s\n", *lang)
	}
}

func healtzHandler(response http.ResponseWriter, request *http.Request) {
	if *alive {
		fmt.Println("ping /healthz => pong [healthy]")
		fmt.Fprintf(response, "Ok\n")
	} else {
		fmt.Println("ping /healthz => pong [unhealthy]")
		http.Error(response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(response, "Error!. App not healthy!\n")
	}
}

func readyHandler(response http.ResponseWriter, request *http.Request) {
	now := time.Now()
	diff := now.Sub(started)
	if int(diff.Seconds()) > 30 {
		fmt.Println("ping /ready => pong [ready]")
		fmt.Fprintf(response, "Ready for service requests...\n")
	} else {
		fmt.Println("ping /ready => pong [notready]")
		http.Error(response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(response, "Error! Service not ready for requests...\n")
	}
}

func flipHandler(response http.ResponseWriter, request *http.Request) {
	var action = request.URL.Query()["action"]
	if action[0] == "kill" {
		fmt.Println("Received kill request. Changing app state to unhealthy...")
		*alive = false
		fmt.Fprintf(response, "Switched app state to unhealthy...\n")
	} else if action[0] == "revive" {
		fmt.Println("Received revive request. Changing app state to healthy...")
		*alive = true
		fmt.Fprintf(response, "Switched app state to healthy...\n")
	}
}
