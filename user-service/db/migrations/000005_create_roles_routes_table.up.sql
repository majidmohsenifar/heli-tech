CREATE TABLE IF NOT EXISTS roles_routes (
    role_id integer REFERENCES roles(id),
    route_id integer REFERENCES routes(id),
   PRIMARY KEY (role_id,route_id)
);

