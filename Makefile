.PHONY: build
build:
	docker build -t netra .

.PHONY: clean
clean:
	docker-compose down
	docker rmi netra

.PHONY: run
run:
	docker-compose up
