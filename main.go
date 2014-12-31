package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
)

func main() {

	log.Printf("Version: %s", APP_VERSION)
	runtime.GOMAXPROCS(1) // don't let it go multi-process

	// make sure we can shutdown gracefully
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	doneDone := make(chan bool)
	errCh := make(chan error, 1)

	go func() {
		c, err := newCollector()
		if err != nil {
			errCh <- err
		}
		errCh <- c.collect()
	}()

	go func() {

		for {
			select {
			case err := <-errCh:
				log.Printf("Error: %v", err)
			case ex := <-sigCh:
				log.Println("Shutting down due to: ", ex)
				doneDone <- true
			default:
				// nothing to do
			}
		}

	}()
	<-doneDone

}
