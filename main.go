package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

var actions = []string{"move", "eat", "load", "unload"}
var directions = []string{"up", "down", "right", "left"}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		data, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		//Hive object from request payload
		var hive Hive
		json.Unmarshal(data, &hive)

		//Loop through ants and give orders
		orders := make(map[int]Order)
		for antID, _ := range hive.Ants {
			orders[antID] = Order{
				Act: actions[rand.Intn(4)],
				Dir: directions[rand.Intn(4)],
			}
		}
		response, _ := json.Marshal(orders)

		log.Println(string(response))
		w.Write(response)
		// json format sample:
		// {"1":{"act":"load","dir":"down"},"17":{"act":"load","dir":"up"}}
	})

	err := http.ListenAndServe(":7070", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type Hive struct {
	Ants map[int]*json.RawMessage
}

type Order struct {
	Act string `json:"act"`
	Dir string `json:"dir"`
}
