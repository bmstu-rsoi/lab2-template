## Library System

Система позволяет пользователю найти интересующую книгу и взять ее в библиотеке. Если у пользователя на руках есть уже N
книг, то он не может взять новую, пока не сдал старые. Если пользователь возвращает книги в хорошем состоянии и сдает их
в срок, то максимальное количество книг у него на руках увеличивается.

### Структура Базы Данных

#### Reservation System

```sql
CREATE TABLE reservation
(
    id              SERIAL PRIMARY KEY,
    reservation_uid uuid UNIQUE NOT NULL,
    username        VARCHAR(80) NOT NULL,
    book_uid        uuid        NOT NULL,
    library_uid     uuid        NOT NULL,
    status          VARCHAR(20) NOT NULL
        CHECK (status IN ('RENTED', 'RETURNED', 'EXPIRED')),
    start_date      TIMESTAMP   NOT NULL,
    till_date       TIMESTAMP   NOT NULL
)
```

#### Library System

```sql
CREATE TABLE library
(
    id          SERIAL PRIMARY KEY,
    library_uid uuid UNIQUE  NOT NULL,
    name        VARCHAR(80)  NOT NULL,
    city        VARCHAR(255) NOT NULL,
    address     VARCHAR(255) NOT NULL
);

CREATE TABLE books
(
    id        SERIAL PRIMARY KEY,
    book_uid  uuid UNIQUE  NOT NULL,
    name      VARCHAR(255) NOT NULL,
    author    VARCHAR(255),
    genre     VARCHAR(255),
    condition VARCHAR(20) DEFAULT 'EXCELLENT'
        CHECK (condition IN ('EXCELLENT', 'GOOD', 'BAD'))
);

CREATE TABLE library_books
(
    book_id    INT REFERENCES books (id),
    library_id INT REFERENCES library (id)
);
```

#### Rating System

```sql
CREATE TABLE rating
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(80) NOT NULL,
    stars    INT         NOT NULL
        CHECK (stars BETWEEN 0 AND 100)
);
```

### Описание API

#### Получить список библиотек в городе

```http request
GET {{baseUrl}}/api/v1/libraries?city={{city}}&page={{page}}&size={{size}}
```

#### Получить список книг в выбранной библиотеке

```http request
GET {{baseUrl}}/api/v1/libraries/{{libraryUid}}/books&page={{page}}&size={{size}}
```

#### Получить информацию по всем взятым в прокат книгам пользователя

```http request
GET {{baseUrl}}/api/v1/reservations
X-User-Name: {{username}}
```

#### Взять книгу в библиотеке

```http request
POST {{baseUrl}}/api/v1/reservations
Content-Type: application/json
X-User-Name: {{username}}

{
  "bookUid": "f7cdc58f-2caf-4b15-9727-f89dcc629b27",
  "libraryUid": "83575e12-7ce0-48ee-9931-51919ff3c9ee",
  "tillDate": "2021-10-11"
}
```

#### Вернуть книгу

```http request
POST {{baseUrl}}/api/v1/reservations/{{reservationUid}}/return
X-User-Name: {{username}}

{
  "condition": "EXCELLENT",
  "date": "2021-10-11"
}
```

#### Получить рейтинг пользователя

```http request
GET {{baseUrl}}/api/v1/rating
X-User-Name: {{username}}
```

Описание в формате [Open API](%5Binst%5D%5Bv4%5D%20Library%20System.yml).