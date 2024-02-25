# Homework 9 - Stream processing

*Событийное взаимодействие с использование брокера сообщений для нотификаций (уведомлений)*

## Сервис бронирования услуг

### Схема работы

Pегистрируется новый бизнес-пользователь, создает конмпанию, размещает услугу.

Регистрируется клиент, ищет услугу, бронирует.

При бронировании создается транзакция: 
- создается бронь услуги, она перезодит в статус "reserved" (микросервис booking)
- создается платеж (микросервис payment)
- eсли платеж прошел успешно списываются деньги с баланса (микросервис payment), бронь переходит в статус "paid" (микросервис booking), eсли не удалось поменять статус, то платеж отменяется (микросервис payment)
- если платеж не прошел, то услуга переходит обратно в статус "open" (микросервис booking) и фиксируется неуспешная транзакция (микросервис payment)

После успешной или неуспешной транзакции отправляется уведомление (микросервис notification)

Сообщение отправляется на почту менеджеру(ам) компании и клиенту.


![Mailhog](./mailhog.jpg?raw=true "Tests")

*Имя хоста для Mailhog `mailhog.booksvc.com`*


#### Схема успешного формирования заказа услуги


![Sequence](./booking_ok.jpeg?raw=true "Tests")




### Описание приложения

Содержит микросервисы: 
- ``company`` 
- ``customer`` 
- ``booking`` 
- ``payment`` 
- ``notification``


`Keycloak` - сервис авторизации

`KrakenD` - API Gateway

`PostgreSQL` - база данных, на каждый сервис своя база.

`Kafka` - брокер сообщений

`MailHog` - почтовый сервер-заглушка для тестирования



Запуск приложения

    make up

Остановка приложения
    
    make down




## Тестирование

### Запуск тестов

    make newman

### Результаты тестов


![Tests](./newman/screen1.jpg?raw=true "Tests")
![Tests](./newman/screen2.jpg?raw=true "Tests")
![Tests](./newman/screen3.jpg?raw=true "Tests")





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

### MailHog

    helm repo add codecentric https://codecentric.github.io/helm-charts
    helm repo update

### Prometheus и Grafana

    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    helm repo update
