APP_IMAGE_NAME=bopoh24/simple_server:latest
MIGRATE_IMAGE_NAME=bopoh24/simple_server_migrate:latest

NAMESPACE=app
RELEASE_NAME=booking-srv

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ### this help information
	@awk 'BEGIN {FS = ":.*##"; printf "\nMakefile help:\n  make \033[36m<target>\033[0m\n"} /^[.a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help



newman: ### run newman tests
	@echo "Running newman tests..."
	newman run newman/postman.json
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


# App ==================================================================================================================
up:
	@echo "Creating namespace..."
	kubectl create namespace ${NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
	@echo "Install helm chart for ${RELEASE_NAME}..."
	@helm install -n ${NAMESPACE} ${RELEASE_NAME} \
		--set-file krakend.krakend.config=deployments/chart/config/krakend/krakend.json \
		./deployments/chart
	@echo "Done!"

upgrade:
	@echo "Upgrading helm chart for ${RELEASE_NAME}..."
	@helm upgrade -n ${NAMESPACE} ${RELEASE_NAME} \
		--set-file krakend.krakend.config=deployments/chart/config/krakend/krakend.json \
		./deployments/chart
	@echo "Done!"

down:
	@echo "Stopping ${RELEASE_NAME}..."
	helm delete -n ${NAMESPACE} ${RELEASE_NAME}
	@echo "Deleting namespace..."
	kubectl delete ns ${NAMESPACE}
	@echo "Done!"


# Ingress controller ===================================================================================================
up_ctrl: ### install ingress-nginx helm chart
	@echo "Creating namespace for nginx-ingress..."
	kubectl create namespace ctrl
	@echo "Install helm chart for nginx-ingress..."
	@-helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
	helm repo update
	helm install nginx ingress-nginx/ingress-nginx --namespace ctrl -f deployments/nginx-ingress.yaml
	@echo "Done!"

down_ctrl:
	@echo "Stopping nginx-ingress..."
	@-helm delete -n ctrl nginx
	@echo "Deleting namespace..."
	kubectl delete ns ctrl
	@echo "Done!"
