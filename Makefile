NAMESPACE=booksvc
RELEASE_NAME=booksvc

CUSTOMER_IMAGE_NAME=bopoh24/booksvc-customer:latest
CUSTOMER_MIGRATE_IMAGE_NAME=bopoh24/booksvc-customer-migrate:latest
COMPANY_IMAGE_NAME=bopoh24/booksvc-company:latest
COMPANY_MIGRATE_IMAGE_NAME=bopoh24/booksvc-company-migrate:latest
BOOKING_IMAGE_NAME=bopoh24/booksvc-booking:latest
BOOKING_MIGRATE_IMAGE_NAME=bopoh24/booksvc-booking-migrate:latest
PAYMENT_IMAGE_NAME=bopoh24/booksvc-payment:latest
PAYMENT_MIGRATE_IMAGE_NAME=bopoh24/booksvc-payment-migrate:latest

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ### this help information
	@awk 'BEGIN {FS = ":.*##"; printf "\nMakefile help:\n  make \033[36m<target>\033[0m\n"} /^[.a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

# Customer service =====================================================================================================

customer_service: customer_test ### build customer service docker image
	@echo "Building image..."
	docker build --platform linux/amd64 -t ${CUSTOMER_IMAGE_NAME} -f ./customer/customer.dockerfile .
	@echo "Image built successfully!"
	docker push ${CUSTOMER_IMAGE_NAME}
	@echo "Image pushed successfully!"
.PHONY:customer_service

customer_migrations: ### build customer service migrations docker image
	@echo "Building image..."
	docker build --platform linux/amd64 -t ${CUSTOMER_MIGRATE_IMAGE_NAME} -f ./customer/customer.migrate.dockerfile .
	@echo "Image built successfully!"
	docker push ${CUSTOMER_MIGRATE_IMAGE_NAME}
	@echo "Image pushed successfully!"
.PHONY:customer_migrations

customer_test: ### run tests for customer service
	@echo "Testing customer service..."
	go test -v ./customer/...
.PHONY:customer_test


# Company service =====================================================================================================

company_service: company_test ### build customer service docker image
	@echo "Building image..."
	docker build --platform linux/amd64 -t ${COMPANY_IMAGE_NAME} -f ./company/company.dockerfile .
	@echo "Image built successfully!"
	docker push ${COMPANY_IMAGE_NAME}
	@echo "Image pushed successfully!"
.PHONY:company_service

company_migrations: ### build customer service migrations docker image
	@echo "Building image..."
	docker build --platform linux/amd64 -t ${COMPANY_MIGRATE_IMAGE_NAME} -f ./company/company.migrate.dockerfile .
	@echo "Image built successfully!"
	docker push ${COMPANY_MIGRATE_IMAGE_NAME}
	@echo "Image pushed successfully!"
.PHONY:company_migrations

company_test: ### run tests for company service
	@echo "Testing company service..."
	go test -v ./company/...
.PHONY:company_test


# Booking service =====================================================================================================

booking_service: booking_test ### build booking service docker image
	@echo "Building image..."
	docker build --platform linux/amd64 -t ${BOOKING_IMAGE_NAME} -f ./booking/booking.dockerfile .
	@echo "Image built successfully!"
	docker push ${BOOKING_IMAGE_NAME}
	@echo "Image pushed successfully!"
.PHONY:booking_service

booking_migrations: ### build booking service migrations docker image
	@echo "Building image..."
	docker build --platform linux/amd64 -t ${BOOKING_MIGRATE_IMAGE_NAME} -f ./booking/booking.migrate.dockerfile .
	@echo "Image built successfully!"
	docker push ${BOOKING_MIGRATE_IMAGE_NAME}
	@echo "Image pushed successfully!"
.PHONY:booking_migrations

booking_test: ### run tests for booking service
	@echo "Testing booking service..."
	go test -v ./booking/...
.PHONY:booking_test


# Payment service =====================================================================================================

payment_service: payment_test ### build payment service docker image
	@echo "Building image..."
	docker build --platform linux/amd64 -t ${PAYMENT_IMAGE_NAME} -f ./payment/payment.dockerfile .
	@echo "Image built successfully!"
	docker push ${PAYMENT_IMAGE_NAME}
	@echo "Image pushed successfully!"
.PHONY:payment_service

payment_migrations: ### build payment service migrations docker image
	@echo "Building image..."
	docker build --platform linux/amd64 -t ${PAYMENT_MIGRATE_IMAGE_NAME} -f ./payment/payment.migrate.dockerfile .
	@echo "Image built successfully!"
	docker push ${PAYMENT_MIGRATE_IMAGE_NAME}
	@echo "Image pushed successfully!"

payment_test: ### run tests for payment service
	@echo "Testing payment service..."
	go test -v ./payment/...


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
	kubectl port-forward pod/${RELEASE_NAME}-postgresql-0 5432:5432 -n ${NAMESPACE}
