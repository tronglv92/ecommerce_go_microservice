CREATE DATABASE /*!32312 IF NOT EXISTS*/ `ecommerce_loan` /*!40100 DEFAULT CHARACTER SET utf8mb3 */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `ecommerce_loan`;


DROP TABLE IF EXISTS loans;

CREATE TABLE `loans` (
  `id` int NOT NULL AUTO_INCREMENT,
  `customer_id` int NOT NULL,
  `start_date` timestamp NOT NULL,
  `loan_type` varchar(100) NOT NULL,
  `total_loan` int NOT NULL,
  `amount_paid` int NOT NULL,
  `outstanding_amount` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);

INSERT INTO `loans` ( `customer_id`, `start_date`, `loan_type`, `total_loan`, `amount_paid`, `outstanding_amount`)
 VALUES ( 1, NOW(), 'Home', 200000, 50000, 150000);
 
INSERT INTO `loans` ( `customer_id`, `start_date`, `loan_type`, `total_loan`, `amount_paid`, `outstanding_amount`)
 VALUES ( 1, NOW(), 'Vehicle', 40000, 10000, 30000);
 
INSERT INTO `loans` ( `customer_id`, `start_date`, `loan_type`, `total_loan`, `amount_paid`, `outstanding_amount`)
 VALUES ( 1, NOW(), 'Home', 50000, 10000, 40000);

INSERT INTO `loans` ( `customer_id`, `start_date`, `loan_type`, `total_loan`, `amount_paid`, `outstanding_amount`)
 VALUES ( 1, NOW(), 'Personal', 10000, 3500, 6500);