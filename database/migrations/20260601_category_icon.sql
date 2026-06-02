USE `billSoftware`;

ALTER TABLE `category_templates`
  ADD COLUMN `icon` VARCHAR(64) DEFAULT NULL AFTER `name`;

ALTER TABLE `user_categories`
  ADD COLUMN `icon` VARCHAR(64) DEFAULT NULL AFTER `name`;

-- 支出分类图标
UPDATE `category_templates` SET `icon` = 'food-bowl' WHERE `name` = '早午晚餐';
UPDATE `category_templates` SET `icon` = 'food-utensils' WHERE `name` = '餐饮';
UPDATE `category_templates` SET `icon` = 'shopping-bag' WHERE `name` = '购物';
UPDATE `category_templates` SET `icon` = 'daily-home' WHERE `name` = '日用';
UPDATE `category_templates` SET `icon` = 'drink-cup' WHERE `name` = '奶茶';
UPDATE `category_templates` SET `icon` = 'transport-car' WHERE `name` = '交通';
UPDATE `category_templates` SET `icon` = 'food-carrot' WHERE `name` = '蔬菜';
UPDATE `category_templates` SET `icon` = 'food-apple' WHERE `name` = '水果';
UPDATE `category_templates` SET `icon` = 'food-candy' WHERE `name` = '零食';
UPDATE `category_templates` SET `icon` = 'sport-run' WHERE `name` = '运动';
UPDATE `category_templates` SET `icon` = 'entertainment-game' WHERE `name` = '娱乐';
UPDATE `category_templates` SET `icon` = 'comm-phone' WHERE `name` = '通讯';
UPDATE `category_templates` SET `icon` = 'clothes-tshirt' WHERE `name` = '服饰';
UPDATE `category_templates` SET `icon` = 'beauty-mirror' WHERE `name` = '美容';
UPDATE `category_templates` SET `icon` = 'housing-key' WHERE `name` = '住房';
UPDATE `category_templates` SET `icon` = 'home-sofa' WHERE `name` = '居家';
UPDATE `category_templates` SET `icon` = 'family-child' WHERE `name` = '孩子';
UPDATE `category_templates` SET `icon` = 'family-elder' WHERE `name` = '长辈';
UPDATE `category_templates` SET `icon` = 'social-handshake' WHERE `name` = '社交';
UPDATE `category_templates` SET `icon` = 'travel-plane' WHERE `name` = '旅行';
UPDATE `category_templates` SET `icon` = 'drink-wine' WHERE `name` = '烟酒';
UPDATE `category_templates` SET `icon` = 'digital-phone' WHERE `name` = '数码';
UPDATE `category_templates` SET `icon` = 'transport-wheel' WHERE `name` = '汽车';
UPDATE `category_templates` SET `icon` = 'medical-cross' WHERE `name` = '医疗';
UPDATE `category_templates` SET `icon` = 'edu-book' WHERE `name` = '学习';
UPDATE `category_templates` SET `icon` = 'pet-paw' WHERE `name` = '宠物';
UPDATE `category_templates` SET `icon` = 'gift-box' WHERE `name` = '礼物';
UPDATE `category_templates` SET `icon` = 'work-briefcase' WHERE `name` = '办公';
UPDATE `category_templates` SET `icon` = 'repair-wrench' WHERE `name` = '维修';
UPDATE `category_templates` SET `icon` = 'donate-heart' WHERE `name` = '捐赠';
UPDATE `category_templates` SET `icon` = 'lottery-star' WHERE `name` = '彩票';
UPDATE `category_templates` SET `icon` = 'family-people' WHERE `name` = '亲友';
UPDATE `category_templates` SET `icon` = 'delivery-package' WHERE `name` = '快递';

-- 收入分类图标
UPDATE `category_templates` SET `icon` = 'income-salary' WHERE `name` = '工资';
UPDATE `category_templates` SET `icon` = 'income-parttime' WHERE `name` = '兼职';
UPDATE `category_templates` SET `icon` = 'income-invest' WHERE `name` = '理财';
UPDATE `category_templates` SET `icon` = 'income-gift' WHERE `name` = '礼金';
UPDATE `category_templates` SET `icon` = 'income-other' WHERE `name` = '其他';

-- 同步已有用户的分类图标
UPDATE `user_categories` uc
  JOIN `category_templates` ct ON uc.name = ct.name AND uc.category_type = ct.category_type
  SET uc.icon = ct.icon
  WHERE uc.icon IS NULL AND ct.icon IS NOT NULL;
