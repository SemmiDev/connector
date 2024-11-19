.PHONY: run rebuild deploy install-golang

run:
	@go run cmd/*.go

install-golang:
	@chmod +x install_golang.sh && ./install_golang.sh

rebuild:
	@mkdir -p ../connector-production
	@go build -o ../connector-production/g-learning-connector cmd/*.go

deploy:
	@git pull
	@sudo rm /etc/systemd/system/g-learning-connector.service
	@sudo mv g-learning-connector.service /etc/systemd/system
	@sudo systemctl enable g-learning-connector
	@sudo systemctl restart g-learning-connector
