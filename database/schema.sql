
-- schema for user

DROP TABLE IF EXISTS user_models;
CREATE TABLE `user_models`(
    `user_id` VARCHAR(100) PRIMARY KEY,
    `record_id` serial unique,
    `name` CHAR(100),
    `email` CHAR(100) UNIQUE NOT NULL,
    `password` VARCHAR(100),
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `creator` varchar(36),
    `updater` varchar(36),
    `active` tinyint(1) NOT NULL DEFAULT '1'
);

DROP TABLE IF EXISTS subscription_models;
CREATE TABLE `subscription_models`(
    `scheme_id` VARCHAR(100) PRIMARY KEY,
    `record_id` serial unique,
    `name` CHAR(100),
    `price` FLOAT,
    `days` INT,
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `creator` varchar(36),
    `updater` varchar(36),
    `active` tinyint(1) NOT NULL DEFAULT '1'
);

DROP TABLE IF EXISTS user_subscriptions;
CREATE TABLE `user_subscriptions`(
    `record_id` serial unique,
    `scheme_id` VARCHAR(100),
    `user_id` VARCHAR(100),
    `validity` DATETIME,
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `creator` varchar(36),
    `updater` varchar(36),
    `active` tinyint(1) NOT NULL DEFAULT '1'
);

ALTER TABLE user_subscriptions ADD PRIMARY KEY(`user_id`,`scheme_id`);
ALTER TABLE user_subscriptions ADD CONSTRAINT FOREIGN KEY(`scheme_id`) REFERENCES subscription_models(`scheme_id`);
ALTER TABLE user_subscriptions ADD CONSTRAINT FOREIGN KEY(`user_id`) REFERENCES user_models(`user_id`);

DROP TABLE IF EXISTS admin_models;
CREATE TABLE `admin_models`(
    `record_id` serial unique,
    `admin_id` VARCHAR(100) PRIMARY KEY,
    `email` CHAR(100) UNIQUE NOT NULL,
    `password` VARCHAR(100),
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `creator` varchar(36),
    `updater` varchar(36),
    `active` tinyint(1) NOT NULL DEFAULT '1'
);

DROP TABLE IF EXISTS news_letter_models;
CREATE TABLE `news_letter_models`(
    `record_id` serial unique,
    `news_letter_id` VARCHAR(100) PRIMARY KEY,
    `title` TEXT,
    `body` TEXT,
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `creator` varchar(36),
    `updater` varchar(36),
    `active` tinyint(1) NOT NULL DEFAULT '1'
);

DROP TABLE IF EXISTS news_schemes;
CREATE TABLE `news_schemes`(
    `record_id` serial unique,
    `news_letter_id` VARCHAR(100),
    `scheme_id` VARCHAR(100),
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `creator` varchar(36),
    `updater` varchar(36),
    `active` tinyint(1) NOT NULL DEFAULT '1'
);

ALTER TABLE news_schemes ADD PRIMARY KEY(`news_letter_id`,`scheme_id`);
ALTER TABLE news_schemes ADD CONSTRAINT FOREIGN KEY(`news_letter_id`) REFERENCES news_letter_models(`news_letter_id`);
ALTER TABLE news_schemes ADD CONSTRAINT FOREIGN KEY(`scheme_id`) REFERENCES subscription_models(`scheme_id`);

-- adding admin record manually

INSERT INTO admin_models (admin_id,email,password) VALUES('123','adminemail@some.com','$2a$10$XYJryrFLesYulUtNrZXLPOScNZqWci.uMsw5gMrFK0PYj89Tmjvee');