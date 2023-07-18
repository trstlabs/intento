pkill -f "docker-compose .*trst.* logs" | true
pkill -f "/bin/bash.*create_logs.sh" | true
pkill -f "tail .*.log" | true
