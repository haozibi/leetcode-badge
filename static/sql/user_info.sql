/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 50725
 Source Host           : localhost:3306
 Source Schema         : leetcode_badge

 Target Server Type    : MySQL
 Target Server Version : 50725
 File Encoding         : 65001

 Date: 28/10/2019 15:13:08
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user_info
-- ----------------------------
DROP TABLE IF EXISTS `user_info`;
CREATE TABLE `user_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_slug` varchar(50) NOT NULL COMMENT 'path 路径',
  `real_name` varchar(50) NOT NULL COMMENT '用户名称',
  `user_avatar` varchar(255) NOT NULL COMMENT '头像',
  `is_cn` tinyint(1) NOT NULL COMMENT '国区',
  `updated_time` int(11) NOT NULL COMMENT '更新时间',
  `created_time` int(11) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_slug` (`user_slug`,`is_cn`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
