#!/bin/bash
while IFS='' read -r line || [[ -n "$line" ]]; do
    echo "Touching: test_data/$line"
    touch "test_data/$line"
done < "$1"
