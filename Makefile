VERSION := $(shell cat VERSION)

install:
	go install .

release:
	@echo "Releasing version $(VERSION) ..."
	jfrog rt gp go $(VERSION)

.PHONY: release
