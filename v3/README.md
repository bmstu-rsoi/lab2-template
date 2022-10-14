## Car Rental System

Система предоставляет пользователю возможность забронировать автомобиль на выбранные даты.

### Структура Базы Данных

#### Cars Service

Сервис запускается на порту 8070.

```sql
CREATE TABLE cars
(
    id                  SERIAL PRIMARY KEY,
    car_uid             uuid UNIQUE NOT NULL,
    brand               VARCHAR(80) NOT NULL,
    model               VARCHAR(80) NOT NULL,
    registration_number VARCHAR(20) NOT NULL,
    power               INT,
    price               INT         NOT NULL,
    type                VARCHAR(20)
        CHECK (type IN ('SEDAN', 'SUV', 'MINIVAN', 'ROADSTER')),
    availability        BOOLEAN     NOT NULL
);
```

#### Rental Service

Сервис запускается на порту 8060.

```sql
CREATE TABLE rental
(
    id          SERIAL PRIMARY KEY,
    rental_uid  uuid UNIQUE              NOT NULL,
    username    VARCHAR(80)              NOT NULL,
    payment_uid uuid                     NOT NULL,
    car_uid     uuid                     NOT NULL,
    date_from   TIMESTAMP WITH TIME ZONE NOT NULL,
    date_to     TIMESTAMP WITH TIME ZONE NOT NULL,
    status      VARCHAR(20)              NOT NULL
        CHECK (status IN ('IN_PROGRESS', 'FINISHED', 'CANCELED'))
);
```

#### Payment Service

Сервис запускается на порту 8050.

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

### Описание API

#### Получить список всех доступных для бронирования автомобилей

Если передан флаг `showAll = true`, то выводить автомобили в резерве (`availability = false`).

```http request
GET {{baseUrl}}/api/v1/cars&page={{page}}&size={{size}}
```

#### Получить информацию о всех арендах пользователя

```http request
GET {{baseUrl}}/api/v1/rental
X-User-Name: {{username}}
```

#### Информация по конкретной аренде пользователя

При запросе требуется проверить, что аренда принадлежит пользователю.

```http request
GET {{baseUrl}}/api/v1/rental/{{rentalUid}}
X-User-Name: {{username}}
```

#### Забронировать автомобиль

Пользователь вызывает метод `GET {{baseUrl}}/api/v1/cars` и выбирает нужный автомобиль и в запросе на аренду передает:

* `carUid` (UUID автомобиля) – берется из запроса `/cars`;
* `dateFrom` и `dateTo` (дата начала и конца аренды) – задается пользователем.

Система проверяет, что автомобиль с таким `carUid` существует и резервирует его (флаг `availability = false`). При
повторном вызове `GET {{baseUrl}}/api/v1/cars` этот автомобиль будет скрыт в выдаче результатов пока не будет передан
флаг `showAll = true`.

Считается количество дней аренды (`dateFrom` – `dateTo`), вычисляется общая сумма бронирования, выполняется запрос в
Payment Service и создается новая запись об оплате. В сервисе Rental Service создается запись с информацией о
бронировании.

```http request
POST {{baseUrl}}/api/v1/rental
Content-Type: application/json
X-User-Name: {{username}}

{
  "carUid": "109b42f3-198d-4c89-9276-a7520a7120ab",
  "dateFrom": "2021-10-08",
  "dateTo": "2021-10-11"
}
```

#### Завершение аренды автомобиля

* С автомобиля снимается резерв.
* В Rental Service аренда помечается завершенной (статус `FINISHED`).

```http request
POST {{baseUrl}}/api/v1/rental/{{rentalUid}}/finish
X-User-Name: {{username}}
```

#### Отмена аренды автомобиля

* С автомобиля снимается резерв.
* В Rental Service аренда помечается отмененной (статус `CANCELED`).
* В Payment Service запись об оплате помечается отмененной (статус `CANCELED`).

```http request
DELETE {{baseUrl}}/api/v1/rental/{{rentalUid}}
X-User-Name: {{username}}
```

Описание в формате [OpenAPI](%5Binst%5D%5Bv3%5D%20Car%20Rental%20System.yml).

### Данные для тестов

Создать данные для тестов:

```yaml
cars:
  – id: 1
    car_uid: "109b42f3-198d-4c89-9276-a7520a7120ab"
    brand: "Mercedes Benz"
    model: "GLA 250"
    registration_number: "ЛО777Х799"
    power: 249
    type: "SEDAN"
    price: 3500
    available: true
```