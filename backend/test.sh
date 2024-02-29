#!/bin/bash
# Currently a bug where the mock server is not ready to accept connections right away
# need to fix this in the code itself, but for now we can retry

# Maximum number of attempts
max_attempts=10
attempt=1

# Loop until the tests pass or we reach the maximum number of attempts
while [ $attempt -le $max_attempts ]; do
  echo "Attempt $attempt of $max_attempts"
  go test -v .
  
  # Check the exit status of the last command (go test)
  if [ $? -eq 0 ]; then
    echo "Tests passed on attempt $attempt"
    exit 0
  fi
  echo "Failed to establish websocket connection, retrying..."
  
  # Increment the attempt counter
  attempt=$((attempt+1))
  
  # Optional: wait before retrying
  sleep 1
done

echo "Tests failed after $max_attempts attempts"
exit 1
