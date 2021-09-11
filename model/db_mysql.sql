DROP DATABASE pwdmanage;
CREATE DATABASE pwdmanage;
USE pwdmanage;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id` int NOT NULL AUTO_INCREMENT,
    `username` varchar(255) CHARACTER SET utf8mb4 NULL DEFAULT NULL,
    `password` varchar(255) CHARACTER SET utf8mb4 NULL DEFAULT NULL,
    `pri` int NULL DEFAULT NULL,
    `role` varchar(255) CHARACTER SET utf8mb4 NULL DEFAULT NULL,
    `pass_strength` varchar(255) CHARACTER SET utf8mb4 NULL DEFAULT NULL,
    `mobile` varchar(255) CHARACTER SET utf8mb4 NULL DEFAULT NULL,
    `email` varchar(255) CHARACTER SET utf8mb4 NULL DEFAULT NULL,
    `remark` varchar(255) CHARACTER SET utf8mb4 NULL DEFAULT NULL,
    `token` varchar(255) CHARACTER SET utf8mb4 NULL DEFAULT NULL,
    `real_name` varchar(255) CHARACTER SET utf8mb4 NULL DEFAULT NULL,
    `department` varchar(255) CHARACTER SET utf8mb4 NULL DEFAULT NULL,
    `post` varchar(255) CHARACTER SET utf8mb4 NULL DEFAULT NULL,
    `address` varchar(255) CHARACTER SET utf8mb4 NULL DEFAULT NULL,
    `secret` varchar(255) CHARACTER SET utf8mb4 NULL DEFAULT NULL,
    `mfa_status` varchar(255) CHARACTER SET utf8mb4 NULL DEFAULT NULL,
    `created_at` datetime NULL DEFAULT NULL,
    `updated_at` datetime NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;
SET FOREIGN_KEY_CHECKS = 1;

DROP TABLE IF EXISTS `machine_user`;
CREATE TABLE `machine_user`
(
    `id`           int NOT NULL AUTO_INCREMENT,
    `machine_name` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
    `machine_pwd`  varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
    `remark`       varchar(1000) CHARACTER SET utf8mb4 DEFAULT NULL,
    `expired_at`   datetime(6) DEFAULT NULL,
    `created_tm`   datetime(6) DEFAULT NULL,
    `update_tm`    datetime(6) DEFAULT NULL,
    `dn`           varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
    `domain`       varchar(255) DEFAULT NULL,
    PRIMARY        KEY (`ID`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT=DYNAMIC;
SET FOREIGN_KEY_CHECKS = 1;

DROP TABLE IF EXISTS `domain`;
CREATE TABLE `domain`
(
    `id`          int NOT NULL AUTO_INCREMENT,
    `name`        varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
    `dc_hostname` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
    `status`      int         DEFAULT NULL,
    `err_msg`     varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
    `created_tm`  datetime(6) DEFAULT NULL,
    `dns`         varchar(255) DEFAULT NULL,
    `user_name`   varchar(255) DEFAULT NULL,
    `dn`          varchar(255) DEFAULT NULL,
    `user_dn`     varchar(255) DEFAULT NULL,
    `password`    varchar(255) DEFAULT NULL,
    PRIMARY       KEY (`ID`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT=DYNAMIC;
SET FOREIGN_KEY_CHECKS = 1;
