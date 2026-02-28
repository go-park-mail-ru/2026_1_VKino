PHONY: init
init:
	cp .github/hooks/* .git/hooks
	chmod +x .git/hooks/*
	migrate create -ext sql -dir ./migrations migration

MIGRATIONS_DIR := ./migrations