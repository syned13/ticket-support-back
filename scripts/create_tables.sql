CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    user_type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);


CREATE TABLE IF NOT EXISTS tickets (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    decription TEXT NOT NULL,
    ticket_type TEXT NOT NULL,
    severity INT NOT NULL,
    ticket_priority INT NOT NULL,
    ticket_status TEXT NOT NULL,
    creator_id TEXT NOT NULL REFERENCES users (id),
    owner_id TEXT REFERENCES users (id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    resolved_at TIMESTAMP
);