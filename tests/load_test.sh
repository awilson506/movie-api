#!/bin/bash

echo "Load testing production company details rpc"
for i in `seq 1 2000`; do curl -s http://localhost:8080/production-company-details -o /dev/null; done
echo "Load test complete"