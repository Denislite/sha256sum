package utils

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func CreateDocs() {
	flag.Usage = func() {
		_, err := fmt.Fprintln(os.Stderr, "Options for tool:")
		if err != nil {
			return
		}
		flag.VisitAll(func(f *flag.Flag) {
			_, err := fmt.Fprintf(os.Stderr, "-%v %v  \n", f.Name, f.Usage)
			if err != nil {
				log.Println(ErrorDocs)
				return
			}
		})
	}
}

func CheckSignal(signals chan os.Signal) {
	signal.Notify(signals, os.Interrupt)
	go func() {
		for sig := range signals {
			log.Printf("request canceled by signal %d \n", sig)
			os.Exit(1)
		}
	}()
}
