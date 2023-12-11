CREATE TABLE IF NOT EXISTS url
(
    alias      VARCHAR(10) PRIMARY KEY,
    original   VARCHAR(255) NOT NULL,
    expires_at DATE         NOT NULL,
    visits     INTEGER
);

