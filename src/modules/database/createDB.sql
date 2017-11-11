-- phpMyAdmin SQL Dump
-- version 4.6.6deb4
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Generation Time: Nov 09, 2017 at 05:08 PM
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
-- Table structure for table `login_token`
--

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

CREATE TABLE `role` (
  `rid` int(11) NOT NULL,
  `name` varchar(64) DEFAULT NULL,
  `weight` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `sessions`
--

CREATE TABLE `sessions` (
  `uid` int(11) NOT NULL,
  `sid` int(11) NOT NULL,
  `hostname` varchar(128) DEFAULT NULL,
  `timestamp` int(11) NULL DEFAULT CURRENT_TIMESTAMP,
  `cache` varchar(45) DEFAULT NULL,
  `sesssion` varchar(45) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `user_roles`
--

CREATE TABLE `user_roles` (
  `uid` int(11) NOT NULL,
  `rid` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `uid` int(11) NOT NULL,
  `name` varchar(45) NOT NULL UNIQUE,
  `pass` varchar(45) NOT NULL,
  `mail` varchar(45) NOT NULL UNIQUE,
  `signature` varchar(65) DEFAULT NULL,
  `created` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `access` varchar(45) DEFAULT NULL,
  `login` varchar(45) DEFAULT NULL,
  `status` varchar(45) DEFAULT NULL,
  `timezone` varchar(45) DEFAULT NULL,
  `language` varchar(45) DEFAULT NULL,
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Indexes for table `login_token`
--

ALTER TABLE `login_token`
  ADD PRIMARY KEY (`tid`) USING BTREE,
  ADD KEY `uid` (`uid`) USING BTREE;

--
-- Indexes for table `role`
--

ALTER TABLE `role`
  ADD PRIMARY KEY (`rid`) USING BTREE;

--
-- Indexes for table `sessions`
--

ALTER TABLE `sessions`
  ADD PRIMARY KEY (`sid`) USING BTREE,
  ADD KEY `uid` (`uid`) USING BTREE;

--
-- Indexes for table `user_roles`
--

ALTER TABLE `user_roles`
  ADD PRIMARY KEY (`rid`) USING BTREE,
  ADD KEY `uid` (`uid`) USING BTREE;

--
-- Indexes for table `users`
--

ALTER TABLE `users`
  ADD PRIMARY KEY (`uid`) USING BTREE;

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
  MODIFY `uid` int(11) NOT NULL AUTO_INCREMENT;