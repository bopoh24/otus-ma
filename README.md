# Homework 3


### Зависимости

    brew install helm
    helm repo add bitnami https://charts.bitnami.com/bitnami

### Установка PostgreSQL

    make helm_install_postgres


### Удаление PostgreSQL

    make helm_delete_postgres






Применить манифесты можно командой, создает неймспейс "app" и применяет манифесты:

    make apply

Удалить неймспейc "app" можно командой:

    make delete


Билд образа:

    make image_build


Пуш образа:

    make image_push

