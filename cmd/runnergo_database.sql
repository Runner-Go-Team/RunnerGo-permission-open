
DROP TABLE IF EXISTS `auto_plan`;
CREATE TABLE `auto_plan` (
                             `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
                             `plan_id` varchar(100) NOT NULL COMMENT '计划ID',
                             `rank_id` bigint(10) NOT NULL DEFAULT '0' COMMENT '序号ID',
                             `team_id` varchar(100) NOT NULL COMMENT '团队ID',
                             `plan_name` varchar(255) NOT NULL COMMENT '计划名称',
                             `task_type` tinyint(2) NOT NULL DEFAULT '1' COMMENT '计划类型：1-普通任务，2-定时任务',
                             `status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '计划状态1:未开始,2:进行中',
                             `create_user_id` varchar(100) NOT NULL COMMENT '创建人id',
                             `run_user_id` varchar(100) NOT NULL COMMENT '运行人id',
                             `remark` text COMMENT '备注',
                             `run_count` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '运行次数',
                             `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                             `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                             `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                             PRIMARY KEY (`id`),
                             KEY `idx_plan_id` (`plan_id`),
                             KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='自动化测试-计划表';

-- ----------------------------
-- Table structure for auto_plan_email
-- ----------------------------
DROP TABLE IF EXISTS `auto_plan_email`;
CREATE TABLE `auto_plan_email` (
                                   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
                                   `plan_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '计划ID',
                                   `team_id` varchar(100) NOT NULL COMMENT '团队ID',
                                   `email` varchar(255) NOT NULL COMMENT '邮箱',
                                   `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                   `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                   `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='自动化测计划—收件人邮箱表';

-- ----------------------------
-- Table structure for auto_plan_report
-- ----------------------------
DROP TABLE IF EXISTS `auto_plan_report`;
CREATE TABLE `auto_plan_report` (
                                    `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                    `report_id` varchar(100) NOT NULL COMMENT '报告ID',
                                    `report_name` varchar(125) NOT NULL COMMENT '报告名称',
                                    `plan_id` varchar(100) NOT NULL COMMENT '计划ID',
                                    `rank_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '序号ID',
                                    `plan_name` varchar(255) NOT NULL COMMENT '计划名称',
                                    `team_id` varchar(100) NOT NULL COMMENT '团队ID',
                                    `task_type` int(11) NOT NULL DEFAULT '0' COMMENT '任务类型',
                                    `task_mode` int(11) NOT NULL DEFAULT '0' COMMENT '运行模式：1-按测试用例运行',
                                    `control_mode` tinyint(2) NOT NULL DEFAULT '0' COMMENT '控制模式：0-集中模式，1-单独模式',
                                    `scene_run_order` tinyint(2) NOT NULL DEFAULT '1' COMMENT '场景运行次序：1-顺序执行，2-同时执行',
                                    `test_case_run_order` tinyint(2) NOT NULL DEFAULT '1' COMMENT '测试用例运行次序：1-顺序执行，2-同时执行',
                                    `run_duration_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '任务运行持续时长',
                                    `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '报告状态1:进行中，2:已完成',
                                    `run_user_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '启动人id',
                                    `remark` text NOT NULL COMMENT '备注',
                                    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间（执行时间）',
                                    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                    PRIMARY KEY (`id`),
                                    KEY `idx_report_id` (`report_id`),
                                    KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='自动化测试计划-报告表';

-- ----------------------------
-- Table structure for auto_plan_task_conf
-- ----------------------------
DROP TABLE IF EXISTS `auto_plan_task_conf`;
CREATE TABLE `auto_plan_task_conf` (
                                       `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '配置ID',
                                       `plan_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '计划ID',
                                       `team_id` varchar(100) NOT NULL COMMENT '团队ID',
                                       `task_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '任务类型：1-普通模式，2-定时任务',
                                       `task_mode` tinyint(2) NOT NULL DEFAULT '1' COMMENT '运行模式：1-按照用例执行',
                                       `scene_run_order` tinyint(2) NOT NULL DEFAULT '1' COMMENT '场景运行次序：1-顺序执行，2-同时执行',
                                       `test_case_run_order` tinyint(2) NOT NULL DEFAULT '1' COMMENT '用例运行次序：1-顺序执行，2-同时执行',
                                       `run_user_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '运行人用户ID',
                                       `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                       `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                       `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                       PRIMARY KEY (`id`),
                                       KEY `idx_plan_id` (`plan_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='自动化测试—普通任务配置表';

-- ----------------------------
-- Table structure for auto_plan_timed_task_conf
-- ----------------------------
DROP TABLE IF EXISTS `auto_plan_timed_task_conf`;
CREATE TABLE `auto_plan_timed_task_conf` (
                                             `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '表id',
                                             `plan_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '计划id',
                                             `team_id` varchar(100) NOT NULL COMMENT '团队id',
                                             `frequency` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '任务执行频次: 0-一次，1-每天，2-每周，3-每月',
                                             `task_exec_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '任务执行时间',
                                             `task_close_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '任务结束时间',
                                             `task_type` tinyint(2) NOT NULL DEFAULT '2' COMMENT '任务类型：1-普通任务，2-定时任务',
                                             `task_mode` tinyint(2) NOT NULL DEFAULT '1' COMMENT '运行模式：1-按照用例执行',
                                             `scene_run_order` tinyint(2) NOT NULL DEFAULT '1' COMMENT '场景运行次序：1-顺序执行，2-同时执行',
                                             `test_case_run_order` tinyint(2) NOT NULL DEFAULT '1' COMMENT '测试用例运行次序：1-顺序执行，2-同时执行',
                                             `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '任务状态：0-未启用，1-运行中，2-已过期',
                                             `run_user_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '运行人用户ID',
                                             `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                             `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                             `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                             PRIMARY KEY (`id`),
                                             KEY `idx_plan_id` (`plan_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='自动化测试-定时任务配置表';

-- ----------------------------
-- Table structure for company
-- ----------------------------
DROP TABLE IF EXISTS `company`;
CREATE TABLE `company` (
                           `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                           `company_id` varchar(100) NOT NULL COMMENT '企业id',
                           `name` varchar(100) NOT NULL DEFAULT '' COMMENT '企业名称',
                           `logo` varchar(255) NOT NULL DEFAULT '' COMMENT '企业logo',
                           `expire_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '服务到期时间',
                           `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                           `deleted_at` datetime DEFAULT NULL,
                           PRIMARY KEY (`id`),
                           KEY `idx_company_id` (`company_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='企业表';

-- ----------------------------
-- Table structure for element
-- ----------------------------
DROP TABLE IF EXISTS `element`;
CREATE TABLE `element` (
                           `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
                           `element_id` varchar(100) NOT NULL COMMENT '全局唯一ID',
                           `element_type` varchar(10) NOT NULL COMMENT '类型：文件夹，元素',
                           `team_id` varchar(100) NOT NULL COMMENT '团队id',
                           `name` varchar(255) NOT NULL COMMENT '名称',
                           `parent_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '父级ID',
                           `locators` json DEFAULT NULL COMMENT '定位元素属性',
                           `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
                           `version` int(11) NOT NULL DEFAULT '0' COMMENT '产品版本号',
                           `created_user_id` varchar(100) NOT NULL COMMENT '创建人ID',
                           `description` text NOT NULL COMMENT '备注',
                           `source` tinyint(4) NOT NULL DEFAULT '0' COMMENT '数据来源：0-元素管理，1-场景管理',
                           `source_id` varchar(100) NOT NULL DEFAULT '' COMMENT '引用来源ID',
                           `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                           `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                           PRIMARY KEY (`id`),
                           KEY `idx_element_id` (`element_id`),
                           KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='元素表';

-- ----------------------------
-- Table structure for global_variable
-- ----------------------------
DROP TABLE IF EXISTS `global_variable`;
CREATE TABLE `global_variable` (
                                   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                                   `team_id` varchar(100) NOT NULL COMMENT '团队id',
                                   `type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '变量类型：1-全局变量，2-场景变量',
                                   `var` varchar(255) NOT NULL COMMENT '变量名',
                                   `val` text NOT NULL COMMENT '变量值',
                                   `description` text NOT NULL COMMENT '描述',
                                   `status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '开关状态：1-开启，2-关闭',
                                   `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                   `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                   `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                   PRIMARY KEY (`id`),
                                   KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='全局变量表';

-- ----------------------------
-- Table structure for machine
-- ----------------------------
DROP TABLE IF EXISTS `machine`;
CREATE TABLE `machine` (
                           `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                           `region` varchar(64) NOT NULL COMMENT '所属区域',
                           `ip` varchar(16) NOT NULL COMMENT '机器IP',
                           `port` int(11) unsigned NOT NULL COMMENT '端口',
                           `name` varchar(200) NOT NULL COMMENT '机器名称',
                           `cpu_usage` float unsigned NOT NULL DEFAULT '0' COMMENT 'CPU使用率',
                           `cpu_load_one` float unsigned NOT NULL DEFAULT '0' COMMENT 'CPU-1分钟内平均负载',
                           `cpu_load_five` float unsigned NOT NULL DEFAULT '0' COMMENT 'CPU-5分钟内平均负载',
                           `cpu_load_fifteen` float unsigned NOT NULL DEFAULT '0' COMMENT 'CPU-15分钟内平均负载',
                           `mem_usage` float unsigned NOT NULL DEFAULT '0' COMMENT '内存使用率',
                           `disk_usage` float unsigned NOT NULL DEFAULT '0' COMMENT '磁盘使用率',
                           `max_goroutines` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '最大协程数',
                           `current_goroutines` bigint(20) NOT NULL DEFAULT '0' COMMENT '已用协程数',
                           `server_type` tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '机器类型：1-主力机器，2-备用机器',
                           `status` tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '机器状态：1-使用中，2-已卸载',
                           `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                           `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                           PRIMARY KEY (`id`),
                           KEY `machine_region_ip_status_index` (`region`,`ip`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='压力测试机器表';

-- ----------------------------
-- Table structure for migrations
-- ----------------------------
DROP TABLE IF EXISTS `migrations`;
CREATE TABLE `migrations` (
                              `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                              `version` varchar(50) NOT NULL COMMENT '版本号',
                              `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                              `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                              `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                              PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for mock_target
-- ----------------------------
DROP TABLE IF EXISTS `mock_target`;
CREATE TABLE `mock_target` (
                               `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
                               `target_id` varchar(100) NOT NULL COMMENT '全局唯一ID',
                               `team_id` varchar(100) NOT NULL COMMENT '团队id',
                               `target_type` varchar(10) NOT NULL COMMENT '类型：文件夹，接口，分组，场景,测试用例',
                               `name` varchar(255) NOT NULL COMMENT '名称',
                               `parent_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '父级ID',
                               `method` varchar(16) NOT NULL COMMENT '方法',
                               `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
                               `type_sort` int(11) NOT NULL DEFAULT '0' COMMENT '类型排序',
                               `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '回收站状态：1-正常，2-回收站',
                               `version` int(11) NOT NULL DEFAULT '0' COMMENT '产品版本号',
                               `created_user_id` varchar(100) NOT NULL COMMENT '创建人ID',
                               `recent_user_id` varchar(100) NOT NULL COMMENT '最近修改人ID',
                               `description` text NOT NULL COMMENT '备注',
                               `source` tinyint(4) NOT NULL DEFAULT '0' COMMENT '数据来源：0-mock管理',
                               `plan_id` varchar(100) NOT NULL COMMENT '计划id',
                               `source_id` varchar(100) NOT NULL COMMENT '引用来源ID',
                               `is_checked` tinyint(2) NOT NULL DEFAULT '1' COMMENT '是否开启：1-开启，2-关闭',
                               `is_disabled` tinyint(2) NOT NULL DEFAULT '0' COMMENT '运行计划时是否禁用：0-不禁用，1-禁用',
                               `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                               `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                               `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                               PRIMARY KEY (`id`),
                               KEY `idx_target_id` (`target_id`),
                               KEY `idx_plan_id` (`plan_id`),
                               KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='创建目标';

-- ----------------------------
-- Table structure for mock_target_debug_log
-- ----------------------------
DROP TABLE IF EXISTS `mock_target_debug_log`;
CREATE TABLE `mock_target_debug_log` (
                                         `id` bigint(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                                         `target_id` varchar(100) NOT NULL COMMENT '目标唯一ID',
                                         `target_type` tinyint(2) NOT NULL COMMENT '目标类型：1-api，2-scene',
                                         `team_id` varchar(100) NOT NULL COMMENT '团队ID',
                                         `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                         `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                         `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                         PRIMARY KEY (`id`),
                                         KEY `idx_target_id` (`target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='mock目标调试日志表';

-- ----------------------------
-- Table structure for permission
-- ----------------------------
DROP TABLE IF EXISTS `permission`;
CREATE TABLE `permission` (
                              `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                              `permission_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '权限ID',
                              `title` varchar(100) NOT NULL DEFAULT '' COMMENT '权限内容',
                              `mark` varchar(100) NOT NULL DEFAULT '' COMMENT '权限标识',
                              `url` varchar(100) NOT NULL DEFAULT '' COMMENT '权限url',
                              `type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '类型（1：权限   2：功能）',
                              `group_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属权限组（1：企业成员管理  2：团队管理  3：角色管理）',
                              `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                              `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                              `deleted_at` datetime DEFAULT NULL,
                              PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';

-- ----------------------------
-- Table structure for preinstall_conf
-- ----------------------------
DROP TABLE IF EXISTS `preinstall_conf`;
CREATE TABLE `preinstall_conf` (
                                   `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                                   `conf_name` varchar(100) NOT NULL COMMENT '配置名称',
                                   `team_id` varchar(100) NOT NULL COMMENT '团队ID',
                                   `user_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '用户ID',
                                   `user_name` varchar(64) NOT NULL COMMENT '用户名称',
                                   `task_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '任务类型',
                                   `task_mode` tinyint(2) NOT NULL DEFAULT '0' COMMENT '压测模式',
                                   `control_mode` tinyint(2) NOT NULL DEFAULT '0' COMMENT '控制模式：0-集中模式，1-单独模式',
                                   `debug_mode` varchar(100) NOT NULL DEFAULT 'stop' COMMENT 'debug模式：stop-关闭，all-开启全部日志，only_success-开启仅成功日志，only_error-开启仅错误日志',
                                   `mode_conf` text NOT NULL COMMENT '压测配置详情',
                                   `timed_task_conf` text NOT NULL COMMENT '定时任务相关配置',
                                   `is_open_distributed` tinyint(2) NOT NULL DEFAULT '0' COMMENT '是否开启分布式调度：0-关闭，1-开启',
                                   `machine_dispatch_mode_conf` text NOT NULL COMMENT '分布式压力机配置',
                                   `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                   `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                   `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                   PRIMARY KEY (`id`),
                                   KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='预设配置表';

-- ----------------------------
-- Table structure for public_function
-- ----------------------------
DROP TABLE IF EXISTS `public_function`;
CREATE TABLE `public_function` (
                                   `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                                   `function` varchar(255) NOT NULL COMMENT '函数',
                                   `function_name` varchar(255) NOT NULL COMMENT '函数名称',
                                   `remark` text NOT NULL COMMENT '备注',
                                   `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                   `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                   `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='公共函数表';

-- ----------------------------
-- Table structure for report_machine
-- ----------------------------
DROP TABLE IF EXISTS `report_machine`;
CREATE TABLE `report_machine` (
                                  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                                  `report_id` varchar(100) NOT NULL COMMENT '报告id',
                                  `plan_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '计划ID',
                                  `team_id` varchar(100) NOT NULL COMMENT '团队ID',
                                  `ip` varchar(15) NOT NULL COMMENT '机器ip',
                                  `concurrency` bigint(20) NOT NULL DEFAULT '0' COMMENT '并发数',
                                  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                  PRIMARY KEY (`id`),
                                  KEY `idx_report_id` (`report_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
                        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                        `role_id` varchar(100) NOT NULL COMMENT '角色id',
                        `role_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '角色分类（1：企业  2：团队）',
                        `name` varchar(100) NOT NULL DEFAULT '' COMMENT '角色名称',
                        `company_id` varchar(100) NOT NULL DEFAULT '' COMMENT '企业id',
                        `level` tinyint(2) NOT NULL DEFAULT '0' COMMENT '角色层级（1:超管/团队管理员 2:管理员/团队成员 3:普通成员/只读成员/自定义角色） ',
                        `is_default` tinyint(2) NOT NULL DEFAULT '2' COMMENT '是否是默认角色  1：是   2：自定义角色',
                        `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        `deleted_at` datetime DEFAULT NULL,
                        PRIMARY KEY (`id`),
                        KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- ----------------------------
-- Table structure for role_permission
-- ----------------------------
DROP TABLE IF EXISTS `role_permission`;
CREATE TABLE `role_permission` (
                                   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                                   `role_id` varchar(100) NOT NULL COMMENT '角色id',
                                   `permission_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '权限id',
                                   `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                   `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                   `deleted_at` datetime DEFAULT NULL,
                                   PRIMARY KEY (`id`),
                                   KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限表';

-- ----------------------------
-- Table structure for role_type_permission
-- ----------------------------
DROP TABLE IF EXISTS `role_type_permission`;
CREATE TABLE `role_type_permission` (
                                        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                                        `role_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '角色分类（1：企业  2：团队）',
                                        `permission_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '权限id',
                                        `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                        `deleted_at` datetime DEFAULT NULL,
                                        PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色分类可拥有的权限';

-- ----------------------------
-- Table structure for scene_variable
-- ----------------------------
DROP TABLE IF EXISTS `scene_variable`;
CREATE TABLE `scene_variable` (
                                  `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                  `team_id` varchar(100) NOT NULL COMMENT '团队id',
                                  `scene_id` varchar(100) NOT NULL COMMENT '场景ID',
                                  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '使用范围：1-全局变量，2-场景变量',
                                  `var` varchar(255) NOT NULL COMMENT '变量名',
                                  `val` text NOT NULL COMMENT '变量值',
                                  `description` text NOT NULL COMMENT '描述',
                                  `status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '开关状态：1-开启，2-关闭',
                                  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                  PRIMARY KEY (`id`),
                                  KEY `idx_team_id` (`team_id`),
                                  KEY `idx_scene_id` (`scene_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='设置变量表';

-- ----------------------------
-- Table structure for setting
-- ----------------------------
DROP TABLE IF EXISTS `setting`;
CREATE TABLE `setting` (
                           `id` bigint(20) NOT NULL AUTO_INCREMENT,
                           `user_id` varchar(100) NOT NULL COMMENT '用户id',
                           `team_id` varchar(100) NOT NULL COMMENT '当前团队id',
                           `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                           `deleted_at` datetime DEFAULT NULL,
                           PRIMARY KEY (`id`),
                           KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='设置表';

-- ----------------------------
-- Table structure for sms_log
-- ----------------------------
DROP TABLE IF EXISTS `sms_log`;
CREATE TABLE `sms_log` (
                           `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
                           `type` tinyint(2) NOT NULL COMMENT '短信类型: 1-注册，2-登录，3-找回密码',
                           `mobile` char(11) NOT NULL DEFAULT '' COMMENT '手机号',
                           `content` varchar(200) NOT NULL COMMENT '短信内容',
                           `verify_code` varchar(20) NOT NULL COMMENT '验证码',
                           `verify_code_expiration_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '验证码有效时间',
                           `client_ip` varchar(100) NOT NULL DEFAULT '' COMMENT '客户端IP',
                           `send_status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '发送状态：1-成功 2-失败',
                           `verify_status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '校验状态：1-未校验 2-已校验',
                           `send_response` text NOT NULL COMMENT '短信服务响应',
                           `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                           `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                           PRIMARY KEY (`id`),
                           KEY `idx_type_mobile_verify_code` (`type`,`mobile`,`verify_code`,`verify_code_expiration_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='短信发送记录表';

-- ----------------------------
-- Table structure for stress_plan
-- ----------------------------
DROP TABLE IF EXISTS `stress_plan`;
CREATE TABLE `stress_plan` (
                               `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                               `plan_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '计划ID',
                               `team_id` varchar(100) NOT NULL COMMENT '团队ID',
                               `rank_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '序号ID',
                               `plan_name` varchar(255) NOT NULL COMMENT '计划名称',
                               `task_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '计划类型：1-普通任务，2-定时任务',
                               `task_mode` tinyint(2) NOT NULL DEFAULT '0' COMMENT '压测类型: 1-并发模式，2-阶梯模式，3-错误率模式，4-响应时间模式，5-每秒请求数模式，6-每秒事务数模式',
                               `status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '计划状态1:未开始,2:进行中',
                               `create_user_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '创建人id',
                               `run_user_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '运行人id',
                               `remark` text NOT NULL COMMENT '备注',
                               `run_count` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '运行次数',
                               `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                               `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                               `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                               PRIMARY KEY (`id`),
                               KEY `idx_plan_id` (`plan_id`),
                               KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='性能计划表';

-- ----------------------------
-- Table structure for stress_plan_email
-- ----------------------------
DROP TABLE IF EXISTS `stress_plan_email`;
CREATE TABLE `stress_plan_email` (
                                     `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
                                     `plan_id` varchar(100) NOT NULL COMMENT '计划ID',
                                     `team_id` varchar(100) NOT NULL COMMENT '团队ID',
                                     `email` varchar(255) DEFAULT NULL COMMENT '邮箱',
                                     `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                     `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                     `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                     PRIMARY KEY (`id`),
                                     KEY `idx_plan_id` (`plan_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='性能计划收件人';

-- ----------------------------
-- Table structure for stress_plan_report
-- ----------------------------
DROP TABLE IF EXISTS `stress_plan_report`;
CREATE TABLE `stress_plan_report` (
                                      `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                      `report_id` varchar(100) NOT NULL COMMENT '报告ID',
                                      `report_name` varchar(125) NOT NULL COMMENT '报告名称',
                                      `team_id` varchar(100) NOT NULL COMMENT '团队ID',
                                      `plan_id` varchar(100) NOT NULL COMMENT '计划ID',
                                      `rank_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '序号ID',
                                      `plan_name` varchar(255) NOT NULL COMMENT '计划名称',
                                      `scene_id` varchar(100) NOT NULL COMMENT '场景ID',
                                      `scene_name` varchar(255) NOT NULL COMMENT '场景名称',
                                      `task_type` int(11) NOT NULL COMMENT '任务类型',
                                      `task_mode` int(11) NOT NULL COMMENT '压测模式',
                                      `control_mode` tinyint(2) NOT NULL DEFAULT '0' COMMENT '控制模式：0-集中模式，1-单独模式',
                                      `debug_mode` varchar(100) NOT NULL DEFAULT 'stop' COMMENT 'debug模式：stop-关闭，all-开启全部日志，only_success-开启仅成功日志，only_error-开启仅错误日志',
                                      `run_duration_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '任务运行持续时长',
                                      `status` tinyint(4) NOT NULL COMMENT '报告状态1:进行中，2:已完成',
                                      `remark` text NOT NULL COMMENT '备注',
                                      `run_user_id` varchar(100) NOT NULL COMMENT '启动人id',
                                      `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间（执行时间）',
                                      `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                      `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                      PRIMARY KEY (`id`),
                                      KEY `idx_report_id` (`report_id`),
                                      KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='性能测试报告表';

-- ----------------------------
-- Table structure for stress_plan_task_conf
-- ----------------------------
DROP TABLE IF EXISTS `stress_plan_task_conf`;
CREATE TABLE `stress_plan_task_conf` (
                                         `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '配置ID',
                                         `plan_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '计划ID',
                                         `team_id` varchar(100) NOT NULL COMMENT '团队ID',
                                         `scene_id` varchar(100) NOT NULL COMMENT '场景ID',
                                         `task_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '任务类型：1-普通模式，2-定时任务',
                                         `task_mode` tinyint(2) NOT NULL DEFAULT '0' COMMENT '压测模式：1-并发模式，2-阶梯模式，3-错误率模式，4-响应时间模式，5-每秒请求数模式，6-每秒事务数模式',
                                         `control_mode` tinyint(2) NOT NULL DEFAULT '0' COMMENT '控制模式：0-集中模式，1-单独模式',
                                         `debug_mode` varchar(100) NOT NULL DEFAULT 'stop' COMMENT 'debug模式：stop-关闭，all-开启全部日志，only_success-开启仅成功日志，only_error-开启仅错误日志',
                                         `mode_conf` text NOT NULL COMMENT '压测模式配置详情',
                                         `is_open_distributed` tinyint(2) NOT NULL DEFAULT '0' COMMENT '是否开启分布式调度：0-关闭，1-开启',
                                         `machine_dispatch_mode_conf` text NOT NULL COMMENT '分布式压力机配置',
                                         `run_user_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '运行人用户ID',
                                         `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                         `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                         `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                         PRIMARY KEY (`id`),
                                         KEY `idx_plan_id` (`plan_id`),
                                         KEY `idx_scene_id` (`scene_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='性能计划—普通任务配置表';

-- ----------------------------
-- Table structure for stress_plan_timed_task_conf
-- ----------------------------
DROP TABLE IF EXISTS `stress_plan_timed_task_conf`;
CREATE TABLE `stress_plan_timed_task_conf` (
                                               `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '表id',
                                               `plan_id` varchar(100) NOT NULL COMMENT '计划id',
                                               `scene_id` varchar(100) NOT NULL COMMENT '场景id',
                                               `team_id` varchar(100) NOT NULL COMMENT '团队id',
                                               `user_id` varchar(100) NOT NULL COMMENT '用户ID',
                                               `frequency` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '任务执行频次: 0-一次，1-每天，2-每周，3-每月',
                                               `task_exec_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '任务执行时间',
                                               `task_close_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '任务结束时间',
                                               `task_type` tinyint(2) NOT NULL DEFAULT '2' COMMENT '任务类型：1-普通任务，2-定时任务',
                                               `task_mode` tinyint(2) NOT NULL DEFAULT '1' COMMENT '压测模式：1-并发模式，2-阶梯模式，3-错误率模式，4-响应时间模式，5-每秒请求数模式，6 -每秒事务数模式',
                                               `control_mode` tinyint(2) NOT NULL DEFAULT '0' COMMENT '控制模式：0-集中模式，1-单独模式',
                                               `debug_mode` varchar(100) NOT NULL DEFAULT 'stop' COMMENT 'debug模式：stop-关闭，all-开启全部日志，only_success-开启仅成功日志，only_error-开启仅错误日志',
                                               `mode_conf` text NOT NULL COMMENT '压测详细配置',
                                               `is_open_distributed` tinyint(2) NOT NULL DEFAULT '0' COMMENT '是否开启分布式调度：0-关闭，1-开启',
                                               `machine_dispatch_mode_conf` text NOT NULL COMMENT '分布式压力机配置',
                                               `run_user_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '运行人ID',
                                               `status` tinyint(11) NOT NULL DEFAULT '0' COMMENT '任务状态：0-未启用，1-运行中，2-已过期',
                                               `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                               `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                               `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                               PRIMARY KEY (`id`),
                                               KEY `idx_plan_id` (`plan_id`),
                                               KEY `idx_scene_id` (`scene_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='性能计划-定时任务配置表';

-- ----------------------------
-- Table structure for target
-- ----------------------------
DROP TABLE IF EXISTS `target`;
CREATE TABLE `target` (
                          `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
                          `target_id` varchar(100) NOT NULL COMMENT '全局唯一ID',
                          `team_id` varchar(100) NOT NULL COMMENT '团队id',
                          `target_type` varchar(10) NOT NULL COMMENT '类型：文件夹，接口，分组，场景,测试用例',
                          `name` varchar(255) NOT NULL COMMENT '名称',
                          `parent_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '父级ID',
                          `method` varchar(16) NOT NULL COMMENT '方法',
                          `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
                          `type_sort` int(11) NOT NULL DEFAULT '0' COMMENT '类型排序',
                          `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '回收站状态：1-正常，2-回收站',
                          `version` int(11) NOT NULL DEFAULT '0' COMMENT '产品版本号',
                          `created_user_id` varchar(100) NOT NULL COMMENT '创建人ID',
                          `recent_user_id` varchar(100) NOT NULL COMMENT '最近修改人ID',
                          `description` text NOT NULL COMMENT '备注',
                          `source` tinyint(4) NOT NULL DEFAULT '0' COMMENT '数据来源：0-测试对象，1-场景管理，2-性能，3-自动化测试， 4-mock',
                          `plan_id` varchar(100) NOT NULL COMMENT '计划id',
                          `source_id` varchar(100) NOT NULL COMMENT '引用来源ID',
                          `is_checked` tinyint(2) NOT NULL DEFAULT '1' COMMENT '是否开启：1-开启，2-关闭',
                          `is_disabled` tinyint(2) NOT NULL DEFAULT '0' COMMENT '运行计划时是否禁用：0-不禁用，1-禁用',
                          `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                          `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                          PRIMARY KEY (`id`),
                          KEY `idx_target_id` (`target_id`),
                          KEY `idx_plan_id` (`plan_id`),
                          KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='创建目标';

-- ----------------------------
-- Table structure for target_debug_log
-- ----------------------------
DROP TABLE IF EXISTS `target_debug_log`;
CREATE TABLE `target_debug_log` (
                                    `id` bigint(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                                    `target_id` varchar(100) NOT NULL COMMENT '目标唯一ID',
                                    `target_type` tinyint(2) NOT NULL COMMENT '目标类型：1-api，2-scene',
                                    `team_id` varchar(100) NOT NULL COMMENT '团队ID',
                                    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                    PRIMARY KEY (`id`),
                                    KEY `idx_target_id` (`target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='目标调试日志表';

-- ----------------------------
-- Table structure for team
-- ----------------------------
DROP TABLE IF EXISTS `team`;
CREATE TABLE `team` (
                        `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                        `team_id` varchar(100) NOT NULL COMMENT '团队ID',
                        `name` varchar(64) NOT NULL COMMENT '团队名称',
                        `description` text COMMENT '团队描述',
                        `company_id` varchar(100) NOT NULL DEFAULT '' COMMENT '所属企业id',
                        `type` tinyint(4) NOT NULL COMMENT '团队类型 1: 私有团队；2: 普通团队',
                        `trial_expiration_date` datetime NOT NULL COMMENT '试用有效期',
                        `is_vip` tinyint(2) NOT NULL DEFAULT '1' COMMENT '是否为付费团队 1-否 2-是',
                        `vip_expiration_date` datetime NOT NULL COMMENT '付费有效期',
                        `vum_num` bigint(20) NOT NULL DEFAULT '0' COMMENT '当前可用VUM总数',
                        `max_user_num` bigint(20) NOT NULL DEFAULT '0' COMMENT '当前团队最大成员数量',
                        `created_user_id` varchar(100) NOT NULL COMMENT '创建者id',
                        `team_buy_version_type` int(10) NOT NULL DEFAULT '1' COMMENT '团队套餐类型：1-个人版，2-团队版，3-企业版，4-私有化部署',
                        `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        `deleted_at` datetime DEFAULT NULL,
                        PRIMARY KEY (`id`),
                        KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='团队表';

-- ----------------------------
-- Table structure for team_env
-- ----------------------------
DROP TABLE IF EXISTS `team_env`;
CREATE TABLE `team_env` (
                            `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                            `team_id` varchar(100) NOT NULL COMMENT '团队id',
                            `name` varchar(100) NOT NULL COMMENT '环境名称',
                            `created_user_id` varchar(100) NOT NULL COMMENT '创建人id',
                            `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                            `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                            PRIMARY KEY (`id`),
                            KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='环境管理表';

-- ----------------------------
-- Table structure for team_env_database
-- ----------------------------
DROP TABLE IF EXISTS `team_env_database`;
CREATE TABLE `team_env_database` (
                                     `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                                     `team_id` varchar(100) NOT NULL COMMENT '团队id',
                                     `team_env_id` bigint(20) NOT NULL COMMENT '环境变量id',
                                     `type` varchar(100) NOT NULL COMMENT '数据库类型',
                                     `server_name` varchar(100) NOT NULL COMMENT 'mysql服务名称',
                                     `host` varchar(200) NOT NULL COMMENT '服务地址',
                                     `port` int(11) NOT NULL COMMENT '端口号',
                                     `user` varchar(100) NOT NULL COMMENT '账号',
                                     `password` varchar(200) NOT NULL COMMENT '密码',
                                     `db_name` varchar(100) NOT NULL COMMENT '数据库名称',
                                     `charset` varchar(100) NOT NULL DEFAULT 'utf8mb4' COMMENT '字符编码集',
                                     `created_user_id` varchar(100) NOT NULL COMMENT '创建人id',
                                     `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                     `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                     `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                     PRIMARY KEY (`id`),
                                     KEY `idx_team_id` (`team_id`),
                                     KEY `idx_team_env_id` (`team_env_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Sql数据库服务基础信息表';

-- ----------------------------
-- Table structure for team_env_service
-- ----------------------------
DROP TABLE IF EXISTS `team_env_service`;
CREATE TABLE `team_env_service` (
                                    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                                    `team_id` varchar(100) NOT NULL COMMENT '团队id',
                                    `team_env_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '环境id',
                                    `name` varchar(100) NOT NULL COMMENT '服务名称',
                                    `content` varchar(200) NOT NULL COMMENT '服务URL',
                                    `created_user_id` varchar(100) NOT NULL COMMENT '创建人id',
                                    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                    PRIMARY KEY (`id`),
                                    KEY `idxx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='团队环境服务管理';

-- ----------------------------
-- Table structure for team_user_queue
-- ----------------------------
DROP TABLE IF EXISTS `team_user_queue`;
CREATE TABLE `team_user_queue` (
                                   `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                   `email` varchar(255) NOT NULL COMMENT '邮箱',
                                   `team_id` varchar(100) NOT NULL COMMENT '团队id',
                                   `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                   `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                   `deleted_at` datetime DEFAULT NULL,
                                   PRIMARY KEY (`id`),
                                   KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='邀请待注册队列';

-- ----------------------------
-- Table structure for third_notice
-- ----------------------------
DROP TABLE IF EXISTS `third_notice`;
CREATE TABLE `third_notice` (
                                `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                                `notice_id` varchar(100) NOT NULL COMMENT '通知id',
                                `name` varchar(100) NOT NULL DEFAULT '' COMMENT '通知名称',
                                `channel_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '三方通知渠道id',
                                `params` json DEFAULT NULL COMMENT '通知参数',
                                `status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '1:启用 2:禁用',
                                `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                `deleted_at` datetime DEFAULT NULL,
                                PRIMARY KEY (`id`),
                                KEY `idx_notice_id` (`notice_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='三方通知设置';

-- ----------------------------
-- Table structure for third_notice_channel
-- ----------------------------
DROP TABLE IF EXISTS `third_notice_channel`;
CREATE TABLE `third_notice_channel` (
                                        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                                        `name` varchar(100) NOT NULL DEFAULT '' COMMENT '名称',
                                        `type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '类型 1:飞书  2:企业微信  3:邮箱  4:钉钉',
                                        `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                        `deleted_at` datetime DEFAULT NULL,
                                        PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='三方通知渠道';

-- ----------------------------
-- Table structure for third_notice_group
-- ----------------------------
DROP TABLE IF EXISTS `third_notice_group`;
CREATE TABLE `third_notice_group` (
                                      `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                                      `group_id` varchar(100) NOT NULL COMMENT '通知组id',
                                      `name` varchar(100) NOT NULL DEFAULT '' COMMENT '通知组名称',
                                      `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                      `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                      `deleted_at` datetime DEFAULT NULL,
                                      PRIMARY KEY (`id`),
                                      KEY `idx_group_id` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='三方通知组表';

-- ----------------------------
-- Table structure for third_notice_group_event
-- ----------------------------
DROP TABLE IF EXISTS `third_notice_group_event`;
CREATE TABLE `third_notice_group_event` (
                                            `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                                            `group_id` varchar(100) NOT NULL DEFAULT '' COMMENT '通知组id',
                                            `event_id` int(11) NOT NULL DEFAULT '0' COMMENT '事件id',
                                            `plan_id` varchar(100) NOT NULL DEFAULT '' COMMENT '计划ID',
                                            `team_id` varchar(100) NOT NULL DEFAULT '' COMMENT '团队ID',
                                            `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                            `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                            `deleted_at` datetime DEFAULT NULL,
                                            PRIMARY KEY (`id`),
                                            KEY `idx_group_id` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='三方通知组触发事件表';

-- ----------------------------
-- Table structure for third_notice_group_relate
-- ----------------------------
DROP TABLE IF EXISTS `third_notice_group_relate`;
CREATE TABLE `third_notice_group_relate` (
                                             `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                                             `group_id` varchar(100) NOT NULL COMMENT '通知组id',
                                             `notice_id` varchar(100) NOT NULL COMMENT '通知id',
                                             `params` json DEFAULT NULL COMMENT '通知目标参数',
                                             `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                             `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                             `deleted_at` datetime DEFAULT NULL,
                                             PRIMARY KEY (`id`),
                                             KEY `idx_notice_id` (`notice_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='三方通知组通知关联表';

-- ----------------------------
-- Table structure for ui_plan
-- ----------------------------
DROP TABLE IF EXISTS `ui_plan`;
CREATE TABLE `ui_plan` (
                           `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                           `plan_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '计划ID',
                           `team_id` varchar(100) NOT NULL COMMENT '团队ID',
                           `rank_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '序号ID',
                           `name` varchar(255) NOT NULL COMMENT '计划名称',
                           `task_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '计划类型：1-普通任务，2-定时任务',
                           `create_user_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '创建人id',
                           `head_user_id` varchar(1000) NOT NULL DEFAULT '0' COMMENT '负责人id ,用分割',
                           `run_count` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '运行次数',
                           `init_strategy` tinyint(2) NOT NULL DEFAULT '1' COMMENT '初始化策略：1-计划执行前重启浏览器，2-场景执行前重启浏览器，3-无初始化',
                           `description` text NOT NULL COMMENT '备注',
                           `browsers` json DEFAULT NULL COMMENT '浏览器信息',
                           `ui_machine_key` varchar(255) NOT NULL DEFAULT '' COMMENT '指定机器key',
                           `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                           `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                           PRIMARY KEY (`id`),
                           KEY `idx_plan_id` (`plan_id`),
                           KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='UI计划表';

-- ----------------------------
-- Table structure for ui_plan_report
-- ----------------------------
DROP TABLE IF EXISTS `ui_plan_report`;
CREATE TABLE `ui_plan_report` (
                                  `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                  `report_id` varchar(100) NOT NULL COMMENT '报告ID',
                                  `report_name` varchar(125) NOT NULL COMMENT '报告名称',
                                  `plan_id` varchar(100) NOT NULL COMMENT '计划ID',
                                  `plan_name` varchar(255) NOT NULL COMMENT '计划名称',
                                  `team_id` varchar(100) NOT NULL COMMENT '团队ID',
                                  `rank_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '序号ID',
                                  `task_type` int(11) NOT NULL DEFAULT '0' COMMENT '任务类型',
                                  `scene_run_order` tinyint(2) NOT NULL DEFAULT '1' COMMENT '场景运行次序：1-顺序执行，2-同时执行',
                                  `run_duration_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '任务运行持续时长',
                                  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '报告状态1:进行中，2:已完成',
                                  `run_user_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '启动人id',
                                  `remark` text NOT NULL COMMENT '备注',
                                  `browsers` json DEFAULT NULL COMMENT '浏览器信息',
                                  `ui_machine_key` varchar(255) NOT NULL DEFAULT '' COMMENT '指定机器key',
                                  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间（执行时间）',
                                  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                  PRIMARY KEY (`id`),
                                  KEY `idx_report_id` (`report_id`),
                                  KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='UI自动化测试计划-报告表';

-- ----------------------------
-- Table structure for ui_plan_task_conf
-- ----------------------------
DROP TABLE IF EXISTS `ui_plan_task_conf`;
CREATE TABLE `ui_plan_task_conf` (
                                     `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '配置ID',
                                     `plan_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '计划ID',
                                     `team_id` varchar(100) NOT NULL COMMENT '团队ID',
                                     `task_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '任务类型：1-普通模式，2-定时任务',
                                     `scene_run_order` tinyint(2) NOT NULL DEFAULT '1' COMMENT '场景运行次序：1-顺序执行，2-同时执行',
                                     `run_user_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '运行人用户ID',
                                     `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                     `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                     `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                     PRIMARY KEY (`id`),
                                     KEY `idx_plan_id` (`plan_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='UI自动化测试—普通任务配置表';

-- ----------------------------
-- Table structure for ui_plan_timed_task_conf
-- ----------------------------
DROP TABLE IF EXISTS `ui_plan_timed_task_conf`;
CREATE TABLE `ui_plan_timed_task_conf` (
                                           `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '表id',
                                           `plan_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '计划id',
                                           `team_id` varchar(100) NOT NULL COMMENT '团队id',
                                           `frequency` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '任务执行频次: 0-一次，1-每天，2-每周，3-每月，4-固定时间间隔',
                                           `task_exec_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '任务执行时间',
                                           `task_close_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '任务结束时间',
                                           `fixed_interval_start_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '固定时间间隔开始时间',
                                           `fixed_interval_time` int(10) NOT NULL DEFAULT '0' COMMENT '固定间隔时间',
                                           `fixed_run_num` int(10) NOT NULL DEFAULT '0' COMMENT '固定执行次数',
                                           `fixed_interval_time_type` int(10) NOT NULL DEFAULT '0' COMMENT '固定间隔时间类型：0-分钟，1-小时',
                                           `task_type` tinyint(2) NOT NULL DEFAULT '2' COMMENT '任务类型：1-普通任务，2-定时任务',
                                           `scene_run_order` tinyint(2) NOT NULL DEFAULT '1' COMMENT '场景运行次序：1-顺序执行，2-同时执行',
                                           `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '任务状态：0-未启用，1-运行中，2-已过期',
                                           `run_user_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '运行人用户ID',
                                           `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                           `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                           `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                           PRIMARY KEY (`id`),
                                           KEY `idx_plan_id` (`plan_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='UI自动化测试-定时任务配置表';

-- ----------------------------
-- Table structure for ui_scene
-- ----------------------------
DROP TABLE IF EXISTS `ui_scene`;
CREATE TABLE `ui_scene` (
                            `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
                            `scene_id` varchar(100) NOT NULL COMMENT '全局唯一ID',
                            `scene_type` varchar(10) NOT NULL COMMENT '类型：文件夹，场景',
                            `team_id` varchar(100) NOT NULL COMMENT '团队id',
                            `name` varchar(255) NOT NULL COMMENT '名称',
                            `parent_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '父级ID',
                            `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
                            `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '回收站状态：1-正常，2-回收站',
                            `version` int(11) NOT NULL DEFAULT '0' COMMENT '产品版本号',
                            `source` tinyint(2) NOT NULL DEFAULT '1' COMMENT '数据来源：1-场景管理，2-计划',
                            `plan_id` varchar(255) NOT NULL DEFAULT '' COMMENT '计划ID',
                            `created_user_id` varchar(100) NOT NULL COMMENT '创建人ID',
                            `recent_user_id` varchar(100) NOT NULL COMMENT '最近修改人ID',
                            `description` text NOT NULL COMMENT '备注',
                            `ui_machine_key` varchar(255) NOT NULL DEFAULT '' COMMENT '指定执行的UI自动化机器key',
                            `source_id` varchar(100) NOT NULL COMMENT '引用来源ID',
                            `browsers` json DEFAULT NULL COMMENT '浏览器信息',
                            `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                            `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                            PRIMARY KEY (`id`),
                            KEY `idx_scene_id` (`scene_id`),
                            KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='UI自动化场景';

-- ----------------------------
-- Table structure for ui_scene_element
-- ----------------------------
DROP TABLE IF EXISTS `ui_scene_element`;
CREATE TABLE `ui_scene_element` (
                                    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
                                    `scene_id` varchar(100) NOT NULL COMMENT '场景ID',
                                    `operator_id` varchar(100) NOT NULL COMMENT '操作ID',
                                    `element_id` varchar(100) NOT NULL COMMENT '元素ID',
                                    `team_id` varchar(100) NOT NULL COMMENT '团队id',
                                    `status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '状态 1：正常  2：回收站',
                                    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                    PRIMARY KEY (`id`),
                                    KEY `idx_scene_id` (`scene_id`),
                                    KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='UI自动化场景元素关联表';

-- ----------------------------
-- Table structure for ui_scene_operator
-- ----------------------------
DROP TABLE IF EXISTS `ui_scene_operator`;
CREATE TABLE `ui_scene_operator` (
                                     `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
                                     `operator_id` varchar(100) NOT NULL COMMENT '全局唯一ID',
                                     `scene_id` varchar(100) NOT NULL COMMENT '场景ID',
                                     `name` varchar(255) NOT NULL COMMENT '名称',
                                     `parent_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '父级ID',
                                     `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
                                     `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：1-正常，2-禁用',
                                     `type` varchar(100) NOT NULL DEFAULT '' COMMENT '步骤类型',
                                     `action` varchar(100) NOT NULL DEFAULT '' COMMENT '步骤方法',
                                     `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                     `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                     `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                     PRIMARY KEY (`id`),
                                     KEY `idx_scene_id` (`scene_id`),
                                     KEY `idx_operator_id` (`operator_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='UI自动化场景步骤';

-- ----------------------------
-- Table structure for ui_scene_sync
-- ----------------------------
DROP TABLE IF EXISTS `ui_scene_sync`;
CREATE TABLE `ui_scene_sync` (
                                 `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                                 `scene_id` varchar(100) NOT NULL COMMENT '场景ID',
                                 `source_scene_id` varchar(100) NOT NULL COMMENT '引用场景ID',
                                 `team_id` varchar(100) NOT NULL COMMENT '团队id',
                                 `sync_mode` tinyint(2) NOT NULL DEFAULT '0' COMMENT '状态：1-实时，2-手动,已场景为准   3-手动,已计划为准',
                                 `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                 `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                 `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                 PRIMARY KEY (`id`),
                                 KEY `idx_scene_id` (`scene_id`),
                                 KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='UI场景同步关系表';

-- ----------------------------
-- Table structure for ui_scene_trash
-- ----------------------------
DROP TABLE IF EXISTS `ui_scene_trash`;
CREATE TABLE `ui_scene_trash` (
                                  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
                                  `scene_id` varchar(100) NOT NULL COMMENT '场景ID',
                                  `team_id` varchar(100) NOT NULL COMMENT '团队id',
                                  `created_user_id` varchar(100) NOT NULL COMMENT '创建人ID',
                                  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                  PRIMARY KEY (`id`),
                                  KEY `idx_scene_id` (`scene_id`),
                                  KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='UI自动化场景回收站';

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
                        `id` bigint(20) NOT NULL AUTO_INCREMENT,
                        `user_id` varchar(100) NOT NULL COMMENT '用户id',
                        `account` varchar(100) NOT NULL DEFAULT '' COMMENT '账号',
                        `email` varchar(100) NOT NULL COMMENT '邮箱',
                        `mobile` char(11) NOT NULL COMMENT '手机号',
                        `password` varchar(255) NOT NULL COMMENT '密码',
                        `nickname` varchar(64) NOT NULL COMMENT '昵称',
                        `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
                        `wechat_open_id` varchar(100) NOT NULL COMMENT '微信开放的唯一id',
                        `utm_source` varchar(50) NOT NULL COMMENT '渠道来源',
                        `last_login_at` datetime DEFAULT NULL COMMENT '最近登录时间',
                        `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        `deleted_at` datetime DEFAULT NULL,
                        PRIMARY KEY (`id`),
                        KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- ----------------------------
-- Table structure for user_collect_info
-- ----------------------------
DROP TABLE IF EXISTS `user_collect_info`;
CREATE TABLE `user_collect_info` (
                                     `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                                     `user_id` varchar(100) NOT NULL COMMENT '用户id',
                                     `industry` varchar(100) NOT NULL COMMENT '所属行业',
                                     `team_size` varchar(20) NOT NULL COMMENT '团队规模',
                                     `work_type` varchar(20) NOT NULL COMMENT '工作岗位',
                                     `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                     `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                     `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                     PRIMARY KEY (`id`),
                                     KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for user_company
-- ----------------------------
DROP TABLE IF EXISTS `user_company`;
CREATE TABLE `user_company` (
                                `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                                `user_id` varchar(100) NOT NULL COMMENT '用户id',
                                `company_id` varchar(100) NOT NULL COMMENT '企业id',
                                `invite_user_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '邀请人id',
                                `invite_time` datetime DEFAULT NULL COMMENT '邀请时间',
                                `status` tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '状态：1-正常，2-已禁用',
                                `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                `deleted_at` datetime DEFAULT NULL,
                                PRIMARY KEY (`id`),
                                KEY `idx_company_id` (`company_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户企业关系表';


-- ----------------------------
-- Table structure for user_role
-- ----------------------------
DROP TABLE IF EXISTS `user_role`;
CREATE TABLE `user_role` (
                             `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                             `role_id` varchar(100) NOT NULL COMMENT '角色id',
                             `user_id` varchar(100) NOT NULL COMMENT '用户id',
                             `company_id` varchar(100) NOT NULL DEFAULT '' COMMENT '企业id',
                             `team_id` varchar(100) NOT NULL DEFAULT '' COMMENT '团队id',
                             `invite_user_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '邀请人id',
                             `invite_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '邀请时间',
                             `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                             `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                             `deleted_at` datetime DEFAULT NULL,
                             PRIMARY KEY (`id`),
                             KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表（企业角色、团队角色）';

-- ----------------------------
-- Table structure for user_team
-- ----------------------------
DROP TABLE IF EXISTS `user_team`;
CREATE TABLE `user_team` (
                             `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                             `user_id` varchar(100) NOT NULL COMMENT '用户ID',
                             `team_id` varchar(100) NOT NULL COMMENT '团队id',
                             `role_id` bigint(20) NOT NULL COMMENT '角色id1:超级管理员，2成员，3管理员',
                             `invite_user_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '邀请人id',
                             `invite_time` datetime DEFAULT NULL COMMENT '邀请时间',
                             `sort` int(11) NOT NULL DEFAULT '0',
                             `is_show` tinyint(2) NOT NULL DEFAULT '1' COMMENT '是否展示到团队列表  1:展示   2:不展示',
                             `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                             `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                             `deleted_at` datetime DEFAULT NULL,
                             PRIMARY KEY (`id`),
                             KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户团队关系表';

-- ----------------------------
-- Table structure for user_team_collection
-- ----------------------------
DROP TABLE IF EXISTS `user_team_collection`;
CREATE TABLE `user_team_collection` (
                                        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
                                        `user_id` varchar(100) NOT NULL COMMENT '用户ID',
                                        `team_id` varchar(100) NOT NULL COMMENT '团队id',
                                        `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                        `deleted_at` datetime DEFAULT NULL,
                                        PRIMARY KEY (`id`),
                                        KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户收藏团队表';

-- ----------------------------
-- Table structure for variable
-- ----------------------------
DROP TABLE IF EXISTS `variable`;
CREATE TABLE `variable` (
                            `id` bigint(20) NOT NULL AUTO_INCREMENT,
                            `team_id` varchar(100) NOT NULL COMMENT '团队id',
                            `scene_id` varchar(100) NOT NULL COMMENT '场景ID',
                            `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '使用范围：1-全局变量，2-场景变量',
                            `var` varchar(255) NOT NULL COMMENT '变量名',
                            `val` text NOT NULL COMMENT '变量值',
                            `description` text NOT NULL COMMENT '描述',
                            `status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '开关状态：1-开启，2-关闭',
                            `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                            `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                            PRIMARY KEY (`id`),
                            KEY `idx_team_id` (`team_id`),
                            KEY `idx_scene_id` (`scene_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='设置变量表';

-- ----------------------------
-- Table structure for variable_import
-- ----------------------------
DROP TABLE IF EXISTS `variable_import`;
CREATE TABLE `variable_import` (
                                   `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                   `team_id` varchar(100) NOT NULL COMMENT '团队id',
                                   `scene_id` varchar(100) NOT NULL DEFAULT '0' COMMENT '场景id',
                                   `name` varchar(128) NOT NULL COMMENT '文件名称',
                                   `url` varchar(255) NOT NULL COMMENT '文件地址',
                                   `uploader_id` varchar(100) NOT NULL COMMENT '上传人id',
                                   `status` tinyint(2) NOT NULL DEFAULT '1' COMMENT '开关状态：1-开，2-关',
                                   `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                   `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                   `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                   PRIMARY KEY (`id`),
                                   KEY `idx_team_id` (`team_id`),
                                   KEY `idx_scene_id` (`scene_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='导入变量表';
