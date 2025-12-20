BINARY=terraform-ui-core

all: build testing cleanup

build:
	@echo "~RUNNING BUILD~"
	go build -o $(BINARY) cmd/terraform-ui/main.go

testing:
	@echo "~RUNNING TEST~"
	go test --cover ./...

cleanup:
	@echo "~RUNNING CLEANUP~"
	rm -f ./terraform-ui