USE `billSoftware`;

CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(64) NOT NULL,
  `nickname` VARCHAR(64) NOT NULL,
  `phone` VARCHAR(20) NOT NULL,
  `email` VARCHAR(128) DEFAULT NULL,
  `password_hash` VARCHAR(255) NOT NULL,
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0=disabled,1=enabled',
  `avatar_original` LONGTEXT DEFAULT NULL,
  `avatar_compressed` LONGTEXT DEFAULT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_username` (`username`),
  UNIQUE KEY `uk_users_phone` (`phone`),
  KEY `idx_users_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `user_sessions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `session_token` VARCHAR(128) NOT NULL,
  `client_type` VARCHAR(32) NOT NULL DEFAULT 'pc',
  `is_active` TINYINT NOT NULL DEFAULT 1,
  `login_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_sessions_token` (`session_token`),
  KEY `idx_user_sessions_user_id` (`user_id`),
  KEY `idx_user_sessions_active` (`user_id`, `is_active`),
  CONSTRAINT `fk_user_sessions_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `category_templates` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `category_type` VARCHAR(16) NOT NULL COMMENT 'income|expense',
  `name` VARCHAR(64) NOT NULL,
  `sort_order` INT NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_category_templates_type` (`category_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `user_categories` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `category_type` VARCHAR(16) NOT NULL COMMENT 'income|expense',
  `name` VARCHAR(64) NOT NULL,
  `sort_order` INT NOT NULL DEFAULT 0,
  `is_system` TINYINT NOT NULL DEFAULT 1,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_categories_unique` (`user_id`, `category_type`, `name`),
  KEY `idx_user_categories_user_type` (`user_id`, `category_type`),
  CONSTRAINT `fk_user_categories_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `ledger_records` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `record_type` VARCHAR(16) NOT NULL COMMENT 'income|expense',
  `category_id` BIGINT UNSIGNED DEFAULT NULL,
  `amount` DECIMAL(18, 2) NOT NULL DEFAULT 0.00,
  `remark` VARCHAR(255) DEFAULT NULL,
  `record_date` DATE NOT NULL,
  `image_url` VARCHAR(500) DEFAULT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_ledger_records_user_date` (`user_id`, `record_date`),
  KEY `idx_ledger_records_user_type_date` (`user_id`, `record_type`, `record_date`),
  KEY `idx_ledger_records_category_id` (`category_id`),
  CONSTRAINT `fk_ledger_records_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_ledger_records_category_id` FOREIGN KEY (`category_id`) REFERENCES `user_categories` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `budgets` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `budget_type` VARCHAR(16) NOT NULL COMMENT 'month|year',
  `period_key` VARCHAR(16) NOT NULL COMMENT '2026-04 or 2026',
  `total_amount` DECIMAL(18, 2) NOT NULL DEFAULT 0.00,
  `used_amount` DECIMAL(18, 2) NOT NULL DEFAULT 0.00,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_budgets_unique` (`user_id`, `budget_type`, `period_key`),
  CONSTRAINT `fk_budgets_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `budget_items` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `budget_id` BIGINT UNSIGNED NOT NULL,
  `category_id` BIGINT UNSIGNED NOT NULL,
  `amount` DECIMAL(18, 2) NOT NULL DEFAULT 0.00,
  `used_amount` DECIMAL(18, 2) NOT NULL DEFAULT 0.00,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_budget_items_unique` (`budget_id`, `category_id`),
  CONSTRAINT `fk_budget_items_budget_id` FOREIGN KEY (`budget_id`) REFERENCES `budgets` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_budget_items_category_id` FOREIGN KEY (`category_id`) REFERENCES `user_categories` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `asset_accounts` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `account_type` VARCHAR(32) NOT NULL COMMENT 'cash|debit_card|credit_card|virtual|investment|liability|receivable|custom_asset',
  `sub_type` VARCHAR(32) DEFAULT NULL COMMENT 'wechat|alipay',
  `name` VARCHAR(64) NOT NULL,
  `remark` VARCHAR(255) DEFAULT NULL,
  `card_no` VARCHAR(64) DEFAULT NULL,
  `balance_nature` VARCHAR(16) NOT NULL COMMENT 'asset|liability',
  `balance` DECIMAL(18, 2) NOT NULL DEFAULT 0.00,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_asset_accounts_user_id` (`user_id`),
  CONSTRAINT `fk_asset_accounts_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `asset_account_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `account_id` BIGINT UNSIGNED NOT NULL,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `change_type` VARCHAR(16) NOT NULL COMMENT 'adjust|increase|decrease|edit',
  `amount_before` DECIMAL(18, 2) NOT NULL DEFAULT 0.00,
  `amount_change` DECIMAL(18, 2) NOT NULL DEFAULT 0.00,
  `amount_after` DECIMAL(18, 2) NOT NULL DEFAULT 0.00,
  `remark` VARCHAR(255) DEFAULT NULL,
  `log_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_asset_account_logs_account_id` (`account_id`),
  KEY `idx_asset_account_logs_user_id` (`user_id`),
  CONSTRAINT `fk_asset_account_logs_account_id` FOREIGN KEY (`account_id`) REFERENCES `asset_accounts` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_asset_account_logs_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `families` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `family_uid` VARCHAR(64) NOT NULL,
  `family_name` VARCHAR(64) NOT NULL,
  `creator_user_id` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_families_uid` (`family_uid`),
  CONSTRAINT `fk_families_creator_user_id` FOREIGN KEY (`creator_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `family_members` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `family_id` BIGINT UNSIGNED NOT NULL,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `joined_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_family_members_unique` (`family_id`, `user_id`),
  KEY `idx_family_members_user_id` (`user_id`),
  CONSTRAINT `fk_family_members_family_id` FOREIGN KEY (`family_id`) REFERENCES `families` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_family_members_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

