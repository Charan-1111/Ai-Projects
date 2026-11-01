#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
services=(
	"Ask-My-Documents"
	"LLM-Playground"
	"Semantic-Search"
	"Tool-Calling"
)
pids=()
selected=()

usage() {
	printf 'Usage: %s <service-name> [service-name ...]\n' "$0"
	printf 'Available services: %s\n' "${services[*]}"
}

service_port() {
	local service="$1"
	local config_file
	local port

	case "$service" in
		"LLM-Playground")
			printf '8000\n'
			return
			;;
		"Ask-My-Documents")
			config_file="${CONFIG_PATH:-config/config.json}"
			;;
		"Semantic-Search")
			config_file="${CONFIG_FILE_PATH:-config/local/config.json}"
			;;
		"Tool-Calling")
			config_file="${CONFIG_PATH:-config/config.json}"
			;;
	esac

	if [[ "$config_file" != /* ]]; then
		config_file="$ROOT_DIR/$service/$config_file"
	fi
	if [[ ! -f "$config_file" ]]; then
		printf 'Configuration file not found for %s: %s\n' "$service" "$config_file" >&2
		return 1
	fi

	port="$(sed -nE 's/.*"port"[[:space:]]*:[[:space:]]*"?[:]?[[:space:]]*([0-9]+)"?.*/\1/p' "$config_file" | head -n 1)"
	if [[ -z "$port" ]]; then
		printf 'Could not determine listening port for %s from %s\n' "$service" "$config_file" >&2
		return 1
	fi

	printf '%s\n' "$port"
}

if [[ $# -eq 0 ]]; then
	usage >&2
	exit 2
fi

if [[ $# -eq 1 && ( "$1" == "-h" || "$1" == "--help" ) ]]; then
	usage
	exit 0
fi

for requested in "$@"; do
	valid=false
	for service in "${services[@]}"; do
		if [[ "$requested" == "$service" ]]; then
			valid=true
			break
		fi
	done

	if [[ "$valid" != true ]]; then
		printf 'Unknown service: %s\n' "$requested" >&2
		usage >&2
		exit 2
	fi

	for service in "${selected[@]:-}"; do
		if [[ "$requested" == "$service" ]]; then
			printf 'Service specified more than once: %s\n' "$requested" >&2
			exit 2
		fi
	done
	selected+=("$requested")
done

stop_services() {
	for pid in "${pids[@]:-}"; do
		kill "$pid" 2>/dev/null || true
	done
}

trap stop_services EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

for service in "${selected[@]}"; do
	port="$(service_port "$service")"
	if nc -z 127.0.0.1 "$port" >/dev/null 2>&1; then
		printf 'Port %s is already in use; cannot start %s\n' "$port" "$service" >&2
		exit 1
	fi

	printf 'Starting %s...\n' "$service"
	"$ROOT_DIR/$service/run.sh" &
	pid="$!"
	pids+=("$pid")

	printf 'Waiting for %s to listen on port %s...\n' "$service" "$port"
	started=false
	for ((attempt = 0; attempt < 60; attempt++)); do
		if ! kill -0 "$pid" 2>/dev/null; then
			printf '%s exited before listening on port %s\n' "$service" "$port" >&2
			exit 1
		fi
		if nc -z 127.0.0.1 "$port" >/dev/null 2>&1; then
			started=true
			break
		fi
		sleep 1
	done

	if [[ "$started" != true ]]; then
		printf 'Timed out waiting for %s to listen on port %s\n' "$service" "$port" >&2
		exit 1
	fi
done

wait
