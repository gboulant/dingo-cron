all: test

test:
	@go test ${GOTEST_OPTIONS}

# -----------------------------------------------------------------------
# Module management
GONAME="galuma.net/systemd/cron"
go.mod:
	go mod init ${GONAME} && \
	go mod tidy

# -----------------------------------------------------------------------
# Clean the workspace
clean:
	@rm -f go.mod go.sum *~


