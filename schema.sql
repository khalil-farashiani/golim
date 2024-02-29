CREATE TABLE role (
                         id                 INTEGER PRIMARY KEY AUTO_INCREMENT,
                         endpoint           varchar(255)      NOT NULL,
                         operation          varchar(255) NOT NULL,
                         bucket_size        INTEGER NOT NULL,
                         add_token_per_sec  INTEGER NOT NULL,
                         initial_tokens     INTEGER NOT NULL,
                         created_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                         updated_at         TIMESTAMP,
                         deleted_at         TIMESTAMP
);