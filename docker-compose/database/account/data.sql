CREATE DATABASE /*!32312 IF NOT EXISTS*/ `ecommerce_account` /*!40100 DEFAULT CHARACTER SET utf8mb3 */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `ecommerce_account`;


DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS accounts;

CREATE TABLE customers (
  id SERIAL,
  name varchar(100) NOT NULL,
  email varchar(100) NOT NULL,
  mobile_number varchar(20) NOT NULL,
  created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,

  PRIMARY KEY (id)
);

CREATE TABLE accounts (
  id SERIAL,
  account_number int NOT NULL,
  customer_id int NOT NULL,
  account_type varchar(100) NOT NULL,
  branch_address varchar(200) NOT NULL,
  created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  
   PRIMARY KEY (id)
  
);
CREATE INDEX customer_id_btree_index
ON accounts (customer_id);

INSERT INTO customers (name,email,mobile_number)
 VALUES ('Eazy Bytes','tutor@eazybytes.com','9876548337');

 
INSERT INTO accounts (customer_id, account_number, account_type, branch_address)
 VALUES (1, 186576453, 'Savings', '123 Main Street, New York');
 
