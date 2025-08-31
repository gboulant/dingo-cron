all: test

test.pkg:
	@go test

test.cmd:
	@make -C demos/demo.cron.01.crontab test

test: test.pkg
testall: test.pkg test.cmd

# -----------------------------------------------------------------------
# Clean the workspace
clean:
	@rm -f *~
	@make -C demos/demo.cron.01.crontab clean


