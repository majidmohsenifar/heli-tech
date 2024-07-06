INSERT INTO roles_routes(role_id,route_id) 
SELECT roles.id AS role_id, routes.id AS route_id FROM roles JOIN routes ON roles.id >0;
