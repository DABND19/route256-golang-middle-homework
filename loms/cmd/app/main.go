package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const port = ":8081"

type Item struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type RequestPayload struct {
	User  int64  `json:"user"`
	Items []Item `json:"items"`
}

type ResponsePayload struct {
	OrderID int64 `json:"orderID"`
}

type ErrorResponsePayload struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/createOrder", func(w http.ResponseWriter, r *http.Request) {
		reqPayload := RequestPayload{}
		err := json.NewDecoder(r.Body).Decode(&reqPayload)
		if err != nil {
			resPayload := ErrorResponsePayload{
				Message: "Invalid request payload: " + err.Error(),
			}

			resBody, err := json.Marshal(resPayload)
			if err != nil {
				log.Fatal("Failed to encode error response payload:", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(resBody)
			return
		}

		resPayload := ResponsePayload{OrderID: 1}
		resBody, err := json.Marshal(resPayload)
		if err != nil {
			log.Fatal("Failed to encode response payload:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resBody)
	})

	log.Println("Starting a server...")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Couldn't start a server.")
	}
}
