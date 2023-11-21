# Homework 5
## Простой сервис с auth proxy и api-gateway


Запуск приложения

    make up

Остановка приложения
    
    make down


### Схема приложения

![Schema](./screenshots/app_scheme.jpg?raw=true "App Scheme")


## Тестирование

### Запуск тестов

    make newman

### Результаты тестов


![Tests](./newman/newman_1.jpg?raw=true "Tests")
![Tests](./newman/newman_2.jpg?raw=true "Tests")





### Установка ingress-nginx контроллера с метриками

    kubectl create namespace m && \
    helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx/ && \
    helm repo update && helm install nginx ingress-nginx/ingress-nginx --namespace m -f nginx-ingress.yaml



## Остальные команды

    make help

## Установка зависимостей

    brew install helm
    helm repo add bitnami https://charts.bitnami.com/bitnami

### KrakenD

    helm repo add equinixmetal https://helm.equinixmetal.com
    helm repo update

### Prometheus и Grafana

    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    helm repo update
