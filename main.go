package main


import (
	"fmt"
	"net/http"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"log"
	"time"
	"context"
)

const GCP_PROJECT_ID="carrefour-ecommerce"
var client *http.Client

func main() {

	// Create an register a OpenCensus
	// Stackdriver Trace exporter.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: GCP_PROJECT_ID,
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	_, span := trace.StartSpan(context.Background(), "sample-papi")

	client = &http.Client{
		Transport: &ochttp.Transport{
			// Use Google Cloud propagation format.
			Propagation: &propagation.HTTPFormat{},
		},
	}

	http.HandleFunc("/",takeSomeTime)
	log.Fatal(http.ListenAndServe(":6060", &ochttp.Handler{}))
	defer span.End()
}

func takeSomeTime(w http.ResponseWriter, req *http.Request) {

	start := time.Now()
	_, span := trace.StartSpan(context.Background(), "sample-papi")
	reqOtherService, _ := http.NewRequest("GET", "https://www.google.es", nil)
	reqOtherService = reqOtherService.WithContext(req.Context())

	// The outgoing request will be traced with r's trace ID.
	if _, err := client.Do(reqOtherService); err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(start)
	fmt.Fprintln(w, "Time to request a second service: ",elapsed)
	span.End()

}
