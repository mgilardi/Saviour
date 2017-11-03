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
  `Expire` INT(11) NULL,
  `Created` INT(11) NULL,
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


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
