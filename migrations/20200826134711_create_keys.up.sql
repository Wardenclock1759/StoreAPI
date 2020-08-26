CREATE TABLE game_code (
    game_id UUID NOT NULL REFERENCES game,
    code varchar NOT NULL
);

CREATE INDEX game_code_game
    ON game_code (game_id);