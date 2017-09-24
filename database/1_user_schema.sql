/* Application database schema
 * RDBMS: PostgreSQL 9.5
 *
 * This file defines the schema for user accounts. It relies on the core schema.
 */
CREATE EXTENSION pgcrypto;

CREATE TABLE roles (
    id serial NOT NULL PRIMARY KEY,
    role text NOT NULL UNIQUE
);
INSERT INTO roles (role) VALUES
('admin'), ('user');

CREATE TABLE users (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    role_id int NOT NULL REFERENCES roles(id),
    name text NOT NULL UNIQUE,
    target_language_id int REFERENCES languages(id),
    password text NOT NULL
);

/* create_user creates a new user with a default 'user' role.
 *
 * params:
 *      _name should be unique across all users
 *      _password should be plaintext. It will be hashed and salted here.
 *
 * returns: the id (a uuid) of the created user
 * raises:  an exception if _name is not unique
 */
CREATE OR REPLACE FUNCTION create_user(_name TEXT, _password TEXT)
RETURNS uuid AS $$
DECLARE
    _user_id uuid;
BEGIN
    INSERT INTO users (role_id, name, password)
    VALUES (
        (SELECT id FROM roles WHERE role = 'user'),
        _name,
        crypt(_password, gen_salt('bf', 8))
    )
    RETURNING id INTO _user_id;
    RETURN _user_id;
EXCEPTION
    WHEN unique_violation THEN
        RAISE EXCEPTION 'user % already exists', _name;
END;
$$ LANGUAGE plpgsql;
