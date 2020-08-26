CREATE TABLE "user" (
    user_id UUID not null primary key,
    email varchar not null,
    encrypted_password varchar not null
);

CREATE UNIQUE INDEX user_email
    ON "user" (email);