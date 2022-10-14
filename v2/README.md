## Hotels Booking System

Система предоставляет пользователю сервис поиска и бронирования отелей на интересующие даты. В зависимости от количества
заказов система лояльности дает скидку пользователю на новые бронирования.

### Структура Базы Данных

#### Reservation Service

Сервис запускается на порту 8070.

```sql
CREATE TABLE reservation
(
    id              SERIAL PRIMARY KEY,
    reservation_uid uuid UNIQUE NOT NULL,
    username        VARCHAR(80) NOT NULL,
    payment_uid     uuid        NOT NULL,
    hotel_id        INT REFERENCES hotels (id),
    status          VARCHAR(20) NOT NULL
        CHECK (status IN ('PAID', 'CANCELED')),
    start_date      TIMESTAMP WITH TIME ZONE,
    end_data        TIMESTAMP WITH TIME ZONE
);

CREATE TABLE hotels
(
    id        SERIAL PRIMARY KEY,
    hotel_uid uuid         NOT NULL UNIQUE,
    name      VARCHAR(255) NOT NULL,
    country   VARCHAR(80)  NOT NULL,
    city      VARCHAR(80)  NOT NULL,
    address   VARCHAR(255) NOT NULL,
    stars     INT,
    price     INT          NOT NULL
);
```

#### Payment Service

Сервис запускается на порту 8060.

```sql
CREATE TABLE payment
(
    id          SERIAL PRIMARY KEY,
    payment_uid uuid        NOT NULL,
    status      VARCHAR(20) NOT NULL
        CHECK (status IN ('PAID', 'CANCELED')),
    price       INT         NOT NULL
);
```

#### Loyalty Service

Сервис запускается на порту 8050.

```sql
CREATE TABLE loyalty
(
    id                SERIAL PRIMARY KEY,
    username          VARCHAR(80) NOT NULL UNIQUE,
    reservation_count INT         NOT NULL DEFAULT 0,
    status            VARCHAR(80) NOT NULL DEFAULT 'BRONZE'
        CHECK (status IN ('BRONZE', 'SILVER', 'GOLD')),
    discount          INT         NOT NULL
);
```

### Описание API

#### Получить список отелей

```http request
GET {{baseUrl}}/api/v1/hotels&page={{page}}&size={{size}}
```

#### Получить полную информацию о пользователе

Возвращается информация о бронированиях и статусе в системе лояльности.

```http request
GET {{baseUrl}}/api/v1/me
X-User-Name: {{username}}
```

#### Информация по всем бронированиям пользователя

```http request
GET {{baseUrl}}/api/v1/reservations
X-User-Name: {{username}}
```

#### Информация по конкретному бронированию

При запросе требуется проверить, что бронирование принадлежит пользователю.

```http request
GET {{baseUrl}}/api/v1/reservations/{{reservationUid}}
X-User-Name: {{username}}
```

#### Забронировать отель

Пользователь вызывает метод `GET {{baseUrl}}/api/v1/hotels` и выбирает нужный отель и в запросе на бронирование
передает:

* `hotelUid` (UUID отеля) – берется из запроса `/hotels`;
* `startDate` и `endDate` (дата начала и конца бронирования) – задается пользователем.

Система проверяет, что отель с таким `hotelUid` существует. Считаем что в отеле бесконечное количество мест.

Считается количество ночей (`endDate` – `startDate`), вычисляется общая сумма бронирования, выполняется обращение в
Loyalty Service и получается скидка в зависимости от статуса клиента:

* BRONZE – 5%
* SILVER – 7%
* GOLD – 10%

После применения скидки выполняется запрос в Payment Service и создается новая запись об оплате. После этого выполняется
обращение в сервис Loyalty Service, увеличивается счетчик бронирований. По-умолчанию у клиента статус `BRONZE`,
статус `SILVER` присваивается после 10 бронирований, `GOLD` после 20.

```http request
POST {{baseUrl}}/api/v1/reservations
Content-Type: application/json
X-User-Name: {{username}}

{
  "hotelUid": "049161bb-badd-4fa8-9d90-87c9a82b0668",
  "startDate": "2021-10-08",
  "endDate": "2021-10-11"
}
```

#### Отменить бронирование

* Статус бронирования помечается как `CANCELED`.
* В Payment Service запись об оплате помечается отмененной (статус `CANCELED`).
* Loyalty Service уменьшается счетчик бронирований. Так же возможно понижение статуса лояльности, если счетчик стал ниже
  границы уровня.

```http request
DELETE {{baseUrl}}/api/v1/reservations/{{reservationUid}}
X-User-Name: {{username}}
```

#### Получить информацию о статусе в программе лояльности

```http request
GET {{baseUrl}}/api/v1/loyalty
X-User-Name: {{username}}
```

Описание в формате [OpenAPI](%5Binst%5D%5Bv2%5D%20Hotels%20Booking%20System.yml).

### Данные для тестов

Создать данные для тестов:

```yaml
hotels:
  – id: 1
    hotelUid: "049161bb-badd-4fa8-9d90-87c9a82b0668"
    name: "Ararat Park Hyatt Moscow"
    country: "Россия"
    city: "Москва"
    address: "Неглинная ул., 4"
    stars: 5,
    price: 10000

loyalty:
  - id: 1
    username: "Test Max"
    reservation_count: 25
    status: "GOLD"
    discount: 10
```
