# BOOKSVC


![Sequence](./scheme.jpeg?raw=true "App scheme")


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

*Если платеж не прошел, то услуга переходит обратно в статус "open" (микросервис booking) и отправляется уведомление о неуспешной транзакции (микросервис notification)*



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

Установка ingress-nginx контроллера

    make up_ctrl



## Тестирование

### Запуск тестов

    make newman

### Результаты тестов


![Tests](./newman/screen1.jpg?raw=true "Tests")
![Tests](./newman/screen2.jpg?raw=true "Tests")
![Tests](./newman/screen3.jpg?raw=true "Tests")


# Мониторинг


    make grafana


### Services

![Services](./dashboard1.jpg?raw=true "Services")

### PostgreSQL


![PostgreSQL](./dashboard2.jpg?raw=true "PostgreSQL")

### Kafka

![Kafka](./dashboard3.jpg?raw=true "Kafka")

