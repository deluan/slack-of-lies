.PHONY: run
run:
	@reflex -s -d none -r '\.go$\' -- go run main.go
