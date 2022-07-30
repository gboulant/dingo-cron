all: test

test.pkg:
	@go test ${GOTEST_OPTIONS}

test.cmd: go.mod
	@make -C demos/demo.cron.01.crontab test

test: test.pkg
testall: test.pkg test.cmd

# -----------------------------------------------------------------------
# Module management
GONAME="galuma.net/systemd/cron"
go.mod:
	go mod init ${GONAME} && \
	go mod tidy

# -----------------------------------------------------------------------
# Clean the workspace
clean:
	@rm -f *~
	@make -C demos/demo.cron.01.crontab clean


