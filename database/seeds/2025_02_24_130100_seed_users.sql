USE gohst;

-- Password hash for 'admin123' using bcrypt
INSERT INTO users (firstname, lastname, email, password_hash, role_id, active) VALUES
    ('Admin', 'User', 'admin@example.com', '$2a$10$xVF19REj6G9o1bLH4luR2.tRhPzs4KZoQzsMCj14OUo8taeIkaYLy',
      (SELECT id FROM roles WHERE name = 'admin'), 1),
    ('Test', 'Manager', 'manager@example.com', '$2a$10$xVF19REj6G9o1bLH4luR2.tRhPzs4KZoQzsMCj14OUo8taeIkaYLy',
      (SELECT id FROM roles WHERE name = 'manager'), 1),
    ('Regular', 'User', 'user@example.com', '$2a$10$xVF19REj6G9o1bLH4luR2.tRhPzs4KZoQzsMCj14OUo8taeIkaYLy',
      (SELECT id FROM roles WHERE name = 'user'), 1);
