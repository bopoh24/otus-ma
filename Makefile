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

build_image_api_gateway: ### build api gateway docker image
	@echo "Building image..."
	docker build --platform linux/amd64 -t bopoh24/api_gateway:latest -f api_gateway.dockerfile .
	@echo "Image built successfully!"

newman: ### run newman tests
	@echo "Running newman tests..."
	newman run newman/postman.json
.PHONY:newman


keycloak_up: ### start keycloak
	@echo "Starting keycloak..."
	@kubectl create namespace auth --dry-run=client -o yaml | kubectl apply -f -
	# apply keycloak manifests
	@kubectl apply -n auth -f keycloak/manifests
	@helm install auth-server -n auth  oci://registry-1.docker.io/bitnamicharts/keycloak -f keycloak/values.yaml
	@echo "Done!"

keycloak_down: ### stop keycloak
	@echo "Stopping keycloak..."
	helm delete -n auth auth-server
	@kubectl delete ns auth
	@echo "Done!"


krakend_up: ### start krakend
	@echo "Starting krakend..."
	@kubectl create namespace gateway --dry-run=client -o yaml | kubectl apply -f -
	@helm install krakend -n gateway equinixmetal/krakend -f api-gateway/values.yaml
	@echo "Done!"

krakend_down: ### stop krakend
	@echo "Stopping krakend..."
	helm delete -n gateway krakend
	@kubectl delete ns gateway
	@echo "Done!"
