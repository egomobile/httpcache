// This file is part of the ego-cli distribution.
// Copyright (c) Next.e.GO Mobile SE, Aachen, Germany (https://e-go-mobile.com/)
//
// ego-cli is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as
// published by the Free Software Foundation, version 3.
//
// ego-cli is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

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
			sendResponse(w, 404, []byte("key doesn't exist"))
		}
	} else {
		var dataSets []Dataset
		for _, value := range DATASETS {
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

	sendResponse(w, 404, []byte("dataset sucessfully saved"))
}

func sendJsonResponse(w http.ResponseWriter, status int, dataSets interface{}) {
	j, err := json.Marshal(dataSets)
	if err != nil {
		fmt.Println(err)
		sendResponse(w, 500, []byte("an error occured while marshalling the response data"))
	}
	sendResponse(w, status, j)
}

func sendResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write(data)
}

func Server(port string, name string) {
	router := httprouter.New()
	router.GET("/datasets", getDatasets)
	router.PUT("/datasets", putDataset)

	fmt.Println("Server", name, "listening on port", port, "...")
	log.Fatal(http.ListenAndServe(":" + port, router))
}
