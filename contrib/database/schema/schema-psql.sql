CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL UNIQUE,
    password_hash VARCHAR(256) NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL UNIQUE
);

CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    category VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    resource VARCHAR(100) NOT NULL,
    UNIQUE(resource, action)
);

INSERT INTO permissions (id, category, resource, action)
VALUES
    (1, 'system', 'health', 'get'),
    (2, 'auth', 'user', 'add'),
    (3, 'auth', 'group', 'add'),
    (4, 'auth', 'usergroup', 'add'),
    (5, 'auth', 'grouppermission', 'add'); 

CREATE TABLE user_groups (
    user_id  INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, group_id)
);

CREATE TABLE group_permissions (
    group_id      INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, permission_id)
);
