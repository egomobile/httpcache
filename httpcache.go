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

var dataSets = map[string]Dataset{}

func datasetsAsJson(key string) Dataset {
	keys := make([]string, 0, len(dataSets))
	values := make([]Dataset, 0, len(dataSets))

	for k, v := range dataSets {
		keys = append(keys, k)
		values = append(values, v)
	}

	// jsonDatasets, err := json.Marshal(values)

	var d Dataset

	if key != "" {
		value, ok := dataSets[key]
		if ok {
			d = value
		} else {
			fmt.Println("Value doesn't exist")
		}
	}

	return d

	/* if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return nil
	} else {
		return jsonDatasets
	} */
}

func getDatasets(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	key := r.URL.Query().Get("key")
	if key != "" {
		value := datasetsAsJson(key)
		// if value != nil {
		b, err := json.Marshal(value)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))

		sendResponse(w, 200, b)
		// } else {
		// sendResponse(w, 404, []byte("Value does not exist for this key!"))
		//}
	//} else {
		// sendResponse(w, 200, []byte(datasetsAsJson("")))
	//}
}

/* func getStringForDataset(key string, value []byte) string {
	return "{\"key\": \"" + string(key) + "\", \"value\": \"" + string(value) + "\"}"
}

func getStringForDatasets() string {
	dataSetsStr := "["
	dataSetsLength := len(dataSets)
	println(dataSetsLength)
	i := 0
	for key, value := range dataSets {
		dataSetsStr += getStringForDataset(key, []byte(value))
		if i < dataSetsLength-1 {
			dataSetsStr += ","
		}
		i++
	}
	dataSetsStr += "]"
	return dataSetsStr
}*/

func putDataset(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var d Dataset

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataSets[d.Key] = d

	// sendResponse(w, 200, []byte(datasetsAsJson("")))
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
	dataSets[d1.Key] = d1

	var d2 Dataset
	d2.Key = "foo2"
	d2.Value = "bar2"
	dataSets[d2.Key] = d2

	router := httprouter.New()
	router.GET("/datasets", getDatasets)
	router.PUT("/datasets", putDataset)

	log.Fatal(http.ListenAndServe(":80", router))
}
