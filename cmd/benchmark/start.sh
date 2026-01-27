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

curl_pids=()

# Goroutine Profile
curl -s "${BASE_URL}/debug/pprof/goroutine" > goroutine.pprof & curl_pids+=($!)
# Save inuse_objects (currently live objects)
curl -s "${BASE_URL}/debug/pprof/heap" > heap_inuse_objects.pprof & curl_pids+=($!)
# Save inuse_space (bytes currently in use)
curl -s "${BASE_URL}/debug/pprof/heap" > heap_inuse_space.pprof & curl_pids+=($!)

wait "${curl_pids[@]}"

# The ?gc=1 parameter triggers garbage collection before sampling
# Save alloc_objects (total allocations during profile)
curl -s "${BASE_URL}/debug/pprof/heap?gc=1" > heap_alloc_objects.pprof
# Save alloc_space (bytes allocated)
curl -s "${BASE_URL}/debug/pprof/heap?gc=1" > heap_alloc_space.pprof

kill "$go_pid" 2>/dev/null
