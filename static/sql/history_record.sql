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

 Date: 28/10/2019 15:13:25
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for history_record
-- ----------------------------
DROP TABLE IF EXISTS `history_record`;
CREATE TABLE `history_record` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_slug` varchar(50) NOT NULL,
  `is_cn` tinyint(1) NOT NULL COMMENT '国区',
  `ranking` int(11) NOT NULL COMMENT '排名',
  `solved_num` int(11) NOT NULL COMMENT '通过提交数量',
  `zero_time` int(11) NOT NULL COMMENT '零点时刻',
  `created_time` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `slug-zerotime` (`user_slug`,`zero_time`,`is_cn`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
