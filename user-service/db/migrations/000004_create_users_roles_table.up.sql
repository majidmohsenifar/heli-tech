CREATE TABLE IF NOT EXISTS users_roles (
    user_id BIGINT REFERENCES users(id),
    role_id integer REFERENCES roles(id),
   PRIMARY KEY (user_id, role_id)
);
