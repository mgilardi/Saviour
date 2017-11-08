-- -----------------------------------------------------
-- Table `saviour`.'users'
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `saviour`.`users` (
  `uid` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NULL,
  `pass` VARCHAR(45) NULL,
  `mail` VARCHAR(45) NULL,
  `signature` VARCHAR(65) NULL,
  `created` DATE NULL,
  `access` VARCHAR(45) NULL,
  `login` VARCHAR(45) NULL,
  `status` VARCHAR(45) NULL,
  `timezone` VARCHAR(45) NULL,
  `language` VARCHAR(45) NULL,
  `picture` VARCHAR(45) NULL,
  PRIMARY KEY (`uid`),
  INDEX (`uid`))
  ENGINE = InnoDB;

-- -----------------------------------------------------
-- Table 'saviour.'sessions'
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `saviour`.`sessions` (
  `uid` INT NOT NULL AUTO_INCREMENT,
  `sid` VARCHAR(128) NOT NULL,
  `ssid` VARCHAR(128) NULL,
  `hostname` VARCHAR(128) NULL,
  `timestamp` INT(11) NULL,
  `cache` VARCHAR(45),
  `sesssion` VARCHAR(45) NULL,
  PRIMARY KEY (`uid`, `sid`),
  INDEX (`uid`))
  ENGINE = InnoDB;

-- -----------------------------------------------------
-- Table 'saviour.'users_roles'
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `saviour`.`user_roles` (
  `uid` INT NOT NULL,
  `rid` INT(11) NOT NULL,
  PRIMARY KEY (`uid`,`rid`),
  INDEX (`uid`, `rid`))
  ENGINE = InnoDB;

-- -----------------------------------------------------
-- Table 'saviour.'role'
-- -----------------------------------------------------

CREATE TABLE IF NOT EXISTS `saviour`.`role` (
  `rid` INT(11) NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64),
  `weight` INT(11),
  PRIMARY KEY (`rid`),
  INDEX (`rid`))
  ENGINE = InnoDB;

-- -----------------------------------------------------
-- Table 'saviour.'token'
-- -----------------------------------------------------

CREATE TABLE IF NOT EXISTS `saviour`.`login_token` (
  `uid` INT(11) NOT NULL,
  `tid` INT(11) NOT NULL AUTO_INCREMENT,
  `token` VARCHAR(40) NOT NULL UNIQUE,
  `created` DATETIME NOT NULL,
  PRIMARY KEY (`tid`),
  INDEX (`uid`))
  ENGINE = InnoDB;
