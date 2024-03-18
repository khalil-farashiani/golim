CREATE TABLE role (
                         id                 INTEGER PRIMARY KEY AUTO_INCREMENT,
                         endpoint           varchar(255)      NOT NULL,
                         operation          varchar(255) NOT NULL,
                         bucket_size        INTEGER NOT NULL,
                         created_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                         add_token_per_min  INTEGER NOT NULL,
                         initial_tokens     INTEGER NOT NULL,
                         updated_at         TIMESTAMP,
                         deleted_at         TIMESTAMP,
                         rate_limiter_id    INTEGER NOT NULL,
                         FOREIGN KEY (rate_limiter_id)
                             REFERENCES rate_limiter (id)
);


CREATE TABLE rate_limiter (
                      id                 INTEGER PRIMARY KEY AUTO_INCREMENT,
                      name               varchar(255)      NOT NULL,
                      created_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                      updated_at         TIMESTAMP,
                      deleted_at         TIMESTAMP
);