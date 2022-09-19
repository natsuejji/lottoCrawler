DROP TABLE IF EXISTS `lotto-info`;
CREATE TABLE `lotto-info`(
    `id` int NOT NULL AUTO_INCREMENT,
    `lotto-type` VARCHAR(10) NOT NULL,
    `draw-date` DATETIME NOT NULL,
    `lotto-period` VARCHAR(10)NOT NULL,
    `lotto-number-1` VARCHAR(2) ,
    `lotto-number-2` VARCHAR(2) ,
    `lotto-number-3` VARCHAR(2) ,
    `lotto-number-4` VARCHAR(2) ,
    PRIMARY KEY (`id`)
)