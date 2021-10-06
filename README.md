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
1. Каждый сервис должен быть завернут в docker.
1. В [classroom.yml](.github/workflows/classroom.yml) дописать шаги на сборку, прогон unit-тестов и деплой каждого
   сервиса на heroku.

### Пояснения

1. Для локальной разработки можно использовать Postgres в docker, для этого нужно запустить docker compose up -d,
   поднимется контейнер с Postgres 13, будет создана БД `services` и пользователь `program`:`test`.

   Для создания схем нужно прописать в [20-create-schemas.sh](postgres/20-create-schemas.sh):
   ```
   CREATE SCHEMA IF NOT EXISTS <schema-name>;
   ```
1. Горизонтальную коммуникацию между сервисами делать нельзя. Предположим, у нас сервисы `UserService`, `OrderService`,
   `WarehouseService` и `Gateway`:
    * На `Gateway` от пользователя `Alex` приходит запрос `Купить товар с productName: 'Lego Technic 42129`.
    * Для оформления заказа требуется проверить что пользователь c `userId` существует, получить его имя, сходить
      в `WarehouseService` получить `itemUid` по .
    * `Gateway` -> `UserService` получаем `userUid` пользователя по `login: Alex`.
    * `Gateway` -> `WarehouseService` получаем `itemUid` товара по `productName` и резервируем его для заказа.
    * `Gateway` -> `OrderService` с `userUid` и `itemUid` и создаем заказ с `orderUid`.
    * `Gateway` -> `WarehouseService` с `orderUid` и переводим товар `itemUid` из статуса `Зарезервировано` в
      статус `Заказан` и прописываем ссылку на `orderUid`.

### Прием задания

1. При получении задания у вас создается fork этого репозитория для вашего пользователя.
1. После того, как все тесты успешно завершатся, в Github Classroom на Dashboard будет отмечен успешное выполнение
   тестов.

### Варианты заданий

#### Flight Booking System

Система предоставляет пользователю возможность поиска и покупки билетов. В зависимости от количества выполненных
перелетов, пользователю предоставляется скидка на перелет и начисляются баллы, которые он может использовать для оплаты.

##### Ticket Service

```sql
CREATE TABLE ticket
(
    id         SERIAL PRIMARY KEY,
    ticket_uid uuid UNIQUE NOT NULL,
    username   VARCHAR(80) NOT NULL,
    flight_uid uuid        NOT NULL,
    price      INT         NOT NULL
);
```

##### Flight Service

```sql
CREATE TABLE flight
(
    flight_number   VARCHAR(80)              NOT NULL,
    datetime        TIMESTAMP WITH TIME ZONE NOT NULL,
    from_airport_id INT REFERENCES airport (id),
    to_airport_id   INT REFERENCES airport (id)
);

CREATE TABLE airport
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(255),
    country VARCHAR(255)
);
```

##### Bonus Service

```sql
CREATE TABLE privilege
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(80) NOT NULL UNIQUE,
    balance  INT
);

CREATE TABLE privilege_history
(
    id             SERIAL      NOT NULL,
    privilege_id   INT REFERENCES privilege (id),
    datetime       TIMESTAMP   NOT NULL,
    balance_diff   INT         NOT NULL,
    operation_type VARCHAR(20) NOT NULL
        CHECK (operation_type IN ('FILL_IN_BALANCE', 'DEBIT_THE_ACCOUNT', 'FILLED_BY_MONEY'))
);
```

[OpenAPI definition](%5Binst%5D%5Bv1%5D%20Flight%20Booking%20System.yml)

#### Hotels Booking System

Система предоставляет пользователю сервис поиска и бронирования отелей на интересующие даты. В зависимости от количества
заказов система лояльности дает скидку пользователю на новые бронирования.

##### Reservation Service

```sql
CREATE TABLE reservation
(
    id              SERIAL PRIMARY KEY,
    reservation_uid uuid UNIQUE NOT NULL,
    username        VARCHAR(80) NOT NULL,
    payment_uid     uuid        NOT NULL,
    hotel_id        INT REFERENCES hotels (id),
    start_date      TIMESTAMP WITH TIME ZONE,
    end_data        TIMESTAMP WITH TIME ZONE
);

CREATE TABLE hotels
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(255) NOT NULL,
    country VARCHAR(80)  NOT NULL,
    address VARCHAR(255) NOT NULL,
    stars   INT
);
```

##### Payment Service

```sql
CREATE TABLE payment
(
    id          SERIAL PRIMARY KEY,
    payment_uid uuid        NOT NULL,
    status      VARCHAR(20) NOT NULL
        CHECK (status IN ('RESERVED', 'PAID', 'REVERSED', 'CANCELED')),
    price       INT         NOT NULL
);
```

##### Loyalty Service

```sql
CREATE TABLE loyalty
(
    id       SERIAL PRIMARY KEY,
    status   VARCHAR(80) NOT NULL DEFAULT 'BRONZE'
        CHECK (status IN ('BRONZE', 'SILVER', 'GOLD')),
    discount INT         NOT NULL
);
```

[OpenAPI definition](%5Binst%5D%5Bv1%5D%20Hotels%20Booking%20System.yml)

#### Car Rental System

##### Cars Service

Система предоставляет пользователю возможность забронировать автомобиль на выбранные даты.

```sql
CREATE TABLE cars
(
    id                  SERIAL PRIMARY KEY,
    car_uid             uuid UNIQUE NOT NULL,
    brand               VARCHAR(80) NOT NULL,
    model               VARCHAR(80) NOT NULL,
    registration_number VARCHAR(20) NOT NULL,
    power               INT,
    type                VARCHAR(20)
        CHECK (type IN ('SEDAN', 'SUV', 'MINIVAN', 'ROADSTER')),
    availability        BOOLEAN     NOT NULL
);
```

##### RentalService

```sql
CREATE TABLE rental
(
    id          SERIAL PRIMARY KEY,
    username    VARCHAR(80)              NOT NULL,
    payment_uid uuid                     NOT NULL,
    car_uid     uuid                     NOT NULL,
    date_from   TIMESTAMP WITH TIME ZONE NOT NULL,
    date_to     TIMESTAMP WITH TIME ZONE NOT NULL,
    status      VARCHAR(20)              NOT NULL
        CHECK (status IN ('NEW', 'IN_PROGRESS', 'FINISHED', 'CANCELED', 'EXPIRED'))
);
```

##### Payment Service

```sql
CREATE TABLE payment
(
    id          SERIAL PRIMARY KEY,
    payment_uid uuid        NOT NULL,
    status      VARCHAR(20) NOT NULL
        CHECK (status IN ('RESERVED', 'PAID', 'REVERSED', 'CANCELED')),
    price       INT         NOT NULL
);
```

[OpenAPI definition](%5Binst%5D%5Bv1%5D%20Car%20Rental%20System.yml)

#### Library System

Система позволяет пользователю найти интересующую книгу и взять ее в библиотеке. Если у пользователя на руках есть уже N
книг, то он не может взять новую, пока не сдал старые. Если пользователь возвращает книги в хорошем состоянии и сдает их
в срок, то максимальное количество книг у него на руках увеличивается.

##### Reservation System

```sql
CREATE TABLE reservation
(
    id        SERIAL PRIMARY KEY,
    username  VARCHAR(80) NOT NULL,
    book_uid  uuid        NOT NULL,
    till_date TIMESTAMP   NOT NULL
)
```

##### Library System

```sql
CREATE TABLE library
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(80)  NOT NULL,
    address VARCHAR(255) NOT NULL
);

CREATE TABLE books
(
    id     SERIAL PRIMARY KEY,
    name   VARCHAR(255) NOT NULL,
    author VARCHAR(255),
    genre  VARCHAR(255)
);

CREATE TABLE library_books
(
    id         SERIAL PRIMARY KEY,
    book_uid   uuid NOT NULL,
    book_id    INT  NOT NULL
        REFERENCES books (id),
    library_id INT  NOT NULL
        REFERENCES library (id),
    condition  VARCHAR(20) DEFAULT 'EXCELLENT'
        CHECK (condition IN ('EXCELLENT', 'GOOD', 'BAD'))
);
```

##### Rating System

```sql
CREATE TABLE rating
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(80) NOT NULL,
    stars    INT         NOT NULL
        CHECK (stars BETWEEN 1 AND 100)
);
```

[Open API definition](%5Binst%5D%5Bv1%5D%20Library%20System.yml)