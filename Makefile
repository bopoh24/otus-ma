APP_IMAGE_NAME=bopoh24/simple_server:latest
MIGRATE_IMAGE_NAME=bopoh24/simple_server_migrate:latest

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ### this help information
	@awk 'BEGIN {FS = ":.*##"; printf "\nMakefile help:\n  make \033[36m<target>\033[0m\n"} /^[.a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

helm_up: ### install helm chart
	@echo "Creating namespace..."
	kubectl create namespace app --dry-run=client -o yaml | kubectl apply -f -
	@echo "Install helm chart..."
	helm install -n app simple-server ./chart
	@echo "Done!"
.PHONY:helm_up

helm_down: ### delete helm chart
	@echo "Uninstalling helm chart..."
	helm delete -n app simple-server
	@echo "Deleting namespace..."
	kubectl delete ns app
	@echo "Done!"


up: ### create namespace "app" and run app
	@echo "Creating namespace..."
	kubectl create namespace app --dry-run=client -o yaml | kubectl apply -f -
	@echo "Installing postgresql..."
	helm install postgresql bitnami/postgresql -n app --version 12.12.10 -f pg_values.yaml
	@echo "Applying manifests..."
	kubectl apply -f ./manifests
	@echo "Done!"
.PHONY:up

down: ### stop app and delete namespace "app"
	@echo "Deleting postgresql..."
	helm delete postgresql -n app
	@echo "Deleting k8s manifests..."
	kubectl delete ns app
	kubectl delete -n app pv postgres-pv
	@echo "Done!"
.PHONY:down

build_images: build_image_app build_image_migrate ### build docker images
.PHONY:build_images

push_images: ### push docker images to docker hub
	@echo "Pushing image..."
	docker push ${APP_IMAGE_NAME}
	docker push ${MIGRATE_IMAGE_NAME}
	@echo "Image pushed successfully!"
.PHONY:push_images

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

newman: ### run newman tests
	@echo "Running newman tests..."
	newman run newman/postman.json
.PHONY:newman
