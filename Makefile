PHONY: init
init:
	cp .github/hooks/* .git/hooks
	chmod +x .git/hooks/*
