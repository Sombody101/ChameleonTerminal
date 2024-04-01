#!/bin/bash

output_path="build/gecko"

if ! hash go; then
    echo "Go is not installed!"
    exit 1
fi

go build -o "$output_path"
echo "Built to $output_path"