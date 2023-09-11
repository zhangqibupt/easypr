VERSION := $(shell cat VERSION)

release:
	@echo "Releasing version $(VERSION) ..."
	jfrog rt gp go $(VERSION)

.PHONY: release
