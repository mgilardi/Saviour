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
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `expires` bigint(11) DEFAULT NULL COMMENT 'Unix Time When Expires When NULL never expires'
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
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Token Created Timestamp',
  `expires` bigint(20) NOT NULL COMMENT 'Token expire time UNIXTIME'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `role`
--

CREATE TABLE `role` (
  `rid` int(11) NOT NULL COMMENT 'Role ID',
  `name` varchar(64) DEFAULT NULL COMMENT 'Role Name',
  `weight` int(11) DEFAULT NULL COMMENT 'Table Weight'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='User Roles are contained in this table';

INSERT INTO role (`rid`, `name`, `weight`) VALUES (1, 'admin', 1);
INSERT INTO role (`rid`, `name`, `weight`) VALUES (2, 'user', 2);

-- --------------------------------------------------------

--
-- Table structure for table `sessions`
--

CREATE TABLE `sessions` (
  `uid` int(11) NOT NULL COMMENT 'Associated UserID',
  `sid` int(11) NOT NULL COMMENT 'SessionID',
  `hostname` varchar(128) DEFAULT NULL COMMENT 'Current Session Hostname',
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created Timestamp',
  `expires` bigint(20) NOT NULL COMMENT 'Session Expire Time UNIXTIME',
  `sesssion` varchar(45) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `uid` int(11) NOT NULL COMMENT 'User ID',
  `name` varchar(45) NOT NULL COMMENT 'Username',
  `pass` varchar(255) NOT NULL COMMENT 'Hashed Password',
  `mail` varchar(45) DEFAULT NULL COMMENT 'Email Address',
  `created` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created Timestamp',
  `status` varchar(45) NOT NULL DEFAULT 'Offline' COMMENT 'Online/Offline Status',
  `timezone` varchar(45) DEFAULT NULL COMMENT 'User Selected Timezone'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`uid`, `name`, `pass`, `mail`, `created`, `status`, `timezone`) VALUES
(1, 'Admin', '$2a$14$GzWiSEdfzmjprH5oC6XXqeRIiN/LS3nggWmFSRVHi2eH8Vbgbxqbm', 'ian@diysecurity.com', '2017-11-09 17:12:23', 'Offline', 'Phoenix');

-- --------------------------------------------------------

--
-- Table structure for table `user_permissions`
--

CREATE TABLE `user_permissions` (
  `rid` int(11) NOT NULL COMMENT 'rid associated with the id in role',
  `module` varchar(255) NOT NULL COMMENT 'module name associated with loaded module',
  `allowed` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'Boolean for user access'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='User Permissions Table';

-- --------------------------------------------------------

--
-- Table structure for table `user_roles`
--

CREATE TABLE `user_roles` (
  `uid` int(11) NOT NULL,
  `rid` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='User_Roles Many-to-Many table';

INSERT INTO `user_roles` (`uid`, `rid`) VALUES (1, 1);

--
-- Indexes for dumped tables
--

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
-- Indexes for table `role`
--
ALTER TABLE `role`
  ADD PRIMARY KEY (`rid`);

--
-- Indexes for table `sessions`
--
ALTER TABLE `sessions`
  ADD PRIMARY KEY (`sid`) USING BTREE,
  ADD KEY `uid` (`uid`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`uid`),
  ADD UNIQUE KEY `name` (`name`);

--
-- Indexes for table `user_permissions`
--
ALTER TABLE `user_permissions`
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
ALTER TABLE `role`
  MODIFY `rid` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `sessions`
--
ALTER TABLE `sessions`
  MODIFY `sid` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `uid` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;
