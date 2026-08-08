/*
Navicat MySQL Data Transfer

Source Server         : localhost
Source Server Version : 50740
Source Host           : localhost:3306
Source Database       : blog

Target Server Type    : MYSQL
Target Server Version : 50740
File Encoding         : 65001

Date: 2025-11-04 17:12:05
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for blog_article
-- ----------------------------
DROP TABLE IF EXISTS `blog_article`;
CREATE TABLE `blog_article` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `title` varchar(40) NOT NULL COMMENT '标题',
  `content` text NOT NULL COMMENT '内容',
  `content_format` varchar(16) NOT NULL DEFAULT 'html' COMMENT '内容格式',
  `pic` varchar(300) NOT NULL COMMENT '图片',
  `uid` int(11) NOT NULL COMMENT '用户ID',
  `fid` int(11) NOT NULL COMMENT '分类ID',
  `ctime` int(11) NOT NULL COMMENT '创建时间',
  `edittime` int(11) NOT NULL COMMENT '修改时间',
  `view` int(11) NOT NULL COMMENT '查看次数',
  `status` int(11) NOT NULL COMMENT '当前状态',
  `mp3` varchar(200) DEFAULT NULL,
  `istop` int(11) NOT NULL DEFAULT '0' COMMENT '是否置顶',
  `viewtumb` int(11) NOT NULL DEFAULT '0' COMMENT '是否显示封面 0显示 1不显示',
  `articlepassword` varchar(50) DEFAULT NULL COMMENT '文章查看密码',
  `video` varchar(500) DEFAULT NULL COMMENT '视频连接',
  `file` varchar(500) DEFAULT NULL COMMENT '附件地址',
  `type` int(11) NOT NULL COMMENT '文章样式',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章表';

-- ----------------------------
-- Records of blog_article
-- ----------------------------

-- ----------------------------
-- Table structure for blog_category
-- ----------------------------
DROP TABLE IF EXISTS `blog_category`;
CREATE TABLE `blog_category` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(30) NOT NULL COMMENT '分类名称',
  `fid` int(11) NOT NULL COMMENT '父级ID',
  `type` int(11) NOT NULL COMMENT '分类样式',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '分类排序',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=49 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='分类表';

-- ----------------------------
-- Records of blog_category
-- ----------------------------
INSERT INTO `blog_category` (`id`, `name`, `fid`, `type`, `sort`) VALUES
  (1, '默认分类', 0, 1, 0),
  (2, '技术笔记', 1, 1, 0);

-- ----------------------------
-- Table structure for blog_comment
-- ----------------------------
DROP TABLE IF EXISTS `blog_comment`;
CREATE TABLE `blog_comment` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(20) NOT NULL COMMENT '评论者姓名',
  `email` varchar(30) NOT NULL COMMENT '评论者邮箱',
  `content` varchar(200) NOT NULL COMMENT '评论者内容',
  `uid` int(11) NOT NULL DEFAULT '0' COMMENT '评论者ID',
  `replay` int(11) DEFAULT NULL COMMENT '评论谁',
  `ctime` int(11) NOT NULL COMMENT '评论时间',
  `aid` int(11) NOT NULL COMMENT '文章ID',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态 0显示 1不显示',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=131 DEFAULT CHARSET=utf8 COMMENT='评论表';

-- ----------------------------
-- Records of blog_comment
-- ----------------------------

-- ----------------------------
-- Table structure for blog_friendlink
-- ----------------------------
DROP TABLE IF EXISTS `blog_friendlink`;
CREATE TABLE `blog_friendlink` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '友链名称',
  `url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '友链地址',
  `logo` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '友链图标',
  `desc` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '友链描述',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态 0显示 1隐藏',
  `ctime` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_ctime` (`ctime`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='友情链接表';

-- ----------------------------
-- Records of blog_friendlink
-- ----------------------------
INSERT INTO `blog_friendlink` (`name`, `url`, `logo`, `desc`, `status`, `ctime`) VALUES
  ('GitHub', 'https://github.com', 'https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png', '代码托管与协作平台', 0, 0),
  ('Vue.js', 'https://vuejs.org', 'https://vuejs.org/images/logo.png', '渐进式 JavaScript 框架', 0, 0),
  ('Golang', 'https://go.dev', 'https://go.dev/images/go-logo-blue.svg', 'Go 编程语言官方网站', 0, 0);

-- ----------------------------
-- Table structure for blog_article_like
-- ----------------------------
DROP TABLE IF EXISTS `blog_article_like`;
CREATE TABLE `blog_article_like` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `aid` bigint(20) unsigned DEFAULT NULL,
  `visitor_id` varchar(64) DEFAULT NULL,
  `ctime` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_article_visitor` (`aid`, `visitor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章点赞表';

-- ----------------------------
-- Table structure for blog_statistic
-- ----------------------------
DROP TABLE IF EXISTS `blog_statistic`;
CREATE TABLE `blog_statistic` (
  `date` varchar(10) NOT NULL,
  `pv` bigint(20) DEFAULT NULL,
  `uv` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='访问统计表';

-- ----------------------------
-- Table structure for blog_visitor
-- ----------------------------
DROP TABLE IF EXISTS `blog_visitor`;
CREATE TABLE `blog_visitor` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `date` varchar(10) DEFAULT NULL,
  `visitor_id` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_date_visitor` (`date`, `visitor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='独立访客表';

-- ----------------------------
-- Table structure for blog_site
-- ----------------------------
DROP TABLE IF EXISTS `blog_site`;
CREATE TABLE `blog_site` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `title` varchar(200) NOT NULL COMMENT '网站标题',
  `keywords` text NOT NULL COMMENT '网站关键字',
  `description` text NOT NULL COMMENT '网站描述',
  `logo` varchar(200) NOT NULL COMMENT '网站LOGO',
  `articleSatus` int(11) NOT NULL COMMENT '0 无需审核 1 需要审核',
  `userStatus` int(11) NOT NULL COMMENT '0无需注册码 1需要注册码',
  `admin_email` varchar(100) NOT NULL COMMENT '管理员邮箱',
  `set_content` varchar(50) NOT NULL COMMENT '副标题',
  `name` varchar(50) NOT NULL COMMENT '网站名称',
  `statistics` text NOT NULL COMMENT '网站统计代码',
  `code` text NOT NULL COMMENT '邀请码说明',
  `friend_link` text NOT NULL COMMENT '友情链接说明',
  `icp` varchar(600) NOT NULL COMMENT 'ICP备案号',
  `submission` int(11) NOT NULL COMMENT '是否可以投稿 0可以 1不可以',
  `slides_display` int(11) NOT NULL COMMENT '是否显示幻灯片 0显示 1不显示',
  `file_size` int(11) NOT NULL COMMENT '文件大小限制',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='网站设置';

-- ----------------------------
-- Records of blog_site
-- ----------------------------

-- ----------------------------
-- Table structure for blog_user
-- ----------------------------
DROP TABLE IF EXISTS `blog_user`;
CREATE TABLE `blog_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `username` varchar(191) DEFAULT NULL COMMENT '用户登录名',
  `password` varchar(191) DEFAULT NULL COMMENT '用户登录密码',
  `pic` varchar(191) DEFAULT NULL COMMENT '用户头像',
  `email` varchar(191) DEFAULT NULL COMMENT '邮箱',
  `ctime` bigint(20) DEFAULT NULL COMMENT '创建时间',
  `lasttime` bigint(20) DEFAULT NULL COMMENT '最后一次登录时间',
  `ip` varchar(191) DEFAULT NULL COMMENT '登录ip',
  `status` bigint(20) DEFAULT NULL COMMENT '状态',
  `truename` varchar(191) DEFAULT NULL COMMENT '真实用户名',
  `bio` varchar(120) DEFAULT NULL COMMENT '个人简介',
  `admin` bigint(20) DEFAULT NULL COMMENT '是否管理员',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_blog_user_deleted_at` (`deleted_at`),
  KEY `idx_blog_user_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- ----------------------------
-- Records of blog_user
-- ----------------------------
INSERT INTO `blog_user` (`username`, `password`, `pic`, `email`, `ctime`, `lasttime`, `ip`, `status`, `truename`, `bio`, `admin`) VALUES
  ('admin', 'e10adc3949ba59abbe56e057f20f883e', '', '', 0, 0, '', 0, '管理员', 'Golang Dev', 1);

SET FOREIGN_KEY_CHECKS=1;
