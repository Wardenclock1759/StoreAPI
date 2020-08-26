CREATE TABLE "game" (
    game_id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES "user",
    Name varchar NOT NULL UNIQUE,
    Price int NOT NULL
);

CREATE INDEX game_user
    ON "game" (user_id);