#!/usr/bin/env bash

BASE_URL="http://localhost:6060"
SAMPLING_PERIOD=30

go run ./main.go & go_pid=$!
trap 'kill "$go_pid" 2>/dev/null' EXIT

sleep 3

curl_pids=()

# CPU Profile
curl -s "${BASE_URL}/debug/pprof/profile?seconds=${SAMPLING_PERIOD}" > cpu.pprof & curl_pids+=($!)
# I/O Profile
curl -s "${BASE_URL}/debug/pprof/block?seconds=${SAMPLING_PERIOD}" > block.pprof & curl_pids+=($!)
# Lock Contention Profile
curl -s "${BASE_URL}/debug/pprof/mutex?seconds=${SAMPLING_PERIOD}" > mutex.pprof & curl_pids+=($!)
# # To answer: "How many objects are we creating per request?"
curl -s "${BASE_URL}/debug/pprof/heap?seconds=${SAMPLING_PERIOD}" > heap_rate.pprof & curl_pids+=($!)

wait "${curl_pids[@]}"

# Goroutine Profile
curl -s "${BASE_URL}/debug/pprof/goroutine" > goroutine.pprof
# Save raw heap profile for analysis
curl -s "${BASE_URL}/debug/pprof/heap" > heap_raw.pprof
# The ?gc=1 parameter triggers garbage collection before sampling
curl -s "${BASE_URL}/debug/pprof/heap?gc=1" > heap_with_gc.pprof

kill "$go_pid" 2>/dev/null
