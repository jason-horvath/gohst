-- Password hash for 'Test1234!' using argon
INSERT INTO users (firstname, lastname, email, password_hash, role_id, active) VALUES
    ('Admin', 'User', 'admin@example.com', '$argon2id$v=19$m=65536,t=4,p=2$E/ke48n/idmA7oeI3sI9Pg$8mC3W7VHTOHhlb95bEc+LFtU36m1UGp6myy6A30Em5g',
      (SELECT id FROM roles WHERE name = 'admin'), TRUE),
    ('Test', 'Manager', 'manager@example.com', '$argon2id$v=19$m=65536,t=4,p=2$E/ke48n/idmA7oeI3sI9Pg$8mC3W7VHTOHhlb95bEc+LFtU36m1UGp6myy6A30Em5g',
      (SELECT id FROM roles WHERE name = 'manager'), TRUE),
    ('Regular', 'User', 'user@example.com', '$argon2id$v=19$m=65536,t=4,p=2$E/ke48n/idmA7oeI3sI9Pg$8mC3W7VHTOHhlb95bEc+LFtU36m1UGp6myy6A30Em5g',
      (SELECT id FROM roles WHERE name = 'user'), TRUE);
