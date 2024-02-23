CREATE TABLE todo (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    deadline TEXT NOT NULL,
    is_done BOOLEAN NOT NULL DEFAULT false
);
