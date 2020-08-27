CREATE TYPE user_role as ENUM (
    'seller'
);

CREATE TABLE "user_roles" (
    user_id UUID NOT NULL REFERENCES "user",
    role user_role NOT NULL,
    primary key (user_id, role)
);

CREATE INDEX user_roles_user
    ON "user_roles" (user_id);