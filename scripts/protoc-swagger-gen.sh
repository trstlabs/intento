#!/usr/bin/env bash

set -eo pipefail

mkdir -p ./tmp-swagger-gen
cd proto
proto_dirs=$(find ../proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do

    buf generate --template buf.gen.swagger.yaml $query_file

done

cd ..
# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.json` for merging
swagger-combine ./client/docs/config.json -o ./client/docs/swagger-ui/swagger.yaml -f yaml --continueOnConflictingPaths true --includeDefinitions true

# clean swagger files
rm -rf ./tmp-swagger-gen