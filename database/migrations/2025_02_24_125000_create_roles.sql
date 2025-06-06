CREATE TABLE roles (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(50) NOT NULL UNIQUE,
    description     VARCHAR(255) NULL,
    created_at      TIMESTAMPTZ DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at      TIMESTAMPTZ DEFAULT (NOW() AT TIME ZONE 'UTC')
);

CREATE OR REPLACE FUNCTION update_updated_at_roles()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = (NOW() AT TIME ZONE 'UTC');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_roles_updated_at
BEFORE UPDATE ON roles
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_roles();
