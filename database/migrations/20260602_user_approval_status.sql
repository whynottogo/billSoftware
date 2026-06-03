USE `billSoftware`;

SET @approval_status_exists := (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'users'
    AND COLUMN_NAME = 'approval_status'
);
SET @approval_status_sql := IF(
  @approval_status_exists = 0,
  'ALTER TABLE `users` ADD COLUMN `approval_status` VARCHAR(16) NOT NULL DEFAULT ''pending'' COMMENT ''pending|approved|rejected'' AFTER `status`',
  'SELECT 1'
);
PREPARE approval_status_stmt FROM @approval_status_sql;
EXECUTE approval_status_stmt;
DEALLOCATE PREPARE approval_status_stmt;

SET @approval_updated_at_exists := (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'users'
    AND COLUMN_NAME = 'approval_updated_at'
);
SET @approval_updated_at_sql := IF(
  @approval_updated_at_exists = 0,
  'ALTER TABLE `users` ADD COLUMN `approval_updated_at` DATETIME DEFAULT NULL AFTER `approval_status`',
  'SELECT 1'
);
PREPARE approval_updated_at_stmt FROM @approval_updated_at_sql;
EXECUTE approval_updated_at_stmt;
DEALLOCATE PREPARE approval_updated_at_stmt;

SET @approval_remark_exists := (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'users'
    AND COLUMN_NAME = 'approval_remark'
);
SET @approval_remark_sql := IF(
  @approval_remark_exists = 0,
  'ALTER TABLE `users` ADD COLUMN `approval_remark` VARCHAR(255) DEFAULT NULL AFTER `approval_updated_at`',
  'SELECT 1'
);
PREPARE approval_remark_stmt FROM @approval_remark_sql;
EXECUTE approval_remark_stmt;
DEALLOCATE PREPARE approval_remark_stmt;

UPDATE `users`
SET `approval_status` = CASE
    WHEN `status` = 1 THEN 'approved'
    ELSE 'pending'
  END
WHERE `approval_status` IS NULL
  OR `approval_status` = ''
  OR `approval_status` NOT IN ('pending', 'approved', 'rejected');

UPDATE `users`
SET
  `approval_status` = 'approved',
  `approval_updated_at` = COALESCE(`approval_updated_at`, CURRENT_TIMESTAMP)
WHERE `status` = 1
  AND `approval_status` = 'pending';

SET @approval_status_index_exists := (
  SELECT COUNT(*)
  FROM information_schema.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'users'
    AND INDEX_NAME = 'idx_users_approval_status'
);
SET @approval_status_index_sql := IF(
  @approval_status_index_exists = 0,
  'ALTER TABLE `users` ADD INDEX `idx_users_approval_status` (`approval_status`)',
  'SELECT 1'
);
PREPARE approval_status_index_stmt FROM @approval_status_index_sql;
EXECUTE approval_status_index_stmt;
DEALLOCATE PREPARE approval_status_index_stmt;
