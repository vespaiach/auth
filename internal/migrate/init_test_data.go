package migrate

var initTestData = `
INSERT INTO actions (id, action_name, action_desc) VALUES (1, 'create_action', 'Create a action');
INSERT INTO actions (id, action_name, action_desc) VALUES (2, 'delete_action', 'Delete a action');
INSERT INTO actions (id, action_name, action_desc) VALUES (3, 'update_action', 'Update a action');
INSERT INTO actions (id, action_name, action_desc) VALUES (4, 'view_action', 'View a action');
INSERT INTO actions (id, action_name, action_desc) VALUES (5 ,'list_action', 'List actions');
INSERT INTO actions (id, action_name, action_desc) VALUES (6, 'create_user', 'Create a user');
INSERT INTO actions (id, action_name, action_desc) VALUES (7, 'delete_user', 'Delete a user');
INSERT INTO actions (id, action_name, action_desc) VALUES (8, 'update_user', 'Update a user');
INSERT INTO actions (id, action_name, action_desc) VALUES (9, 'view_user', 'View a user');
INSERT INTO actions (id, action_name, action_desc) VALUES (10 ,'list_user', 'List users');

INSERT INTO roles (id, role_name, role_desc) VALUES (1, 'admin_role', 'Admin role');
INSERT INTO roles (id, role_name, role_desc) VALUES (2, 'staff_role', 'Staff role');

INSERT INTO role_actions (id, role_id, action_id) VALUES (1, 1, 1);
INSERT INTO role_actions (id, role_id, action_id) VALUES (2, 1, 2);
INSERT INTO role_actions (id, role_id, action_id) VALUES (3, 1, 3);
INSERT INTO role_actions (id, role_id, action_id) VALUES (4, 1, 4);
INSERT INTO role_actions (id, role_id, action_id) VALUES (5, 1, 5);
INSERT INTO role_actions (id, role_id, action_id) VALUES (6, 1, 6);
INSERT INTO role_actions (id, role_id, action_id) VALUES (7, 1, 7);
INSERT INTO role_actions (id, role_id, action_id) VALUES (8, 1, 8);
INSERT INTO role_actions (id, role_id, action_id) VALUES (9, 1, 9);
INSERT INTO role_actions (id, role_id, action_id) VALUES (10, 1, 10);

INSERT INTO role_actions (id, role_id, action_id) VALUES (11, 2, 1);
INSERT INTO role_actions (id, role_id, action_id) VALUES (12, 2, 2);
INSERT INTO role_actions (id, role_id, action_id) VALUES (13, 2, 3);
INSERT INTO role_actions (id, role_id, action_id) VALUES (14, 2, 4);
INSERT INTO role_actions (id, role_id, action_id) VALUES (15, 2, 5);

INSERT INTO users (id, full_name, username, hashed, email) VALUES (1, 'full name' ,'admin', 'password', 'admin@test.com');
INSERT INTO users (id, full_name, username, hashed, email) VALUES (2, 'full name' ,'staff', 'password', 'staff@test.com');

INSERT INTO user_roles (id, user_id, role_id) VALUES (1, 1, 1);
INSERT INTO user_roles (id, user_id, role_id) VALUES (2, 1, 2);
INSERT INTO user_roles (id, user_id, role_id) VALUES (3, 2, 2);

INSERT INTO user_actions (id, user_id, action_id) VALUES (1, 2, 10);
`
