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

INSERT INTO users (id, full_name, username, hashed, email) VALUES (1, 'full name' ,'admin', '$2a$10$88y3eBfBma0lQzgEhPg7m.xmZQUE5DhcHqewtz0UvIYIfFZQFnD/G', 'admin@test.com');
INSERT INTO users (id, full_name, username, hashed, email) VALUES (2, 'full name' ,'staff', '$2a$10$88y3eBfBma0lQzgEhPg7m.xmZQUE5DhcHqewtz0UvIYIfFZQFnD/G', 'staff@test.com');

INSERT INTO user_roles (id, user_id, role_id) VALUES (1, 1, 1);
INSERT INTO user_roles (id, user_id, role_id) VALUES (2, 1, 2);
INSERT INTO user_roles (id, user_id, role_id) VALUES (3, 2, 2);

INSERT INTO user_actions (id, user_id, action_id) VALUES (1, 2, 10);

INSERT INTO token_histories (uid, user_id, access_token, refresh_token, remote_addr, x_forwarded_for, x_real_ip, user_agent, created_at, expired_at) VALUES ('aab5d5fd-70c1-11e5-a4fb-b026b977eb28', 1, 'access_token1', 'refresh_token', 'remote_addr', 'x_forwarded_for', 'x_real_ip', 'user_agent', '2019-08-31 20:23:00', '2019-08-31 23:00:00');
INSERT INTO token_histories (uid, user_id, access_token, refresh_token, remote_addr, x_forwarded_for, x_real_ip, user_agent, created_at, expired_at) VALUES ('e65bedae-c17f-11e9-bf92-0242ac120002', 1, 'access_token2', 'refresh_token', 'remote_addr', 'x_forwarded_for', 'x_real_ip', 'user_agent', '2019-08-31 20:23:00', '2019-08-31 23:00:00');
INSERT INTO token_histories (uid, user_id, access_token, refresh_token, remote_addr, x_forwarded_for, x_real_ip, user_agent, created_at, expired_at) VALUES ('f53a2a95-c17f-11e9-bf92-0242ac120002', 1, 'access_token3', 'refresh_token', 'remote_addr', 'x_forwarded_for', 'x_real_ip', 'user_agent', '2019-08-31 20:23:00', '2019-08-31 23:00:00');
INSERT INTO token_histories (uid, user_id, access_token, refresh_token, remote_addr, x_forwarded_for, x_real_ip, user_agent, created_at, expired_at) VALUES ('0b0727a5-c180-11e9-bf92-0242ac120002', 1, 'access_token4', 'refresh_token', 'remote_addr', 'x_forwarded_for', 'x_real_ip', 'user_agent', '2019-08-31 20:23:00', '2019-08-31 23:00:00');
`
var sqlInsertAction = `INSERT INTO actions (action_name, action_desc, active) VALUES (:action_name, :action_desc, :active);`
var sqlInsertRole = `INSERT INTO roles (role_name, role_desc, active) VALUES (:role_name, :role_desc, :active);`
var sqlInsertUser = `INSERT INTO users (full_name, username, email, hashed, verified, active) VALUES (:full_name, :username, :email, :hashed, :verified, :active);`
var sqlInsertRoleAction = `INSERT INTO role_actions (role_id, action_id) VALUES (:role_id, :action_id);`
var sqlInsertUserRole = `INSERT INTO user_roles (role_id, user_id) VALUES (:role_id, :user_id);`
var sqlGetAction = `SELECT * FROM actions WHERE id = ?`
var sqlGetRole = `SELECT * FROM roles WHERE id = ?`
var sqlGetUser = `SELECT * FROM users WHERE id = ?`
var sqlGetRoleAction = `SELECT * FROM role_actions WHERE id = ?`
var sqlGetUserRole = `SELECT * FROM user_roles WHERE id = ?`
