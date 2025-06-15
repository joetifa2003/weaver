#!/bin/bash

# Check if correct number of arguments provided
if [ $# -lt 2 ] || [ $# -gt 3 ]; then
    echo "Usage: $0 <start_command> <port> [--pm2]"
    echo "Example: $0 'node main.js' 3000"
    echo "Example: $0 'npm start' 8080 --pm2"
    exit 1
fi

START_COMMAND="$1"
PORT="$2"
PM2_MODE=false

# Check for optional --pm2 argument
if [ $# -eq 3 ]; then
    if [ "$3" == "--pm2" ]; then
        PM2_MODE=true
    else
        echo "Invalid argument: $3"
        echo "Usage: $0 <start_command> <port> [--pm2]"
        exit 1
    fi
fi

# Function to stop the server by port
stop_server() {
    echo "Stopping server on port $PORT..."
  
    # Find and kill processes using the specified port
    PIDS=$(lsof -ti:$PORT 2>/dev/null)
  
    if [ -n "$PIDS" ]; then
        echo "Found processes using port $PORT: $PIDS"
        for PID in $PIDS; do
            echo "Killing process $PID..."
            kill $PID 2>/dev/null
        done
      
        # Wait a moment for graceful shutdown
        sleep 2
      
        # Force kill if still running
        REMAINING_PIDS=$(lsof -ti:$PORT 2>/dev/null)
        if [ -n "$REMAINING_PIDS" ]; then
            echo "Force killing remaining processes: $REMAINING_PIDS"
            for PID in $REMAINING_PIDS; do
                kill -9 $PID 2>/dev/null
            done
        fi
        echo "Server stopped."
    else
        echo "No processes found using port $PORT"
    fi
}

# Function to stop memory monitor gracefully
stop_memory_monitor() {
    if [ -n "$MEMORY_MONITOR_PID" ]; then
        echo "Stopping memory monitor (PID: $MEMORY_MONITOR_PID)..."
        # Send SIGINT to allow graceful shutdown and JSON saving
        kill -INT $MEMORY_MONITOR_PID 2>/dev/null
        # Wait for it to finish saving
        wait $MEMORY_MONITOR_PID 2>/dev/null
        echo "Memory monitor stopped and data saved."
    fi
}

# Start the server in the background
echo "Starting server with command: $START_COMMAND"
eval "$START_COMMAND" &

# Trap signals for cleanup - memory monitor first to save data
trap 'stop_memory_monitor; stop_server' EXIT INT TERM

# Wait for the server to start up and bind to the port
echo "Waiting for server to start on port $PORT..."
for i in {1..10}; do
    if lsof -i:$PORT >/dev/null 2>&1; then
        echo "Server is running on port $PORT"
        break
    fi
    if [ $i -eq 10 ]; then
        echo "Error: Server failed to start on port $PORT after 10 seconds"
        exit 1
    fi
    sleep 1
done

# Start memory monitoring in the background
echo "Starting memory monitoring..."
if [ "$PM2_MODE" == true ]; then
    echo "Running memory monitor in PM2 mode"
    python3 ../memory_monitor.py $PORT mem.json --pm2 &
else
    echo "Running memory monitor in port mode"
    python3 ../memory_monitor.py $PORT mem.json &
fi
MEMORY_MONITOR_PID=$!

# Wait a moment for memory monitor to start
sleep 2

# Run the plow command with the specified port
echo "Running plow load test on http://localhost:$PORT/user/1..."
plow -c 100 --json -d10s "http://localhost:$PORT/user/1" > stats.json

# Check if plow command was successful
if [ $? -eq 0 ]; then
    echo "Load test completed successfully. Results saved to stats.json"
else
    echo "Load test failed"
fi

# Give memory monitor a moment to collect final data
sleep 1

echo "Script completed"
