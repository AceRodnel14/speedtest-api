package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/julienschmidt/httprouter"
)

// const (
// 	layoutISO = "2006-01-02T15:04:05Z"
// )

type SpeedtestResult struct {
	// TimeStamp time.Time `json:"timestamp"`
	Ping     Latency `json:"ping"`
	Download Stats   `json:"download"`
	Upload   Stats   `json:"upload"`
}

type Latency struct {
	Jitter  float64 `json:"jitter"`
	Latency float64 `json:"latency"`
}

type Stats struct {
	Bandwidth float64 `json:"bandwidth"`
}

type outputData struct {
	// TimeStamp     string  `json:"timestamp"`
	Jitter        float64 `json:"jitter"`
	Latency       float64 `json:"latency"`
	DownBandwidth float64 `json:"download_speed"`
	UpBandwidth   float64 `json:"upload_speed"`
}

func main() {

	router := httprouter.New()
	router.GET("/metrics", speedtestExport("prom"))
	router.GET("/metrics/json", speedtestExport("json"))

	log.Fatal(http.ListenAndServe(":9001", router))

}

// func changeFormat(t time.Time) string {
// 	output := t.Format(layoutISO)
// 	return output
// }

func parseJson(path string) (result SpeedtestResult) {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println("File missing")
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &result)

	return result
}

func performSpeedtest() (result SpeedtestResult) {
	cmd := exec.Command("/bin/sh", "/exec/run")
	path := "/resources/report.json"
	// if fileExists(path) {
	// 	result = parseJson(path)
	// 	loc, _ := time.LoadLocation("UTC")
	// 	t := time.Now().In(loc)

	// 	if t.Sub(result.TimeStamp).Seconds() > 600 {
	// 		cmd.Run()
	// 		result = parseJson(path)
	// 		return result
	// 	}

	// 	if t.Sub(result.TimeStamp).Seconds() > 30 {
	// 		cmd.Start()
	// 		result = parseJson(path)
	// 		return result
	// 	}
	// 	result = parseJson(path)
	// 	return result

	// }
	cmd.Run()
	result = parseJson(path)
	return result
}

func printData(result SpeedtestResult) outputData {
	list := outputData{
		// TimeStamp:     changeFormat(result.TimeStamp),
		Jitter:        result.Ping.Jitter,
		Latency:       result.Ping.Latency,
		DownBandwidth: ((result.Download.Bandwidth * 8) / 1000000),
		UpBandwidth:   ((result.Upload.Bandwidth * 8) / 1000000),
	}
	return list
}

func speedtestExport(format string) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		result := performSpeedtest()
		list := printData(result)

		if format == "prom" {
			w.Header().Set("Content-Type", " text/plain; charset=utf-8")
			output := "jitter %.2f\n" +
				"latency %.2f\n" +
				"download_speed %.2f\n" +
				"upload_speed %.2f"
			fmt.Fprintf(w, output, list.Jitter, list.Latency, list.DownBandwidth, list.UpBandwidth)
		}
		if format == "json" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(list)
		}
	}
}

// func fileExists(filename string) bool {
// 	_, err := os.Stat(filename)
// 	return !os.IsNotExist(err)
// }
