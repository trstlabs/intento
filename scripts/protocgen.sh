#!/usr/bin/env bash

set -eo pipefail

echo "Generating gogo proto code"
cd proto
proto_dirs=$(find ./ -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep go_package "$file" &>/dev/null; then
      buf generate --template buf.gen.gogo.yaml "$file"
    fi
  done
done

cd ..

# Move proto files to the right places
cp -r github.com/trstlabs/intento/* ./
rm -rf github.com

# go mod tidy