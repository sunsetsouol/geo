-- GEO (Generative Engine Optimization) Database Schema

CREATE DATABASE IF NOT EXISTS geo_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE geo_db;

-- 1. 提示词表 (Prompts)
CREATE TABLE `prompts` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `content` TEXT NOT NULL COMMENT '提示词内容',
    `category` VARCHAR(50) DEFAULT 'default' COMMENT '分类标签',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX `idx_category` (`category`)
) ENGINE=InnoDB COMMENT='提示词管理表';

-- 2. 任务表 (Tasks)
CREATE TABLE `tasks` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `prompt_id` BIGINT UNSIGNED NOT NULL COMMENT '关联的提示词ID',
    `status` ENUM('pending', 'processing', 'completed', 'failed') DEFAULT 'pending' COMMENT '任务状态',
    `last_run` DATETIME DEFAULT NULL COMMENT '最近运行时间',
    `retry_count` INT DEFAULT 0 COMMENT '失败重试次数',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `fk_task_prompt` FOREIGN KEY (`prompt_id`) REFERENCES `prompts` (`id`) ON DELETE CASCADE,
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB COMMENT='任务执行记录表';

-- 3. 结果与分析表 (Results)
CREATE TABLE `results` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `task_id` BIGINT UNSIGNED NOT NULL COMMENT '关联的任务ID',
    `response_text` LONGTEXT COMMENT '模型生成的回复内容 (含HTML原始数据)',
    `brand_score` DECIMAL(5, 2) DEFAULT 0.00 COMMENT '品牌曝光评分',
    `analysis_report` JSON COMMENT '详细分析报告 (存储为JSON对象)',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT `fk_result_task` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`) ON DELETE CASCADE,
    UNIQUE KEY `uk_task_id` (`task_id`)
) ENGINE=InnoDB COMMENT='任务执行结果与分析表';

-- 4. 引用链接表 (Citations)
CREATE TABLE `citations` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `task_id` BIGINT UNSIGNED NOT NULL COMMENT '关联的任务ID',
    `url` TEXT NOT NULL COMMENT '引用链接地址',
    `title` VARCHAR(255) DEFAULT NULL COMMENT '引用链接标题',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT `fk_citation_task` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`) ON DELETE CASCADE,
    INDEX `idx_task_id` (`task_id`)
) ENGINE=InnoDB COMMENT='大模型回复中的引用链接明细表';

-- 5. 优化文章表 (Articles)
CREATE TABLE `articles` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `brand_name` VARCHAR(100) NOT NULL COMMENT '目标品牌名称',
    `title` VARCHAR(255) NOT NULL COMMENT '文章标题',
    `content` LONGTEXT NOT NULL COMMENT '文章正文',
    `target_keywords` VARCHAR(255) DEFAULT NULL COMMENT '目标关键词',
    `publish_status` ENUM('pending', 'published') DEFAULT 'pending' COMMENT '发布状态',
    `status` varchar(20) DEFAULT 'pending' COMMENT '生成状态,pending:待生成,processing:生成中,completed:生成完成,failed:生成失败',
    `error` TEXT DEFAULT NULL COMMENT '生成错误信息',
    `published_url` TEXT DEFAULT NULL COMMENT '发布后的外链地址',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB COMMENT='生成的优化品牌曝光的文章表';
