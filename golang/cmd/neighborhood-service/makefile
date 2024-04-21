run-local:
	export ENV="local" && go run main.go

test-cov:
	go test -coverprofile=coverage.out ./service
	go tool cover -html=coverage.out