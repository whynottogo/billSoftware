USE `billSoftware`;

CREATE TABLE IF NOT EXISTS `file_uploads` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `biz_type` VARCHAR(32) NOT NULL COMMENT 'ledger_image|avatar|other',
  `original_name` VARCHAR(255) NOT NULL,
  `mime_type` VARCHAR(100) NOT NULL,
  `size_bytes` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `original_object_key` VARCHAR(500) NOT NULL,
  `thumbnail_object_key` VARCHAR(500) DEFAULT NULL,
  `width` INT DEFAULT NULL,
  `height` INT DEFAULT NULL,
  `status` VARCHAR(32) NOT NULL DEFAULT 'uploaded' COMMENT 'uploaded|bound|deleted',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_file_uploads_user_id` (`user_id`),
  KEY `idx_file_uploads_biz_type` (`biz_type`),
  KEY `idx_file_uploads_status` (`status`),
  CONSTRAINT `fk_file_uploads_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE `ledger_records`
  DROP COLUMN `image_url`,
  ADD COLUMN `image_file_id` BIGINT UNSIGNED DEFAULT NULL AFTER `record_date`,
  ADD KEY `idx_ledger_records_image_file_id` (`image_file_id`),
  ADD CONSTRAINT `fk_ledger_records_image_file_id`
    FOREIGN KEY (`image_file_id`) REFERENCES `file_uploads` (`id`) ON DELETE SET NULL;
