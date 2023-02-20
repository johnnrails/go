#!/bin/bash
for i in {1..1000}; do
  curl -X GET http://localhost:9090/metrics
  #  -H "Content-Type: application/json" -d '{"name": "rails"}'  
done
