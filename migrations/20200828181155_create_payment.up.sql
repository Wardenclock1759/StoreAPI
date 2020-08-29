CREATE TABLE "payment" (
    id UUID NOT NULL PRIMARY KEY,
    time time NOT NULL,
    game_name varchar NOT NULL,
    user_email varchar NOT NULL,
    seller_email varchar NOT NULL,
    total int NOT NULL,
    code varchar NOT NULL
);