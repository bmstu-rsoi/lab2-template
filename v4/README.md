## Library System

Система позволяет пользователю найти интересующую книгу и взять ее в библиотеке. Если у пользователя на руках есть уже N
книг, то он не может взять новую, пока не сдал старые. Если пользователь возвращает книги в хорошем состоянии и сдает их
в срок, то максимальное количество книг у него на руках увеличивается.

### Структура Базы Данных

#### Reservation System

Сервис запускается на порту 8070.

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

Сервис запускается на порту 8060.

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
    book_id         INT REFERENCES books (id),
    library_id      INT REFERENCES library (id),
    available_count INT NOT NULL
);
```

#### Rating System

Сервис запускается на порту 8050.

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

Если передан флаг `showAll = true`, то выводить книги, которые в текущий момент недоступны для
аренды (`available_count = 0`).

```http request
GET {{baseUrl}}/api/v1/libraries/{{libraryUid}}/books&page={{page}}&size={{size}}
```

#### Получить информацию по всем взятым в прокат книгам пользователя

```http request
GET {{baseUrl}}/api/v1/reservations
X-User-Name: {{username}}
```

#### Взять книгу в библиотеке

Пользователь вызывает метод `GET {{baseUrl}}/api/v1/libraries?city={{city}}`, выбирает нужную библиотеку, вызывает
метод `GET {{baseUrl}}/api/v1/{{libraryUid}}/books` и выбирает нужную книгу для аренды.

* `bookUid` (UUID книги) – берется из запроса `/books`;
* `libraryUid` (UUID библиотеки) – берется из запроса `/libraries`;
* `tillDate` (дата окончания бронирования) – задается пользователем.

Перед выдачей книги проверяется количество книг у пользователя на руках (запрос в Reservation Service в
статусе `RENTED`). После этого выполняется запрос в Rating System и запрашивается количество звезд. Количество звезд
определяет максимальное количество книг, которые пользователь может одновременно взять в аренду.

Если условие выполнено, то создается запись в Reservation System в статусе `RENTED` и в Library System уменьшается
количество доступных книг (поле `available_count`).

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

* `condition` (состояние, в котором книгу вернули) – задается пользователем;
* `date` (дата, когда вернули книгу) – задается пользователем.

При возврате книги в Rented System изменяется статус на:

* `EXPIRED` если дата возврата больше `till_date` в записи о резерве;
* `RETURNED` если книгу сдали в срок.

Выполняется запрос в Library Service для увеличения счетчика доступных книг (поле `available_count`).

Если книгу вернули позднее срока или ее состояние на момент выдачи (запись в Reservation System) отличается от
состояния, в котором ее вернули, то у пользователя _уменьшается_ количество звезд на 10 за каждое условие (сдача позднее
срока и в плохом состоянии).

Если книгу вернули в исходном состоянии и в срок, то рейтинг пользователя _увеличивается_ на 1 звезду.

```http request
POST {{baseUrl}}/api/v1/reservations/{{reservationUid}}/return
X-User-Name: {{username}}

{
  "condition": "EXCELLENT",
  "date": "2021-10-11"
}
```

#### Получить рейтинг пользователя

Количество звезд пользователя определяем максимальное количество одновременно арендованных книг. У пользователя может
быть от 1 до 100 звезд, если изменение выходит за эти границы, то устанавливается граничное значение.

```http request
GET {{baseUrl}}/api/v1/rating
X-User-Name: {{username}}
```

Описание в формате [Open API](%5Binst%5D%5Bv4%5D%20Library%20System.yml).

### Данные для тестов

Создать данные для тестов:

```yaml
library:
  – id: 1
    library_uid: "83575e12-7ce0-48ee-9931-51919ff3c9ee"
    name: "Библиотека имени 7 Непьющих"
    city: "Москва"
    address: "2-я Бауманская ул., д.5, стр.1"

books:
  - id: 1
    book_uid: "f7cdc58f-2caf-4b15-9727-f89dcc629b27",
    name: "Краткий курс C++ в 7 томах",
    author: "Бьерн Страуструп",
    genre: "Научная фантастика",
    condition: "EXCELLENT",

library_books:
  - book_id: 1
    library_id: 1
    available_count: 1
```