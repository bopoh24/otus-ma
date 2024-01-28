NAMESPACE=app
RELEASE_NAME=booking-srv

CUSTOMER_IMAGE_NAME=bopoh24/b-srv-customer:latest
CUSTOMER_MIGRATE_IMAGE_NAME=bopoh24/b-srv-customer-migrate:latest


# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ### this help information
	@awk 'BEGIN {FS = ":.*##"; printf "\nMakefile help:\n  make \033[36m<target>\033[0m\n"} /^[.a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

# Customer service =====================================================================================================
build_customer_service: ### build customer service docker image
	@echo "Building image..."
	docker build --platform linux/amd64 -t ${CUSTOMER_IMAGE_NAME} -f ./customer/customer.dockerfile .
	@echo "Image built successfully!"
.PHONY:build_image_app

build_customer_service_migrate: ### build customer service migrate docker image
	@echo "Building image..."
	docker build --platform linux/amd64 -t ${CUSTOMER_MIGRATE_IMAGE_NAME} -f ./customer/customer.migrate.dockerfile .
	@echo "Image built successfully!"
.PHONY:build_image_migrate

push_customer_images: ### push customer service docker image to docker hub
	@echo "Pushing image..."
	docker push ${CUSTOMER_IMAGE_NAME}
	docker push ${CUSTOMER_MIGRATE_IMAGE_NAME}
	@echo "Image pushed successfully!"

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

# Newman
newman: ### run newman tests
	@echo "Running newman tests..."
	newman run newman/postman.json
.PHONY:newman


fwd_db: ### port forward to db
	kubectl port-forward pod/booking-srv-postgresql-0 5432:5432 -n app
