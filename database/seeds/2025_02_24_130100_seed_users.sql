USE gohst;

INSERT INTO users (firstname, lastname, email, password_hash, role_id, active) VALUES
    ('Admin', 'User', 'admin@example.com', '$2a$10$your_hashed_password', (SELECT id FROM roles WHERE name = 'admin'), 1),
    ('Test', 'User', 'test@example.com', '$2a$10$your_hashed_password', (SELECT id FROM roles WHERE name = 'user'), 1);
