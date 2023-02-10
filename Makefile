SHELL := /bin/bash

.PHONY: docker_clean
docker_clean:
	docker container prune -f
	docker volume prune -f
	docker system prune -af

