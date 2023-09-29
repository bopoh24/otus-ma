APP_IMAGE_NAME=bopoh24/simple_server:latest
MIGRATE_IMAGE_NAME=bopoh24/simple_server_migrate:latest

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ### this help information
	@awk 'BEGIN {FS = ":.*##"; printf "\nMakefile help:\n  make \033[36m<target>\033[0m\n"} /^[.a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

up: manifests_init helm_install_postgres manifests_apply ### create namespace "app" and run app
.PHONY:up

down: helm_delete_postgres manifests_delete ### stop app and delete namespace "app"
	@echo "Done!"
.PHONY:down

build: build_image_app build_image_migrate ### build docker images
.PHONY:image_build

build_image_app: ### build app docker image
	@echo "Building image..."
	docker build --platform linux/amd64 -t ${APP_IMAGE_NAME} -f app.dockerfile .
	@echo "Image built successfully!"
.PHONY:build_image_app

build_image_migrate: ### build migrations docker image
	@echo "Building image..."
	docker build --platform linux/amd64 -t ${MIGRATE_IMAGE_NAME} -f migrate.dockerfile .
	@echo "Image built successfully!"
.PHONY:build_image_migrate

push_images: ### push docker images to docker hub
	@echo "Pushing image..."
	docker push ${APP_IMAGE_NAME}
	docker push ${MIGRATE_IMAGE_NAME}
	@echo "Image pushed successfully!"
.PHONY:push_images

manifests_apply: ### create namespace "app" and apply k8s manifests
	@echo "Applying k8s manifests..."
	@kubectl create namespace app --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f ./manifests
	@echo "Done!"
.PHONY:manifests_apply

manifests_delete: ### delete namespace "app"
	@echo "Deleting k8s manifests..."
	kubectl delete ns app
.PHONY:manifests_delete

manifests_init:
	@echo "Creating namespace..."
	kubectl create namespace app --dry-run=client -o yaml | kubectl apply -f -
	@echo "Applying secrets and configmap..."
	kubectl apply -f manifests/conf
	@echo "Applying PV and PVC..."
	kubectl apply -f manifests/pvc
.PHONY:manifests_init

helm_install_postgres: ### install postgresql
	@echo "Installing postgresql..."
	helm install postgresql bitnami/postgresql -n app -f pg_values.yaml
.PHONY:helm_install_postgres

helm_delete_postgres: ### delete postgresql
	@echo "Deleting postgresql..."
	helm delete postgresql -n app
.PHONY:helm_delete_postgres

newman: ### run newman tests
	@echo "Running newman tests..."
	newman run newman/postman.json
.PHONY:newman
