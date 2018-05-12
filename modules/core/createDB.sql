-- phpMyAdmin SQL Dump
-- version 4.6.6deb4
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Generation Time: Nov 30, 2017 at 02:18 PM
-- Server version: 5.7.20-0ubuntu0.17.04.1
-- PHP Version: 7.0.22-0ubuntu0.17.04.1

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";

--
-- Database: `saviour`
--
CREATE DATABASE IF NOT EXISTS `saviour` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
USE `saviour`;

-- --------------------------------------------------------

--
-- Table structure for table `cache`
--

CREATE TABLE `cache` (
  `cid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL COMMENT 'Primary Search Key For Cache',
  `data` longblob NOT NULL COMMENT 'Binary Data For Cache',
  `created` bigint(11) NOT NULL COMMENT 'Created DateTime',
  `expires` bigint(11) NOT NULL COMMENT 'Unix Time When Expires When NULL never expires'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Cache Table takes in converted blob';

-- --------------------------------------------------------

--
-- Table structure for table `logger`
--

CREATE TABLE `logger` (
  `lid` int(11) NOT NULL COMMENT 'Log Primary Key',
  `type` varchar(6) NOT NULL COMMENT 'Log Type: Status, Warn, Error, Fatal',
  `module` varchar(20) NOT NULL COMMENT 'Module The Error Originated',
  `message` varchar(255) NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Logger Table Holds Saviours Error Logs';

-- --------------------------------------------------------

--
-- Table structure for table `login_token`
--

CREATE TABLE `login_token` (
  `uid` int(11) NOT NULL COMMENT 'User ID associated with token',
  `tid` int(11) NOT NULL COMMENT 'Token ID',
  `token` varchar(255) NOT NULL COMMENT 'Token',
  `created` bigint(11) NOT NULL COMMENT 'Token Created Timestamp',
  `expires` bigint(11) NOT NULL COMMENT 'Token expire time UNIXTIME'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `role`
--

CREATE TABLE `roles` (
  `rid` int(11) NOT NULL COMMENT 'Role ID',
  `name` varchar(64) DEFAULT NULL COMMENT 'Role Name',
  `weight` int(11) DEFAULT NULL COMMENT 'Table Weight'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='User Roles are contained in this table';

INSERT INTO roles (`rid`, `name`, `weight`) VALUES (1, 'administrator', 1);
INSERT INTO roles (`rid`, `name`, `weight`) VALUES (2, 'unauthorized.user', 2);
INSERT INTO roles (`rid`, `name`, `weight`) VALUES (3, 'authorized.user', 3);

-- --------------------------------------------------------
--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `uid` int(11) NOT NULL COMMENT 'User ID',
  `name` varchar(45) NOT NULL COMMENT 'Username',
  `pass` varchar(255) NOT NULL COMMENT 'Hashed Password',
  `mail` varchar(45) DEFAULT NULL COMMENT 'Email Address',
  `created` bigint(11) NOT NULL COMMENT 'User Created Timestamp',
  `status` varchar(45) NOT NULL DEFAULT 'Offline' COMMENT 'Online/Offline Status',
  `lastlogin` bigint(11) NULL COMMENT 'Last Login Timestamp',
  `timezone` varchar(45) DEFAULT NULL COMMENT 'User Selected Timezone'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`uid`, `name`, `pass`, `created`) VALUES
(0, 'Unauthorized', " ", UNIX_TIMESTAMP());
INSERT INTO `users` (`uid`, `name`, `pass`, `mail`, `created`, `status`, `timezone`) VALUES
(1, 'Admin', '$2a$14$GzWiSEdfzmjprH5oC6XXqeRIiN/LS3nggWmFSRVHi2eH8Vbgbxqbm', 'ian@diysecurity.com', UNIX_TIMESTAMP(), 'Offline', 'Phoenix');

-- --------------------------------------------------------

--
-- Table structure for table `user_permissions`
--

CREATE TABLE `role_permissions` (
  `rid` int(11) NOT NULL COMMENT 'rid associated with the id in role',
  `module` varchar(255) NOT NULL COMMENT 'module name associated with loaded module',
  `permission` varchar(255) NOT NULL COMMENT 'permission name for user access'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='User Permissions Table';

-- --------------------------------------------------------

--
-- Table structure for table `user_roles`
--

CREATE TABLE `user_roles` (
  `uid` int(11) NOT NULL,
  `rid` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='User_Roles Many-to-Many table';

INSERT INTO `user_roles` (`uid`, `rid`) VALUES (0, 2);
INSERT INTO `user_roles` (`uid`, `rid`) VALUES (1, 1);
INSERT INTO `user_roles` (`uid`, `rid`) VALUES (1, 3);

--
-- Indexes for table `cache`-
--
ALTER TABLE `cache`
  ADD PRIMARY KEY (`cid`);

--
-- Indexes for table `logger`
--
ALTER TABLE `logger`
  ADD PRIMARY KEY (`lid`),
  ADD KEY `type_index` (`type`),
  ADD KEY `module_index` (`module`);

--
-- Indexes for table `login_token`
--
ALTER TABLE `login_token`
  ADD PRIMARY KEY (`tid`),
  ADD UNIQUE KEY `uid` (`uid`) USING BTREE;

--
-- Indexes for table `roles`
--
ALTER TABLE `roles`
  ADD PRIMARY KEY (`rid`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`uid`),
  ADD UNIQUE KEY `name` (`name`);

--
-- Indexes for table `user_permissions`
--
ALTER TABLE `role_permissions`
  ADD KEY `rid` (`rid`,`module`);

--
-- Indexes for table `user_roles`
--
ALTER TABLE `user_roles`
  ADD PRIMARY KEY (`rid`) USING BTREE,
  ADD KEY `uid` (`uid`) USING BTREE;

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `logger`
--
ALTER TABLE `logger`
  MODIFY `lid` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `login_token`
--
ALTER TABLE `login_token`
  MODIFY `tid` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `role`
--
ALTER TABLE `roles`
  MODIFY `rid` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `uid` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;
