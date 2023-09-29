# Homework 3


### Запуск приложения

    make up

### Остановка приложения
    
    make down

### Все манифесты в папке `manifests`


### Коллкция Postman и результаты тестирования в папке `newnan`


На MacOS не получается сделать запросы по hostname, поэтому 127.0.0.1

### Установка PostgreSQL (make helm_install_postgres)

    helm install postgresql bitnami/postgresql -n app -f pg_values.yaml

 
#### Зависимости

    brew install helm
    helm repo add bitnami https://charts.bitnami.com/bitnami



