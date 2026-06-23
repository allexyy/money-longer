CREATE TABLE vaults
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    limit_amount       INT          NOT NULL,
    left_amount INT          NOT NULL,
    icon        VARCHAR(10),
    color       VARCHAR(20),
    expire      TIMESTAMP    NOT NULL
);

CREATE TABLE transactions
(
    id        SERIAL PRIMARY KEY,
    name      VARCHAR(255) NOT NULL,
    amount    INT          NOT NULL,
    is_income BOOLEAN      NOT NULL DEFAULT FALSE,
    vault_id  INT          REFERENCES vaults (id) ON DELETE SET NULL,
    note      TEXT,
    date      TIMESTAMP    NOT NULL DEFAULT NOW()
);