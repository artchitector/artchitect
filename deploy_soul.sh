go build -o ../../bin/soul -C ./services/soul cmd/main.go
pm2 start --name soul