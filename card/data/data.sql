CREATE DATABASE /*!32312 IF NOT EXISTS*/ `ecommerce_card` /*!40100 DEFAULT CHARACTER SET utf8mb3 */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `ecommerce_card`;


DROP TABLE IF EXISTS cards;

CREATE TABLE `cards` (
  `id` int NOT NULL AUTO_INCREMENT,
  `card_number` varchar(100) NOT NULL,
  `customer_id` int NOT NULL,
  `card_type` varchar(100) NOT NULL,
  `total_limit` int NOT NULL,
  `amount_used` int NOT NULL,
  `available_amount` int NOT NULL,
   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);


INSERT INTO `cards` (`card_number`, `customer_id`, `card_type`, `total_limit`, `amount_used`, `available_amount`)
 VALUES ('4565XXXX4656', 1, 'Credit', 10000, 500, 9500);

INSERT INTO `cards` (`card_number`, `customer_id`, `card_type`, `total_limit`, `amount_used`, `available_amount`)
 VALUES ('3455XXXX8673', 1, 'Credit', 7500, 600, 6900);
 
INSERT INTO `cards` (`card_number`, `customer_id`, `card_type`, `total_limit`, `amount_used`, `available_amount`)
 VALUES ('2359XXXX9346', 1, 'Credit', 20000, 4000, 16000);
 
