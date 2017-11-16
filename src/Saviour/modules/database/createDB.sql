-- phpMyAdmin SQL Dump
-- version 4.6.6deb4
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Generation Time: Nov 15, 2017 at 01:06 AM
-- Server version: 5.7.20-0ubuntu0.17.04.1
-- PHP Version: 7.0.22-0ubuntu0.17.04.1

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";

--
-- Database: `saviour`
--
CREATE DATABASE IF NOT EXISTS `saviour` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
USE `saviour`;

-- --------------------------------------------------------

--
-- Table structure for table `cache`
--

DROP TABLE IF EXISTS `cache`;
CREATE TABLE `cache` (
  `cid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL COMMENT 'Primary Search Key For Cache',
  `data` longblob NOT NULL COMMENT 'Binary Data For Cache',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `expires` int(11) DEFAULT NULL COMMENT 'Unix Time When Expires When NULL never expires'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Cache Table takes in converted blob';

-- --------------------------------------------------------

--
-- Table structure for table `logger`
--

DROP TABLE IF EXISTS `logger`;
CREATE TABLE `logger` (
  `lid` int(11) NOT NULL COMMENT 'Log Primary Key',
  `type` varchar(6) NOT NULL COMMENT 'Log Type: Status, Warn, Error, Fatal',
  `module` varchar(20) NOT NULL COMMENT 'Module The Error Originated',
  `message` varchar(255) NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `login_token`
--

DROP TABLE IF EXISTS `login_token`;
CREATE TABLE `login_token` (
  `uid` int(11) NOT NULL,
  `tid` int(11) NOT NULL,
  `token` varchar(255) NOT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `role`
--

DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
  `rid` int(11) NOT NULL,
  `name` varchar(64) DEFAULT NULL,
  `weight` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `sessions`
--

DROP TABLE IF EXISTS `sessions`;
CREATE TABLE `sessions` (
  `uid` int(11) NOT NULL,
  `sid` int(11) NOT NULL,
  `ssid` varchar(128) DEFAULT NULL,
  `hostname` varchar(128) DEFAULT NULL,
  `timestamp` int(11) DEFAULT NULL,
  `cache` varchar(45) DEFAULT NULL,
  `sesssion` varchar(45) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `uid` int(11) NOT NULL,
  `name` varchar(45) DEFAULT NULL,
  `pass` varchar(45) DEFAULT NULL,
  `mail` varchar(45) DEFAULT NULL,
  `signature` varchar(65) DEFAULT NULL,
  `created` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `access` varchar(45) DEFAULT NULL,
  `login` varchar(45) DEFAULT NULL,
  `status` varchar(45) DEFAULT NULL,
  `timezone` varchar(45) DEFAULT NULL,
  `language` varchar(45) DEFAULT NULL,
  `picture` varchar(45) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`uid`, `name`, `pass`, `mail`, `signature`, `created`, `access`, `login`,
  `status`, `timezone`, `language`, `picture`)VALUES(1, 'Admin', 'Password', 'admin@diycardsecurity.com',
     NULL, '2017-11-10 00:12:23', NULL, NULL, NULL, NULL, NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `user_roles`
--

DROP TABLE IF EXISTS `user_roles`;
CREATE TABLE `user_roles` (
  `uid` int(11) NOT NULL,
  `rid` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `cache`
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
  ADD KEY `uid` (`uid`);

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
  ADD PRIMARY KEY (`uid`);

--
-- Indexes for table `user_roles`
--
ALTER TABLE `user_roles`
  ADD PRIMARY KEY (`rid`) USING BTREE,
  ADD KEY `uid` (`uid`) USING BTREE;

--
-- AUTO_INCREMENT for table `logger`
--
ALTER TABLE `logger`
  MODIFY `lid` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Log Primary Key';
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
