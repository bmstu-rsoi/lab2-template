## Library System

Система позволяет пользователю найти интересующую книгу и взять ее в библиотеке. Если у пользователя на руках есть уже N
книг, то он не может взять новую, пока не сдал старые. Если пользователь возвращает книги в хорошем состоянии и сдает их
в срок, то максимальное количество книг у него на руках увеличивается.

### Reservation System

```sql
CREATE TABLE reservation
(
    id        SERIAL PRIMARY KEY,
    username  VARCHAR(80) NOT NULL,
    book_uid  uuid        NOT NULL,
    till_date TIMESTAMP   NOT NULL
)
```

### Library System

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

### Rating System

```sql
CREATE TABLE rating
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(80) NOT NULL,
    stars    INT         NOT NULL
        CHECK (stars BETWEEN 1 AND 100)
);
```

[Open API](%5Binst%5D%5Bv4%5D%20Library%20System.yml)