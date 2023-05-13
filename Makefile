
install:
	go install ./cmd/grace

.PHONY: install

.DEFAULT_GOAL := install
