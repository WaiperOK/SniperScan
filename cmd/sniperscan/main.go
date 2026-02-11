package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/WaiperOK/SniperScan/internal/scanner"
	"github.com/WaiperOK/SniperScan/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "scan":
		runScan(os.Args[2:])
	case "serve":
		runServer(os.Args[2:])
	default:
		usage()
		os.Exit(1)
	}
}

func runScan(args []string) {
	fs := flag.NewFlagSet("scan", flag.ExitOnError)
	target := fs.String("target", "", "Target host or IP")
	portsSpec := fs.String("ports", "22,80,443", "Port spec, e.g. 22,80,8000-8010")
	timeoutMS := fs.Int("timeout-ms", 350, "Connection timeout in milliseconds")
	concurrency := fs.Int("concurrency", 200, "Concurrent workers")
	banner := fs.Bool("banner", true, "Attempt banner grabbing")
	outPath := fs.String("out", "", "Write report to JSON file")
	fs.Parse(args)

	if *target == "" {
		log.Fatal("--target is required")
	}

	ports, err := scanner.ParsePorts(*portsSpec)
	if err != nil {
		log.Fatalf("failed to parse ports: %v", err)
	}

	report, err := scanner.Scan(context.Background(), scanner.Options{
		Target:      *target,
		Ports:       ports,
		Timeout:     time.Duration(*timeoutMS) * time.Millisecond,
		Concurrency: *concurrency,
		GrabBanner:  *banner,
	})
	if err != nil {
		log.Fatalf("scan failed: %v", err)
	}

	payload, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatalf("marshal failed: %v", err)
	}

	if *outPath != "" {
		if err := os.WriteFile(*outPath, payload, 0o644); err != nil {
			log.Fatalf("write report failed: %v", err)
		}
		fmt.Printf("report written to %s\n", *outPath)
		return
	}
	fmt.Println(string(payload))
}

func runServer(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	addr := fs.String("addr", ":8097", "Bind address")
	timeoutMS := fs.Int("timeout-ms", 350, "Connection timeout in milliseconds")
	concurrency := fs.Int("concurrency", 200, "Default concurrent workers")
	fs.Parse(args)

	handler := server.NewHandler(time.Duration(*timeoutMS)*time.Millisecond, *concurrency)
	srv := &http.Server{
		Addr:              *addr,
		Handler:           handler,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Printf("SniperScan API listening on %s", *addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Println("sniperscan <command>\n\nCommands:\n  scan   Run one-shot TCP scan\n  serve  Run HTTP API server")
}
