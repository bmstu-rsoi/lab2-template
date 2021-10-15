# Лабораторная работа #2

![GitHub Classroom Workflow](../../workflows/GitHub%20Classroom%20Workflow/badge.svg?branch=master)

## Microservices

### Формулировка

В рамках второй лабораторной работы _по вариантам_ требуется реализовать систему, состоящую из нескольких
взаимодействующих друг с другом сервисов.

### Требования

1. Каждый сервис имеет свое собственное хранилище, если оно ему нужно. Для учебный целей можно использовать одну базу
   данных, но каждая система работает _только_ со своей схемой. Запросы к другой схеме _запрещены_.
1. Для межсервисного взаимодействия использовать HTTP (придерживаться RESTful). Допускается использовать и другие
   протоколы, например grpc, но это требуется согласовать с преподавателем.
1. Выделить Gateway Service как единую точку входа и межсервисной коммуникации. Горизонтальные запросы между сервисами
   делать _нельзя_.
1. Код хранить на Github, для сборки использовать Github Actions.
1. Для автоматических прогонов тестов в файле [autograding.json](.github/classroom/autograding.json)
   и [classroom.yml](.github/workflows/classroom.yml) заменить `<variant>` на ваш вариант.
1. Каждый сервис должен быть завернут в docker.
1. В [classroom.yml](.github/workflows/classroom.yml) дописать шаги на сборку, прогон unit-тестов и деплой каждого
   сервиса на heroku.

### Пояснения

1. Для локальной разработки можно использовать Postgres в docker, для этого нужно запустить docker compose up -d,
   поднимется контейнер с Postgres 13, будет создана БД `services` и пользователь `program`:`test`. Для создания схем
   нужно прописать в [20-create-schemas.sh](postgres/20-create-schemas.sh) свой вариант задания в переменную `VARIANT`.
   После поднятия контейнера будут созданы схемы, описанные в файлах [schema-$VARIANT](postgres/schemes) по вариантам.
1. Горизонтальную коммуникацию между сервисами делать нельзя.

   ![Services](images/services.png)

   Предположим, у нас сервисы `UserService`, `OrderService`,
   `WarehouseService` и `Gateway`:
    * На `Gateway` от пользователя `Alex` приходит запрос `Купить товар с productName: 'Lego Technic 42129`.
    * `Gateway` -> `UserService` проверяем что пользователь существует и получаем `userUid` пользователя
      по `login: Alex`.
    * `Gateway` -> `WarehouseService` получаем `itemUid` товара по `productName` и резервируем его для заказа.
    * `Gateway` -> `OrderService` с `userUid` и `itemUid` и создаем заказ с `orderUid`.
    * `Gateway` -> `WarehouseService` с `orderUid` и переводим товар `itemUid` из статуса `Зарезервировано` в
      статус `Заказан` и прописываем ссылку на `orderUid`.

### Прием задания

1. При получении задания у вас создается fork этого репозитория для вашего пользователя.
1. После того, как все тесты успешно завершатся, в Github Classroom на Dashboard будет отмечено успешное выполнение
   тестов.

### Варианты заданий

Вариант заданий берутся исходя из номера
в [списке группы](https://docs.google.com/spreadsheets/d/1BT5iLgERiWUPPn4gtOQk4KfHjVOTQbUS7ragAJrl6-Q) mod 4 + 1. В
каждой группе нумерация вариантов начинается с 1.

1. [Flight Booking System](v1/README.md)
1. [Hotels Booking System](v2/README.md)
1. [Car Rental System](v3/README.md)
1. [Library System](v4/README.md)