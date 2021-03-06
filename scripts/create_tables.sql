CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    user_type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

INSERT INTO users 
(name, email, password, user_type, created_at)
VALUES ('Erica Ross', 'erica@erica.com', '$2y$12$ZPmLiyARMnzTFZuvhj42y.7PyPh5TVQfvu4IGpPFOopAs4c9rA1km', 'admin', NOW());

INSERT INTO users 
(name, email, password, user_type, created_at)
VALUES ('Denys Rosario', 'denys@denys.com', '$2y$12$ZPmLiyARMnzTFZuvhj42y.7PyPh5TVQfvu4IGpPFOopAs4c9rA1km', 'admin', NOW());

INSERT INTO users 
(name, email, password, user_type, created_at)
VALUES ('Angelica Pena', 'angelica@angelica.com', '$2y$12$ZPmLiyARMnzTFZuvhj42y.7PyPh5TVQfvu4IGpPFOopAs4c9rA1km', 'admin', NOW());

INSERT INTO users 
(name, email, password, user_type, created_at)
VALUES ('Leiscar Trinidad', 'leiscar@leiscar.com', '$2y$12$ZPmLiyARMnzTFZuvhj42y.7PyPh5TVQfvu4IGpPFOopAs4c9rA1km', 'admin', NOW());


CREATE TABLE IF NOT EXISTS tickets (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    ticket_description TEXT NOT NULL,
    ticket_type TEXT NOT NULL,
    severity INT NOT NULL,
    ticket_priority INT NOT NULL,
    ticket_status TEXT NOT NULL,
    creator_id INT NOT NULL REFERENCES users (id),
    owner_id INT REFERENCES users (id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    resolved_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tickets_changes (
    id SERIAL PRIMARY KEY,
    ticket_id INT NOT NULL REFERENCES tickets (id),
    creator_id INT NOT NULL REFERENCES users (id),
    to_status TEXT NOT NULL,
    changed_at TIMESTAMP NOT NULL
);