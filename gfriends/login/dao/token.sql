
CREATE TABLE `table_login_token` (
  `token` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT'',
  `uid` bigint(20) NOT NULL DEFAULT '0',
  `openid` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT'',
  `unionid` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT'',
  `login_type` int(11) NOT NULL DEFAULT '0',
  `login_scene` int(11) NOT NULL DEFAULT '0',
  `register` int(11) NOT NULL DEFAULT '0',
  `access_token` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT'',
  `refresh_token` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT'',
  `access_token_expire` bigint(20) NOT NULL DEFAULT '0',
  `refresh_token_expire` bigint(20) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `update_time` bigint(20) NOT NULL DEFAULT '0',
  `create_time` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`token`)
);

ALTER TABLE `table_login_token` add column  `expire_time` bigint(20) NOT NULL DEFAULT '0' after `status` 



