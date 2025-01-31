#!/usr/bin/env bash
set -Eeuo pipefail

protoc \
		-I . \
		--go_out=paths=source_relative:. \
		--fastmarshal_out=apiversion=v2,enableunsafedecode=true,paths=source_relative:. \
    ./types.proto
