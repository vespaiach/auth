package migrate

var initProdData = `
INSERT INTO actions (id, action_name, action_desc) VALUES (1, 'create_action', 'Create a action');
INSERT INTO actions (id, action_name, action_desc) VALUES (2, 'update_action', 'Update a action');
INSERT INTO actions (id, action_name, action_desc) VALUES (3, 'get_action', 'Get a action');
INSERT INTO actions (id, action_name, action_desc) VALUES (4, 'query_action', 'Query actions');

INSERT INTO actions (id, action_name, action_desc) VALUES (5, 'create_role', 'Create a role');
INSERT INTO actions (id, action_name, action_desc) VALUES (6, 'update_role', 'Update a role');
INSERT INTO actions (id, action_name, action_desc) VALUES (7, 'get_role', 'Get a role');
INSERT INTO actions (id, action_name, action_desc) VALUES (8, 'query_role', 'Query roles');

INSERT INTO actions (id, action_name, action_desc) VALUES (9, 'create_user', 'Create a user');
INSERT INTO actions (id, action_name, action_desc) VALUES (10, 'update_user', 'Update a user');
INSERT INTO actions (id, action_name, action_desc) VALUES (11, 'get_user', 'Get a user');
INSERT INTO actions (id, action_name, action_desc) VALUES (12, 'query_user', 'Query users');

INSERT INTO actions (id, action_name, action_desc) VALUES (13, 'create_role_action', 'Create a role_action');
INSERT INTO actions (id, action_name, action_desc) VALUES (14, 'delete_role_action', 'Delete a role_action');
INSERT INTO actions (id, action_name, action_desc) VALUES (15, 'get_role_action', 'Get a role_action');
INSERT INTO actions (id, action_name, action_desc) VALUES (16, 'query_role_action', 'Query role_actions');

INSERT INTO actions (id, action_name, action_desc) VALUES (17, 'create_user_role', 'Create a user_role');
INSERT INTO actions (id, action_name, action_desc) VALUES (18, 'delete_user_role', 'Delete a user_role');
INSERT INTO actions (id, action_name, action_desc) VALUES (19, 'get_user_role', 'Get a user_role');
INSERT INTO actions (id, action_name, action_desc) VALUES (20, 'query_user_role', 'Query user_roles');

INSERT INTO roles (id, role_name, role_desc) VALUES (1, 'admin', 'Administrator');

INSERT INTO users (id, full_name, username, email, hashed, verified) VALUES (1, '%s', '%s', '%s', '%s', 1);

INSERT INTO role_actions (role_id, action_id) VALUES (1, 1);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 2);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 3);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 4);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 5);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 6);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 7);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 8);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 9);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 10);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 11);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 12);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 13);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 14);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 15);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 16);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 17);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 18);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 19);
INSERT INTO role_actions (role_id, action_id) VALUES (1, 20);

INSERT INTO user_roles (user_id, role_id) VALUES (1, 1);
`
