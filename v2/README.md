## Hotels Booking System

Система предоставляет пользователю сервис поиска и бронирования отелей на интересующие даты. В зависимости от количества
заказов система лояльности дает скидку пользователю на новые бронирования.

### Reservation Service

```sql
CREATE TABLE reservation
(
    id              SERIAL PRIMARY KEY,
    reservation_uid uuid UNIQUE NOT NULL,
    username        VARCHAR(80) NOT NULL,
    payment_uid     uuid        NOT NULL,
    hotel_id        INT REFERENCES hotels (id),
    status          VARCHAR(20) NOT NULL
        CHECK (status IN ('PAID', 'RESERVED', 'CANCELED')),
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
    stars     INT
);
```

### Payment Service

```sql
CREATE TABLE payment
(
    id          SERIAL PRIMARY KEY,
    payment_uid uuid        NOT NULL,
    status      VARCHAR(20) NOT NULL
        CHECK (status IN ('PAID', 'REVERSED', 'CANCELED')),
    price       INT         NOT NULL
);
```

### Loyalty Service

```sql
CREATE TABLE loyalty
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(80) NOT NULL UNIQUE,
    status   VARCHAR(80) NOT NULL DEFAULT 'BRONZE'
        CHECK (status IN ('BRONZE', 'SILVER', 'GOLD')),
    discount INT         NOT NULL
);
```

[OpenAPI](%5Binst%5D%5Bv2%5D%20Hotels%20Booking%20System.yml)