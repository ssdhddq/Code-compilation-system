. PHONY: launch_services launch_with_tests stop_services build_services

launch_services:
	docker-compose up --build

launch_with_tests:
	docker-compose --profile test build
	docker-compose --profile test up --abort-on-container-exit --exit-code-from app_test

stop_services:
	docker-compose --profile test down
