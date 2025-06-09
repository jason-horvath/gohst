-- Password hash for 'admin123' using bcrypt
INSERT INTO users (firstname, lastname, email, password_hash, role_id, active) VALUES
    ('Admin', 'User', 'admin@example.com', '$argon2id$v=19$m=65536,t=4,p=2$d9hd9l/GVDGIfsA9k1U+Sw$7iyxzBwdWcGfbNp+NtOlYCwTarhQiouSXcAccssDEeo',
      (SELECT id FROM roles WHERE name = 'admin'), TRUE),
    ('Test', 'Manager', 'manager@example.com', '$argon2id$v=19$m=65536,t=4,p=2$d9hd9l/GVDGIfsA9k1U+Sw$7iyxzBwdWcGfbNp+NtOlYCwTarhQiouSXcAccssDEeo',
      (SELECT id FROM roles WHERE name = 'manager'), TRUE),
    ('Regular', 'User', 'user@example.com', '$argon2id$v=19$m=65536,t=4,p=2$d9hd9l/GVDGIfsA9k1U+Sw$7iyxzBwdWcGfbNp+NtOlYCwTarhQiouSXcAccssDEeo',
      (SELECT id FROM roles WHERE name = 'user'), TRUE);
