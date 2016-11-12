// My go boiler plate code available at https://github.com/plod/go-web-boilerplate
package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	//"time"
	//"context"
	"github.com/fatih/color"
	"github.com/skratchdot/open-golang/open"
	"github.com/gorilla/handlers"
)

var port = flag.String("port", "8080", "TCP port to listen on")
//to generate default certificates
//
// go run /usr/local/go/src/crypto/tls/generate_cert.go --host=localhost
//
// depetning on personal paths
var tls = flag.Bool("tls", true, "Use TLS (https) or not")
var cert = flag.String("cert", "crt/cert.pem", "TLS Certificate")
var key = flag.String("key", "crt/key.pem", "TLS Key")
var r = mux.NewRouter()

func main() {

	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgWhite, color.BgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	flag.Parse()
	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	var address string

	if *tls == true {
		address = "https://localhost:" + *port
	} else {
		address = "http://localhost:" + *port
	}
	routing()

	loggedRouter := handlers.CombinedLoggingHandler(os.Stdout, r)

	log.Println(green("INFO"), "Starting HTTP server at", address)
	go func() {
		// service connections
		if *tls == true {
			if err := http.ListenAndServeTLS(":"+*port, *cert, *key, loggedRouter); err != nil {
				log.Printf(red("ERR")+" listen: %s\n", err)
				stopChan <- os.Interrupt
			}
		} else {
			if err := http.ListenAndServe(":"+*port, loggedRouter); err != nil {
				log.Printf(red("ERR")+" listen: %s\n", err)
				stopChan <- os.Interrupt
			}
		}
	}()

	//useful in development cycle to open browser on running
	log.Println(green("INFO"), "Opening browser...")
	err := open.Run(address)
	if(err != nil){
		log.Println(red("ERR"), "unable to open browser", err)
	}

	<-stopChan // wait for SIGINT
	log.Println(yellow("SWARN"), "Shutting down server...")

	// shut down gracefully, but wait no longer than 5 seconds before halting //need go http connection draining added to include
	//ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//http.Shutdown(ctx)

	log.Println(red("Server gracefully stopped"))
}

