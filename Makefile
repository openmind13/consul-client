
go:
	go run cmd/main.go --cfg_path=config.toml

rust:
	CFG_PATH=./config.toml cargo r