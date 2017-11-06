-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

-- -----------------------------------------------------
-- Schema saviour
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `saviour` DEFAULT CHARACTER SET utf8 ;
USE `saviour` ;

-- -----------------------------------------------------
-- Table `saviour`.`Users`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `saviour`.`Users` (
  `UID` INT NOT NULL AUTO_INCREMENT,
  `Username` VARCHAR(45) NULL,
  `Password` VARCHAR(45) NULL,
  `Admin` TINYINT(1) NULL,
  `Manager` TINYINT(1) NULL,
  `ServiceRep` TINYINT(1) NULL,
  `Created` DATE NULL,
  PRIMARY KEY (`UID`))
  ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `saviour`.`CreditCards`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `saviour`.`CreditCards` (
  `CCID` INT NOT NULL AUTO_INCREMENT,
  `Users_UID` INT NOT NULL,
  `Type` VARCHAR(45) NULL,
  `Number` INT NULL,
  `CCV` INT NULL,
  PRIMARY KEY (`CCID`, `Users_UID`),
  CONSTRAINT `fk_CreditCards_Users`
  FOREIGN KEY (`Users_UID`)
  REFERENCES `saviour`.`Users` (`UID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
  ENGINE = InnoDB;

CREATE INDEX `fk_CreditCards_Users_idx` ON `saviour`.`CreditCards` (`Users_UID` ASC);


-- -----------------------------------------------------
-- Table `saviour`.`ActivityLog`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `saviour`.`ActivityLog` (
  `AID` INT NOT NULL AUTO_INCREMENT,
  `Users_UID` INT NOT NULL,
  `Timestamp` DATETIME NULL,
  `Location` VARCHAR(45) NULL,
  `Log` VARCHAR(45) NULL,
  PRIMARY KEY (`AID`, `Users_UID`),
  CONSTRAINT `fk_ActivityLog_Users1`
  FOREIGN KEY (`Users_UID`)
  REFERENCES `saviour`.`Users` (`UID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
  ENGINE = InnoDB;

CREATE INDEX `fk_ActivityLog_Users1_idx` ON `saviour`.`ActivityLog` (`Users_UID` ASC);


-- -----------------------------------------------------
-- Table `saviour`.`Cache`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `saviour`.`Cache` (
  `CID` VARCHAR(255) NOT NULL,
  `Data` LONGBLOB NULL,
  `Created` BIGINT(64) NULL,
  `Expires` BIGINT(64) NULL,
  PRIMARY KEY (`CID`))
  ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `saviour`.`UserData`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `saviour`.`UserData` (
  `UDI` INT NOT NULL AUTO_INCREMENT,
  `Users_UID` INT NOT NULL,
  `Name` VARCHAR(60) NULL,
  `Phone` VARCHAR(50) NULL,
  `Email` VARCHAR(45) NULL,
  `Address` VARCHAR(45) NULL,
  PRIMARY KEY (`UDI`, `Users_UID`),
  CONSTRAINT `fk_UserData_Users1`
  FOREIGN KEY (`Users_UID`)
  REFERENCES `saviour`.`Users` (`UID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
  ENGINE = InnoDB;

CREATE INDEX `fk_UserData_Users1_idx` ON `saviour`.`UserData` (`Users_UID` ASC);


-- -----------------------------------------------------
-- Table `saviour`.`Rules`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `saviour`.`Rules` (
  `RID` INT NOT NULL AUTO_INCREMENT,
  `Users_UID` INT NOT NULL,
  `RName` VARCHAR(45) NULL,
  `RType` VARCHAR(45) NULL,
  `ROptions` VARCHAR(255) NULL,
  `RExtra` VARCHAR(255) NULL,
  PRIMARY KEY (`RID`, `Users_UID`),
  CONSTRAINT `fk_Rules_Users1`
  FOREIGN KEY (`Users_UID`)
  REFERENCES `saviour`.`Users` (`UID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
  ENGINE = InnoDB;

CREATE INDEX `fk_Rules_Users1_idx` ON `saviour`.`Rules` (`Users_UID` ASC);

USE `saviour` ;

-- -----------------------------------------------------
-- procedure CreateUser
-- -----------------------------------------------------

DELIMITER $$
USE `saviour`$$
CREATE PROCEDURE `CreateUser` (IN name VARCHAR(60), IN password VARCHAR(60))
  COMMENT 'Create a normal user'
  BEGIN
    DECLARE admin, manager, servicerep INT DEFAULT 0;
    INSERT INTO Users (Username, Password, Admin, Manager, ServiceRep, Created) VALUES (name, password, admin, manager, servicerep, CURDATE());
  END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure GetUID
-- -----------------------------------------------------

DELIMITER $$
USE `saviour`$$
CREATE PROCEDURE `GetUID` (IN user VARCHAR(60))
  COMMENT 'Get ID For User'
  BEGIN
    SELECT UID FROM Users Where Username = user;
  END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure SetUserData
-- -----------------------------------------------------

DELIMITER $$
USE `saviour`$$
CREATE PROCEDURE `SetUserData` (IN UID INT, IN name VARCHAR(60), IN phoneNum VARCHAR(50), IN email VARCHAR(45), IN address VARCHAR(45))
  COMMENT 'Sets UserData for each user'
BEGIN
INSERT INTO UserData (Users_UID, Name, Phone, Email, Address) VALUES (UID, name, phoneNum, email, address);
END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure GetUserData
-- -----------------------------------------------------

DELIMITER $$
USE `saviour`$$
CREATE PROCEDURE `GetUserData` (IN name VARCHAR(60))
  BEGIN
    SELECT UID, Username, Name, Phone, Email, Address, Created
    FROM Users JOIN UserData WHERE Users.UID = UserData.Users_UID
                                   AND Username = name;
  END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure GetPassword
-- -----------------------------------------------------

DELIMITER $$
USE `saviour`$$
CREATE PROCEDURE `GetPassword` (IN user VARCHAR(60))
  COMMENT 'Retrieve User Password'
  BEGIN
    SELECT Password FROM Users WHERE Username = user;
  END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure SetPassword
-- -----------------------------------------------------

DELIMITER $$
USE `saviour`$$
CREATE PROCEDURE `SetPassword` (IN user VARCHAR(45), IN pass VARCHAR(45))
  BEGIN
    UPDATE Users SET
      Password = pass
    WHERE Username = user;
  END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure MakeAdmin
-- -----------------------------------------------------

DELIMITER $$
USE `saviour`$$
CREATE PROCEDURE `MakeAdmin` (IN user VARCHAR(45))
  BEGIN
    UPDATE USERS SET
      Admin = 1
    Where Username = user;
  END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure WriteCache
-- -----------------------------------------------------

DELIMITER $$
USE `saviour`$$
CREATE PROCEDURE `WriteCache` (IN cachekey VARCHAR(255), IN datab BLOB, IN created DATETIME, IN expires DATETIME)
  BEGIN
    INSERT INTO Cache (CID, Data, Created, Expires)
    VALUES (cachekey, datab, created, expires);
  END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure ReadCache
-- -----------------------------------------------------

DELIMITER $$
USE `saviour`$$
CREATE PROCEDURE `ReadCache` (IN findkey VARCHAR(255))
  BEGIN
    SELECT Data FROM Cache WHERE CID = findkey;
  END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure RemoveCache
-- -----------------------------------------------------

DELIMITER $$
USE `saviour`$$
CREATE PROCEDURE `RemoveCache` (IN findkey VARCHAR(255))
  BEGIN
    DELETE FROM Cache WHERE CID = findkey;
  END$$

DELIMITER ;

CREATE USER 'saviour' IDENTIFIED BY 'Wlzv1iwFSV0W6cL4';

GRANT EXECUTE ON ROUTINE `saviour`.* TO 'saviour';

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
