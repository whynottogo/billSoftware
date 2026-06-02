USE `billSoftware`;

SET @avatar_original_exists := (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'users'
    AND COLUMN_NAME = 'avatar_original'
);
SET @avatar_original_sql := IF(
  @avatar_original_exists = 0,
  'ALTER TABLE `users` ADD COLUMN `avatar_original` LONGTEXT DEFAULT NULL AFTER `status`',
  'SELECT 1'
);
PREPARE avatar_original_stmt FROM @avatar_original_sql;
EXECUTE avatar_original_stmt;
DEALLOCATE PREPARE avatar_original_stmt;

SET @avatar_compressed_exists := (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'users'
    AND COLUMN_NAME = 'avatar_compressed'
);
SET @avatar_compressed_sql := IF(
  @avatar_compressed_exists = 0,
  'ALTER TABLE `users` ADD COLUMN `avatar_compressed` LONGTEXT DEFAULT NULL AFTER `avatar_original`',
  'SELECT 1'
);
PREPARE avatar_compressed_stmt FROM @avatar_compressed_sql;
EXECUTE avatar_compressed_stmt;
DEALLOCATE PREPARE avatar_compressed_stmt;
