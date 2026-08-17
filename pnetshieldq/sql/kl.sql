create database netsheilq;
use netsheilq;
CREATE TABLE users (
                       user_id int PRIMARY KEY AUTO_INCREMENT,
                       username VARCHAR(255) NOT NULL ,
                       password VARCHAR(255) NOT NULL,
                       online_status BOOLEAN,
                       last_seen TIMESTAMP,
                       ip VARCHAR(255),
                       email VARCHAR(255),
                       client_id VARCHAR(255)  -- 新增字段存储MQTT客户端ID
);
ALTER TABLE users ADD UNIQUE (username);

CREATE TABLE message (
                         message_id int PRIMARY KEY AUTO_INCREMENT,
                         sender_id int,
                         receiver_id int,
                         content TEXT,
                         message_type VARCHAR(255),
                         filename TEXT,
                         timestamp TIMESTAMP
);
CREATE TABLE file (
                      file_id VARCHAR(255) PRIMARY KEY,
                      filename VARCHAR(255) NOT NULL,
                      sender_id VARCHAR(255),
                      receiver_id VARCHAR(255),
                      file_url TEXT,
                      timestamp TIMESTAMP,
                      filetype VARCHAR(255)
);

CREATE TABLE log (
                     log_id INT AUTO_INCREMENT PRIMARY KEY,
                     user_ip VARCHAR(255),
                     username VARCHAR(255),
                     message_id VARCHAR(255),
                     message_type VARCHAR(255),
                     action VARCHAR(255),
                     result VARCHAR(255),
                     timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `mqtt_user` (
                             `id` int(11) NOT NULL AUTO_INCREMENT,
                             `username` varchar(255) NOT NULL,
                             `password` varchar(255) NOT NULL,  -- 假设密 码是经过SHA256加密的
                             `is_superuser` tinyint(1) DEFAULT 0,  -- 0 表示非超级用户，1 表示超级用户
                             PRIMARY KEY (`id`),
                             UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `mqtt_acl` (
                            `id` int(11) NOT NULL AUTO_INCREMENT,
                            `allow` tinyint(1) NOT NULL,  -- 1 表示允许，0 表示拒绝
                            `ipaddr` varchar(60) DEFAULT NULL,  -- 客户端IP地址，可为 NULL
                            `username` varchar(255) DEFAULT NULL,  -- 用户名，可为 NULL
                            `clientid` varchar(255) DEFAULT NULL,  -- 客户端ID，可为 NULL
                            `access` int(11) NOT NULL,  -- 访问类型，通常1表示订阅，2表示发布，3表示订阅和发布
                            `topic` varchar(255) NOT NULL,  -- MQTT主题
                            PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
INSERT INTO `mqtt_user` (`username`, `password`, `is_superuser`)
VALUES
    ('user1', SHA2('user1password', 256), 0),
    ('admin', SHA2('adminpassword', 256), 1);

INSERT INTO `mqtt_acl` (`allow`, `ipaddr`, `username`, `clientid`, `access`, `topic`)
VALUES
    (1, NULL, 'user1', NULL, 3, 'user/1'),  -- 用户 user1 可以发布和订阅 'sensors/temperature'
    (1, NULL, 'admin', NULL, 3, '#'),                     -- 超级用户 admin 对所有主题有发布和订阅的权限
    (1, NULL, 'user1', 'device001', 2, 'user/2'),    -- 客户端ID为 device001 的设备可以发布到 'devices/status'
    (0, NULL, 'user1', NULL, 1, 'admin/secret');          -- 用户 user1 被禁止订阅 'admin/secret' 主题

CREATE TABLE friendships (
                             id INT AUTO_INCREMENT PRIMARY KEY,
                             user_id1 INT NOT NULL,
                             user_id2 INT NOT NULL,
                             status ENUM('pending', 'accepted') DEFAULT 'pending',
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                             FOREIGN KEY (user_id1) REFERENCES users(user_id),
                             FOREIGN KEY (user_id2) REFERENCES users(user_id)
);

INSERT INTO users (username, password, online_status, last_seen, ip, email, client_id)
VALUES
    ('user1', 'password1', true, NOW(), '192.168.1.1', 'user1@example.com', 'user_1');