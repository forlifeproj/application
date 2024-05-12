CREATE TABLE `table_account` (
  `uid` bigint(20) NOT NULL AUTO_INCREMENT,
  `openid` varchar(255) COLLATE utf8_bin NOT NULL,
  `unionid` varchar(255) COLLATE utf8_bin NOT NULL,
  `open_type` int(11) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0',
  `update_time` bigint(20) NOT NULL DEFAULT '0',
  `create_time` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`)
);
ALTER TABLE `table_account` AUTO_INCREMENT = 1000000;


CREATE TABLE `table_openid` (
    `openid` varchar(255) COLLATE utf8_bin NOT NULL,
    `open_type` int(11) NOT NULL DEFAULT '0', 
    `uid` bigint(20) NOT NULL DEFAULT '0',
    `update_time` bigint(20) NOT NULL DEFAULT '0',
    `create_time` bigint(20) NOT NULL DEFAULT '0',
    PRIMARY KEY (`openid`,`open_type`)
);


CREATE TABLE `table_unionid` (
    `unionid` varchar(255) COLLATE utf8_bin NOT NULL,
    `open_type` int(11) NOT NULL DEFAULT '0', 
    `uid` bigint(20) NOT NULL DEFAULT '0',
    `update_time` bigint(20) NOT NULL DEFAULT '0',
    `create_time` bigint(20) NOT NULL DEFAULT '0',
    PRIMARY KEY (`unionid`,`open_type`)
);


