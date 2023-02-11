SHELL := /bin/bash

.PHONY: docker_clean
docker_clean:
	docker container prune -f
	docker volume prune -f
	docker system prune -af


.PHONY: kube_start
kube_start:
	minikube start

.PHONY: kube_stop
kube_stop:
	minikube stop
	minikube delete

.PHONY: kube_apply
kube_apply:
	kubectl create -f kubernetes/mysql-secret.yaml
	kubectl apply -f kubernetes/mysql-db-pv.yaml
	kubectl apply -f kubernetes/mysql-db-pvc.yaml
	kubectl apply -f kubernetes/mysql-db-deployment.yaml
	kubectl apply -f kubernetes/mysql-db-service.yaml

