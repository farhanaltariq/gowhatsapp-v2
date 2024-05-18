.PHONY: run test

run:
	swag init && go run .

test:
	cd test && ginkgo