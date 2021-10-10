## Flight Booking System

Система предоставляет пользователю возможность поиска и покупки билетов. В зависимости от количества выполненных
перелетов, пользователю предоставляется скидка на перелет и начисляются баллы, которые он может использовать для оплаты.

### Структура Базы Данных

#### Ticket Service

```sql
CREATE TABLE ticket
(
    id         SERIAL PRIMARY KEY,
    ticket_uid uuid UNIQUE NOT NULL,
    username   VARCHAR(80) NOT NULL,
    flight_uid uuid        NOT NULL,
    price      INT         NOT NULL,
    status     VARCHAR(20) NOT NULL
        CHECK (status IN ('PAID', 'CANCELED'))
);
```

#### Flight Service

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

#### Bonus Service

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
    id             SERIAL      NOT NULL,
    privilege_id   INT REFERENCES privilege (id),
    ticket_uid     uuid        NOT NULL,
    datetime       TIMESTAMP   NOT NULL,
    balance_diff   INT         NOT NULL,
    operation_type VARCHAR(20) NOT NULL
        CHECK (operation_type IN ('FILL_IN_BALANCE', 'DEBIT_THE_ACCOUNT', 'FILLED_BY_MONEY'))
);
```

### Описание API

#### Получить список всех перелетов

```http request
GET {{baseUrl}}/api/v1/flights
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

* `flight_number` (номер рейса) – берется из запроса `/flights`;
* `date` (дата перелета) – вводится пользователем;
* `price` (цена) – берется из запроса `/flights`;
* `paid_from_balance` (оплата бонусами) – флаг, указывающий, что для оплаты билета нужно использовать бонусный счет.

Если при покупке указан флаг `"paid_from_balance": true`, то с бонусного счёта списываются максимальное количество
баллов в отношении 1 балл – 1 рубль.

Т.е. если на бонусном счете было 500 бонусов, билет стоит 1500 рублей и при покупке был указан
флаг `"paid_from_balance": true"`, то со счёта спишется 500 бонусов (в ответе будет указано `"paid_by_bonuses": 500`), а
стоимость билета будет 1000 рублей (в ответе будет указано `"paid_by_money": 1000`). В сервисе Bonus Service в
таблицу `privilege_history` будет добавлена запись о списании со счёта 500 бонусов.

Если при покупке был указан флаг `"paid_from_balance": false`, то в ответе будет `"paid_by_bonuses": 0`, а на бонусный
счет будет начислено бонусов в размере 10% от стоимости заказа. Так же в таблицу `privilege_history` будет добавлена
запись о зачислении бонусов.

```http request
POST {{baseUrl}}/api/v1/tickets
Content-Type: application/json
X-User-Name: {{username}}

{
  "flight_number": "AFL031",
  "date": "2021-10-08T19:59:19Z",
  "price": 1500,
  "paid_from_balance": true
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