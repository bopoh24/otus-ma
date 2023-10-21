# Homework 4
## Простой CRUD сервис с использованием PostgreSQL и мониторингом

## Helm

`./chart` - шаблонизация приложения в helm чартах

Запуск приложения

    make helm_up

Остановка приложения
    
    make helm_down


### Конфигурация Grafana

`dashboard.json`  - дашборд для Grafana

`alerts.json` - правила алертинга для Grafana

### Скриншоты

![Dashboard](./screenshots/dashboard.jpg?raw=true "Dashboard")



![Alerts](./screenshots/alerts.jpg?raw=true "Alerts")
    


### Установка ingress-nginx контроллера с метриками

    kubectl create namespace m && \
    helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx/ && \
    helm repo update && helm install nginx ingress-nginx/ingress-nginx --namespace m -f nginx-ingress.yaml



## Тестирование

`./newman` - коллекция Postman и результаты тестирования


## Остальные команды

    make help

## Установка зависимостей

    brew install helm
    helm repo add bitnami https://charts.bitnami.com/bitnami
    
### Prometheus и Grafana

    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    helm repo update
