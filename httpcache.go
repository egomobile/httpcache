package httpcache

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Dataset struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var DATASETS = map[string]Dataset{}

func getDataset(key string) Dataset {
	var d Dataset

	if key != "" {
		value, ok := DATASETS[key]
		if ok {
			d = value
		} else {
			d.Key = ""
			d.Value = ""
			fmt.Println("Key doesn't exist")
		}
	}

	return d
}

func getDatasets(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	key := r.URL.Query().Get("key")
	if key != "" {
		d := getDataset(key)
		if d.Key != "" {
			sendJsonResponse(w, 200, d)
		} else {
			sendResponse(w, 404, []byte("Key doesn't exist"))
		}
	} else {
		var dataSets []Dataset
		for key, value := range DATASETS {
			print(key)
			d := getDataset(value.Key)
			dataSets = append(dataSets, d)
		}
		sendJsonResponse(w, 200, dataSets)
	}
}

func putDataset(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var d Dataset

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	DATASETS[d.Key] = d

	sendResponse(w, 404, []byte("Dataset sucessfully saved!"))
}

func sendJsonResponse(w http.ResponseWriter, status int, dataSets interface{}) {
	j, err := json.Marshal(dataSets)
	if err != nil {
		fmt.Println(err)
		return
	}
	sendResponse(w, status, j)
}

func sendResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write(data)
}

func Server() {
	var d1 Dataset
	d1.Key = "foo"
	d1.Value = "bar"
	DATASETS[d1.Key] = d1

	var d2 Dataset
	d2.Key = "foo2"
	d2.Value = "bar2"
	DATASETS[d2.Key] = d2

	router := httprouter.New()
	router.GET("/datasets", getDatasets)
	router.PUT("/datasets", putDataset)

	log.Fatal(http.ListenAndServe(":80", router))
}
