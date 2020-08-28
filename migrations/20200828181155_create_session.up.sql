CREATE TABLE "payment" (
    id UUID NOT NULL PRIMARY KEY,
    card varchar NOT NULL,
    time time NOT NULL,
    game_name varchar NOT NULL,
    user_email varchar NOT NULL,
    seller_email varchar NOT NULL
);