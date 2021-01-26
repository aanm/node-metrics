package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const Namespace = "node-metrics"

var (
	Registry *prometheus.Registry

	Address string
)

type bpfCollector struct {
	memTotal     *prometheus.Desc
	memAvailable *prometheus.Desc
}

func newbpfCollector() *bpfCollector {
	return &bpfCollector{
		memTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "mem_total_kilobytes"),
			"Total memory",
			nil, nil,
		),
		memAvailable: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "mem_available_kilobytes"),
			"Total memory available",
			nil, nil,
		),
	}
}

func (s *bpfCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- s.memTotal
	ch <- s.memAvailable
}

func getMemoryUsage(typ string) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", "/usr/bin/get_metric.sh", typ)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("unable to get meminfo output: %w: %s", err, string(out))
	}

	return strconv.ParseUint(strings.TrimSpace(string(out)), 10, 64)
}

func (s *bpfCollector) Collect(ch chan<- prometheus.Metric) {
	mapMem, err := getMemoryUsage("MemTotal:")
	if err != nil {
		log.Println("Error while getting total memory usage:", err)
	} else {
		ch <- prometheus.MustNewConstMetric(
			s.memTotal,
			prometheus.GaugeValue,
			float64(mapMem),
		)
	}

	progMem, err := getMemoryUsage("MemAvailable:")
	if err != nil {
		log.Println("Error while getting available memory usage:", err)
	} else {
		ch <- prometheus.MustNewConstMetric(
			s.memAvailable,
			prometheus.GaugeValue,
			float64(progMem),
		)
	}
}

func main() {
	Registry = prometheus.NewPedanticRegistry()
	Registry.MustRegister(newbpfCollector())
	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/node-metrics", promhttp.HandlerFor(Registry, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(os.Getenv("PROMETHEUS_ADDR"), nil))
}
