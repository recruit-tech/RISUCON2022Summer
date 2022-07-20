DROP TABLE IF EXISTS `meeting_room`;

CREATE TABLE `meeting_room` (
    `id` VARCHAR(26) NOT NULL,
    `room_id` VARCHAR(32) NOT NULL,
    `start_at` DATETIME NOT NULL,
    `end_at` DATETIME NOT NULL,
    PRIMARY KEY(`id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4;
