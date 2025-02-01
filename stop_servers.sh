#!/bin/bash

if [ ! -f backend_pids.txt ]; then
  echo "No PID file found. Servers may not be running."
  exit 1
fi

echo "Stopping backend servers..."
while read pid; do
  echo "Killing process $pid..."
  kill -9 "$pid" 2>/dev/null
done < backend_pids.txt

# Cleanup
rm -f backend_pids.txt
echo "All servers stopped."
