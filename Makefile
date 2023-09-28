include .env
IMAGE_NAME=bopoh24/simple_server:latest

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ### this help information
	@awk 'BEGIN {FS = ":.*##"; printf "\nMakefile help:\n  make \033[36m<target>\033[0m\n"} /^[.a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

image_build: ### build docker image
	@echo "Building image..."
	docker build --platform linux/amd64 -t ${IMAGE_NAME} .
	@echo "Image built successfully!"
.PHONY:image_build

image_push:
	@echo "Pushing image..."
	docker push ${IMAGE_NAME}
	@echo "Image pushed successfully!"
.PHONY:image_push

apply: ### create namespace "app" and apply k8s manifests
	@echo "Applying k8s manifests..."
	@kubectl create namespace app --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f ./manifests
	@echo "Done!"
.PHONY:apply

delete: ### delete namespace "app"
	@echo "Deleting k8s manifests..."
	kubectl delete ns app
.PHONY:delete


newman: ### run newman tests
	@echo "Running newman tests..."
	newman run postman.json
.PHONY:newman


helm_install_postgres: ### help install postgresql
	@kubectl create namespace app --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f manifests/pvc
	@echo "Installing postgresql..."
	helm install postgresql bitnami/postgresql -n app \
	--set primary.persistence.existingClaim=postgres-pvc \
	--set volumePermissions.enabled=true \
	--set global.postgresql.auth.postgresPassword=${POSTGRES_PASSWORD} \
	--set global.postgresql.auth.username=${POSTGRES_USER} \
	--set global.postgresql.auth.password=${POSTGRES_PASSWORD} \
	--set global.postgresql.auth.database=${POSTGRES_DB}
.PHONY:helm_install_postgres

helm_delete_postgres: ### help delete postgresql
	@echo "Deleting postgresql..."
	helm delete postgresql -n app
	kubectl delete pvc -n app postgres-pvc
	kubectl delete pv -n app postgres-pv


migrate_up: ### apply migrations to database
	migrate -path migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost/${POSTGRES_DB}?sslmode=disable up

migrate_down: ### rollback migrations from database
	migrate -path migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost/${POSTGRES_DB}?sslmode=disable down
