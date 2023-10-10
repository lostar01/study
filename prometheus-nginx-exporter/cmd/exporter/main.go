package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	exporter "github.com/lostar01/prometheus-nginx-exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		targetHost = flag.String("target.host", "localhost", "nginx address with basic_status page")
		targetPort = flag.Int("target.port", 8080, "nginx port with basic_status page")
		targetPath = flag.String("target.path", "/status", "URL path to scrap metrics")
		promPort   = flag.Int("prom.port", 9150, "port to expose prometheus metrics")
	)
	flag.Parse()

	uri := fmt.Sprintf("http://%s:%d%s", *targetHost, *targetPort, *targetPath)

	// Called on each collector.Collect
	basicStats := func() ([]exporter.NginxStats, error) {
		var netClient = &http.Client{
			Timeout: time.Second * 10,
		}
		resp, err := netClient.Get(uri)
		if err != nil {
			log.Fatalf("netClient.Get failed %s: %s", uri, err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("io.ReadAll failed: %s", err)
		}
		r := bytes.NewReader(body)
		return exporter.ScanBasicStats(r)
	}

	// Make Prometheus client aware of our collectors
	bc := exporter.NewBasicCollector(basicStats)

	reg := prometheus.NewRegistry()
	// reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(bc)

	mux := http.NewServeMux()
	promHander := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", promHander)

	// Start listening for HTTP connections.
	port := fmt.Sprintf(":%d", *promPort)
	log.Printf("starting nginx exporter on %q/metrics", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("cannot start nginx exporter: %s", err)
	}

}
