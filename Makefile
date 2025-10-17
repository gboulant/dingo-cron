all: test

test.pkg:
	@go test

test.cmd:
	@make -C demos/demo.cron.01.crontab test

test: test.pkg
testall: test.pkg test.cmd

# -----------------------------------------------------------------------
# doc and coverage
doc:
	@go tool doc -short
#	@go tool doc -C demos/demo.cron.01.crontab -cmd -all

cov:
	@go test -coverprofile=output.cov
	@go tool cover -func=output.cov

# -----------------------------------------------------------------------
# Clean the workspace
clean:
	@rm -f *~ output.*
	@make -C demos/demo.cron.01.crontab clean
