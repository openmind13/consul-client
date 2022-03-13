
run:
	go run cmd/main.go --config=config.toml

race:
	go run -race cmd/main.go --config=config.toml