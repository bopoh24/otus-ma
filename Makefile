APP_IMAGE_NAME=bopoh24/simple_server:latest
MIGRATE_IMAGE_NAME=bopoh24/simple_server_migrate:latest

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ### this help information
	@awk 'BEGIN {FS = ":.*##"; printf "\nMakefile help:\n  make \033[36m<target>\033[0m\n"} /^[.a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

up: app_up keycloak_up krakend_up ### start all services
.PHONY:up

down: app_down keycloak_down krakend_down ### stop all services

app_up: ### install app helm chart
	@echo "Creating namespace..."
	kubectl create namespace app --dry-run=client -o yaml | kubectl apply -f -
	@echo "Install helm chart for simple-server..."
	helm install -n app simple-server ./app/chart
	@echo "Done!"
.PHONY:helm_up

app_down: ### delete app helm chart
	@echo "Stopping simple-server..."
	helm delete -n app simple-server
	@echo "Deleting namespace..."
	kubectl delete ns app
	@echo "Done!"

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
	@helm install krakend -n gateway equinixmetal/krakend -f krakend/values.yaml
	@echo "Done!"

krakend_down: ### stop krakend
	@echo "Stopping krakend..."
	helm delete -n gateway krakend
	@kubectl delete ns gateway
	@echo "Done!"

newman: ### run newman tests
	@echo "Running newman tests..."
	newman run newman/postman_6.json
.PHONY:newman

push_images: ### push docker images to docker hub
	@echo "Pushing image..."
	docker push ${APP_IMAGE_NAME}
	docker push ${MIGRATE_IMAGE_NAME}
	@echo "Image pushed successfully!"
.PHONY:push_images

build_images: build_image_app build_image_migrate ### build docker images
.PHONY:build_images

build_image_app: ### build app docker image
	@echo "Building image..."
	docker build --platform linux/amd64 -t ${APP_IMAGE_NAME} -f ./app/app.dockerfile ./app
	@echo "Image built successfully!"
.PHONY:build_image_app

build_image_migrate: ### build migrations docker image
	@echo "Building image..."
	docker build --platform linux/amd64 -t ${MIGRATE_IMAGE_NAME} -f ./app/migrate.dockerfile ./app
	@echo "Image built successfully!"
.PHONY:build_image_migrate
