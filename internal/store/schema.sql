CREATE TABLE role (
                         id                 INTEGER PRIMARY KEY AUTOINCREMENT,
                         endpoint           varchar(255)      NOT NULL,
                         Operation          varchar(255) NOT NULL,
                         bucket_size        INTEGER NOT NULL,
                         add_token_per_min  INTEGER NOT NULL,
                         initial_tokens     INTEGER NOT NULL,
                         deleted_at         TIMESTAMP,
                         rate_limiter_id    INTEGER NOT NULL,
                         FOREIGN KEY (rate_limiter_id)
                             REFERENCES rate_limiter (id),
                         CONSTRAINT endpoint_operation_unique
                             UNIQUE (endpoint, Operation) ON CONFLICT REPLACE
);


CREATE TABLE rate_limiter (
                      id                 INTEGER PRIMARY KEY AUTOINCREMENT,
                      Name               varchar(255)      NOT NULL,
                      Destination        varchar(255)      NOT NULL,
                      deleted_at         TIMESTAMP,
                      CONSTRAINT name_unique UNIQUE (Name, Destination)
);