#!/usr/bin/env bash
#
# request_loop.sh — load generator for the projeto-korp endpoint. Every
# INTERVAL seconds it launches BATCH_SIZE requests in parallel, for DURATION
# seconds, producing request volume to observe in Grafana/Prometheus.
set -euo pipefail

readonly URL="${URL:-http://localhost:80/projeto-korp}"
readonly BATCH_SIZE="${BATCH_SIZE:-100}"
readonly INTERVAL="${INTERVAL:-5}"
readonly DURATION="${DURATION:-50}"

# send_batch fires BATCH_SIZE requests in parallel and waits for them all.
send_batch() {
	local request
	for ((request = 1; request <= BATCH_SIZE; request++)); do
		curl -s -o /dev/null "${URL}" &
	done
	wait
}

echo "Launching ${BATCH_SIZE} requests every ${INTERVAL}s for ${DURATION}s to ${URL}..."

deadline=$((SECONDS + DURATION))
batch=0
while ((SECONDS < deadline)); do
	batch=$((batch + 1))
	send_batch
	echo "batch ${batch}: ${BATCH_SIZE} requests sent (elapsed ${SECONDS}s)"

	# Skip the wait after the final batch so the run ends on time.
	if ((SECONDS < deadline)); then
		sleep "${INTERVAL}"
	fi
done

echo "Done: ${batch} batches, $((batch * BATCH_SIZE)) requests total."
