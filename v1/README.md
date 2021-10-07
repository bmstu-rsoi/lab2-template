## Flight Booking System

Система предоставляет пользователю возможность поиска и покупки билетов. В зависимости от количества выполненных
перелетов, пользователю предоставляется скидка на перелет и начисляются баллы, которые он может использовать для оплаты.

### Ticket Service

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

### Flight Service

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

### Bonus Service

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
    datetime       TIMESTAMP   NOT NULL,
    balance_diff   INT         NOT NULL,
    operation_type VARCHAR(20) NOT NULL
        CHECK (operation_type IN ('FILL_IN_BALANCE', 'DEBIT_THE_ACCOUNT', 'FILLED_BY_MONEY'))
);
```

[OpenAPI](%5Binst%5D%5Bv1%5D%20Flight%20Booking%20System.yml)