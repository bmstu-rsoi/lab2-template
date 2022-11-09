## Flight Booking System

Система предоставляет пользователю возможность поиска и покупки билетов. При покупке билетов пользователю начисляются
баллы, которые он может использовать для оплаты.

### Структура Базы Данных

#### Ticket Service

Сервис запускается на порту 8070.

```sql
CREATE TABLE ticket
(
    id            SERIAL PRIMARY KEY,
    ticket_uid    uuid UNIQUE NOT NULL,
    username      VARCHAR(80) NOT NULL,
    flight_number VARCHAR(20) NOT NULL,
    price         INT         NOT NULL,
    status        VARCHAR(20) NOT NULL
        CHECK (status IN ('PAID', 'CANCELED'))
);
```

#### Flight Service

Сервис запускается на порту 8060.

```sql
CREATE TABLE flight
(
    id              SERIAL PRIMARY KEY,
    flight_number   VARCHAR(20)              NOT NULL,
    datetime        TIMESTAMP WITH TIME ZONE NOT NULL,
    from_airport_id INT REFERENCES airport (id),
    to_airport_id   INT REFERENCES airport (id),
    price           INT                      NOT NULL
);

CREATE TABLE airport
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(255),
    city    VARCHAR(255),
    country VARCHAR(255)
);
```

```json
{
  "balance": 1500,
  "status": "GOLD",
  "history": [
    {
      "date": "2021-10-08T19:59:19Z",
      "ticketUid": "049161bb-badd-4fa8-9d90-87c9a82b0668",
      "balanceDiff": 1500,
      "operationType": "FILL_IN_BALANCE"
    }
  ]
}
```
#### Bonus Service

Сервис запускается на порту 8050.

```sql
CREATE TABLE privilege
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(80) NOT NULL UNIQUE,
    status   VARCHAR(80) NOT NULL DEFAULT 'BRONZE'
        CHECK (status IN ('BRONZE', 'SILVER', 'GOLD')),
    balance  INT
);

CREATE TABLE privilege_history
(
    id             SERIAL PRIMARY KEY,
    privilege_id   INT REFERENCES privilege (id),
    ticket_uid     uuid        NOT NULL,
    datetime       TIMESTAMP   NOT NULL,
    balance_diff   INT         NOT NULL,
    operation_type VARCHAR(20) NOT NULL
        CHECK (operation_type IN ('FILL_IN_BALANCE', 'DEBIT_THE_ACCOUNT'))
);
```

### Описание API

#### Получить список всех перелетов

```http request
GET {{baseUrl}}/api/v1/flights&page={{page}}&size={{size}}
```

#### Получить полную информацию о пользователе

Возвращается информация о билетах и статусе в системе привилегии.

```http request
GET {{baseUrl}}/api/v1/me
X-User-Name: {{username}}
```

#### Получить информацию о всех билетах пользователя

```http request
GET {{baseUrl}}/api/v1/tickets
X-User-Name: {{username}}
```

#### Получить информацию по конкретному билету пользователя

При запросе требуется проверить, что билет принадлежит пользователю.

```http request
GET {{baseUrl}}/api/v1/tickets/{{ticketUid}}
X-User-Name: {{username}}
```

#### Покупка билета

Пользователь вызывает метод `GET {{baseUrl}}/api/v1/flights` выбирает нужный рейс и в запросе на покупку передает:

* `flightNumber` (номер рейса) – берется из запроса `/flights`;
* `price` (цена) – берется из запроса `/flights`;
* `paidFromBalance` (оплата бонусами) – флаг, указывающий, что для оплаты билета нужно использовать бонусный счет.

Система проверяет, что рейс с таким номером существует. Считаем что на рейсе бесконечное количество мест.

Если при покупке указан флаг `"paidFromBalance": true`, то с бонусного счёта списываются максимальное количество баллов
в отношении 1 балл – 1 рубль.

Т.е. если на бонусном счете было 500 бонусов, билет стоит 1500 рублей и при покупке был указан
флаг `"paidFromBalance": true"`, то со счёта спишется 500 бонусов (в ответе будет указано `"paidByBonuses": 500`), а
стоимость билета будет 1000 рублей (в ответе будет указано `"paidByMoney": 1000`). В сервисе Bonus Service в
таблицу `privilegeHistory` будет добавлена запись о списании со счёта 500 бонусов.

Если при покупке был указан флаг `"paidFromBalance": false`, то в ответе будет `"paidByBonuses": 0`, а на бонусный счет
будет начислено бонусов в размере 10% от стоимости заказа. Так же в таблицу `privilegeHistory` будет добавлена запись о
зачислении бонусов.

```http request
POST {{baseUrl}}/api/v1/tickets
Content-Type: application/json
X-User-Name: {{username}}

{
  "flightNumber": "AFL031",
  "price": 1500,
  "paidFromBalance": true
}
```

#### Возврат билета

Билет помечается статусом `CANCELED`, в Bonus Service в зависимости от типа операции выполняется возврат бонусов на счёт
или списание ранее начисленных. При списании бонусный счет не может стать меньше 0.

```http request
DELETE {{baseUrl}}/api/v1/tickets/{{ticketUid}}
X-User-Name: {{username}}
```

#### Получить информацию о состоянии бонусного счета

Пользователю возвращается информация о бонусном счете и истории его изменения.

```http request
GET http://localhost:8080/api/v1/privilege
X-User-Name: {{username}}
```

Описание в формате [OpenAPI](%5Binst%5D%5Bv1%5D%20Flight%20Booking%20System.yml).

### Данные для тестов

Создать данные для тестов:

```yaml
airport:
  – id: 1
    name: Шереметьево
    city: Москва
    country: Россия
  - id: 2
    name: Пулково
    city: Санкт-Петербург
    coutry: Россия

flight:
  - id: 1
    flight_number: "AFL031"
    datetime: "2021-10-08 20:00"
    from_airport_id: 2
    to_airport_id: 1
    price: 1500
```
