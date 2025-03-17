.PHONY: run rebuild deploy install-golang

run:
	@go run cmd/api/*.go

install-golang:
	@chmod +x install_golang.sh && ./install_golang.sh

rebuild:
	@mkdir -p ../connector-production
	@go build -o ../connector-production/g-learning-connector cmd/api/*.go

logs:
	@sudo journalctl -f -u g-learning-connector

deploy:
	@git pull
	@mkdir -p ../connector-production
	@go build -o ../connector-production/g-learning-connector cmd/api/*.go
	@sudo rm -f /etc/systemd/system/g-learning-connector.service
	@sudo cp g-learning-connector.service /etc/systemd/system
	@sudo systemctl enable g-learning-connector
	@sudo systemctl restart g-learning-connector
	@sudo systemctl status g-learning-connector
