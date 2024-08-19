create table `table_name` (
    `id` bigint not null auto_increment comment '主键 ID',
    `status` tinyint(3) unsigned not null default 0 comment '状态',
    `create_time` timestamp NOT NULL DEFAULT current_timestamp comment '创建时间',
    `update_time` timestamp NULL DEFAULT NULL on UPDATE current_timestamp COMMENT '更新时间',
    `create_by` varchar(32) NOT NULL DEFAULT '' comment '创建者',
    `update_by` varchar(32) NOT NULL DEFAULT '' comment '最近更新者',
    primary key (`id`),
    key `idx_ct` (`create_time`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='TableNameCN';


