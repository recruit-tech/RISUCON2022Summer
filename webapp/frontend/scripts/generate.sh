#!/bin/bash

docker run --rm \
  -v "${PWD}/../mock:/local/src" \
  -v "${PWD}/src/apiClient/__generated__:/local/dist" \
  openapitools/openapi-generator-cli generate \
    -g typescript-axios \
    -i /local/src/r-calendar.yaml \
    -o /local/dist
  
