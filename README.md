# Homework 2


## Должен работать rewrite:

    curl arch.homework/otusapp/aeugene/health -> рерайт пути на arch.homework/health 



Применить манифесты можно командой, создает неймспейс "app" и применяет манифесты:

    make apply

Удалить неймспейc "app" можно командой:

    make delete


Билд образа:

    make image_build


Пуш образа:

    make image_push

