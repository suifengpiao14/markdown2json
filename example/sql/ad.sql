CREATE TABLE `ad`.`advertise` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `advertiser_id` varchar(32) NOT NULL COMMENT '广告主',
  `title` varchar(32) NOT NULL COMMENT '广告标题',
  `begin_at` datetime DEFAULT null COMMENT '投放开始时间',
  `end_at` datetime DEFAULT null  COMMENT '投放结束时间',
  `summary` varchar(128) DEFAULT '' COMMENT '广告素材-文字描述',
  `image` varchar(256) DEFAULT '' COMMENT '广告素材-图片地址',
  `link` varchar(512) DEFAULT NULL COMMENT '连接地址',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `type`  enum('text','image','vido') NOT NULL DEFAULT 'text' COMMENT '广告素材(类型),text-文字,image-图片,vido-视频',
  `value_obj` varchar(1024) NOT NULL DEFAULT '' COMMENT 'json扩展,广告的值属性对象',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `deleted_at` datetime  DEFAULT null COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8  COMMENT='广告表';



CREATE TABLE `ad`.`window` (
  `id`  int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `code` varchar(32) NOT NULL COMMENT '位置编码',
  `title` varchar(32) NOT NULL COMMENT '位置名称',
  `remark` varchar(255) NOT NULL COMMENT '位置描述',
  `content_types` varchar(50) DEFAULT 'text' COMMENT '广告素材(类型),text-文字,image-图片,vido-视频,多个逗号分隔',
  `width` smallint(6) DEFAULT 0 COMMENT '橱窗宽度',
  `high` smallint(6) DEFAULT 0 COMMENT '橱窗高度',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `deleted_at` datetime  DEFAULT null COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='橱窗表';


CREATE TABLE `ad`.`window_advertise` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `code` varchar(32) NOT NULL COMMENT '橱窗编码',
  `advertise_id`int(11) NOT NULL DEFAULT 0 COMMENT '广告ID',
  `advertise_priority` int(11) DEFAULT 0 COMMENT '广告优先级(同一个橱窗有多个广告时,按照优先级展示)',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `deleted_at` datetime  DEFAULT null COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code_ad_id` (`code`,`advertise_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='橱窗-广告表';