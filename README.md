# Homework 6
## Идемпотетный метод создания заказа

### Описание

В данном примере при создании заказа `POST http://arch.homework/v1/order` 
используется заголовок `X-Idempotency-Key`.

Если заказ уже был создан с таким ключом, то возвращается код `200` и его id.
Если нет, то создается новый заказ и возвращается его id.

    { "order_id": 1123 }



Запуск приложения

    make up

Остановка приложения
    
    make down


## Тестирование

### Запуск тестов

    make newman

### Результаты тестов


![Tests](./newman/newman_3.jpg?raw=true "Tests")





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
