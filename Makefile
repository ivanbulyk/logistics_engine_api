.PHONY: start

start:
	go run cmd/logistics/*.go


.DEFAULT_GOAL := start