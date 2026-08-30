. PHONY: launch_services launch_with_tests stop_services build_services

launch_services:
	docker compose up --build
	docker build -t code-processor:latest -f codeProcessor/docker/Dockerfile codeProcessor/docker

launch_with_tests:
	docker compose --profile test build
	docker build -t code-processor:latest -f codeProcessor/docker/Dockerfile codeProcessor/docker
	docker compose --profile test up --abort-on-container-exit --exit-code-from app_test

stop_services:
	docker compose --profile test down
