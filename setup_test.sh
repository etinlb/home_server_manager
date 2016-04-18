#!/bin/bash
while IFS='' read -r line || [[ -n "$line" ]]; do
    echo "Touching: test/$line"
    touch "test/$line"
done < "$1"
