# Homework 3
## Простой CRUD сервис с использованием PostgreSQL 

## Helm

`./chart` - шаблонизация приложения в helm чартах

Запуск приложения

    make helm_up

Остановка приложения
    
    make helm_down



## Kubernetes манифесты

`./manifests` - все манифесты


Запуск приложения

    make up

Остановка приложения
    
    make down


## Тестирование

`./newnan` - коллекция Postman и результаты тестирования


#### Установка PostgreSQL (см. Makefile)

    helm install postgresql bitnami/postgresql -n app --version 12.12.10 -f pg_values.yaml


## Остальные команды

    make help

## Установка зависимостей

    brew install helm
    helm repo add bitnami https://charts.bitnami.com/bitnami
