#Устанавливаем переменную среды, где находится файл нфстроек local.yaml
path:=CONFIG_PATH=./config/local.yaml

.PHONY: server
server:
				@echo "Running server"
				$(path) go run ./cmd/myrestapi/main.go
				open http://localhost:8082/