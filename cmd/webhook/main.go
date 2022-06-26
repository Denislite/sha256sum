package main

import (
	"encoding/json"
	"fmt"
	"html"
	admissionv1 "k8s.io/api/admission/v1"
	"log"
	"net/http"
	"sha256sum/pkg/webhook"
	"time"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "hello %q", html.EscapeString(r.URL.Path))
	if err != nil {
		log.Println(err)
	}
}

func handleMutate(w http.ResponseWriter, r *http.Request) {
	admReview, err := webhook.AdmissionReviewFromRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("error getting admission review from request: %v", err)
		return
	}

	admResp, err := webhook.AdmissionResponseFromReview(admReview)

	if err != nil {
		w.WriteHeader(400)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Println(err)
		}
		return
	}

	// the final response will be another admission review
	var admissionReviewResponse admissionv1.AdmissionReview
	admissionReviewResponse.Response = admResp
	admissionReviewResponse.SetGroupVersionKind(admReview.GroupVersionKind())
	admissionReviewResponse.Response.UID = admReview.Request.UID

	resp, err := json.Marshal(admissionReviewResponse)
	if err != nil {
		msg := fmt.Errorf("error marshaling response: %v", err)
		log.Println(msg)
		w.WriteHeader(500)
		_, err := w.Write([]byte(msg.Error()))
		if err != nil {
			log.Println(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	log.Printf("allowing pod as %v", string(resp))
	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	log.Println("starting server...")

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/mutate", handleMutate)

	s := &http.Server{
		Addr:           ":8443",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1048576
	}

	log.Fatal(s.ListenAndServeTLS("./ssl/tcpdump.pem", "./ssl/tcpdump.key"))
}
