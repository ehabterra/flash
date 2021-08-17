file -inlinebatch END_OF_DROP_BATCH

DROP TABLE users                        IF EXISTS;
DROP TABLE banks                        IF EXISTS;
DROP TABLE user_bank_accounts           IF EXISTS;
DROP TABLE transactions                 IF EXISTS;
DROP TABLE transaction_accounts         IF EXISTS;
DROP TABLE transaction_recipients       IF EXISTS;

END_OF_DROP_BATCH

file -inlinebatch END_OF_BATCH

CREATE TABLE users (
    id          varchar(50) NOT NULL,
    username    varchar(250) NOT NULL,
    email       varchar(250) NOT NULL,
    password    varchar(250) NOT NULL,
    balance     DECIMAL default 0,
    CONSTRAINT PK_users PRIMARY KEY (id)
);

CREATE TABLE banks (
    id varchar(50) NOT NULL,
    name varchar(250) NOT NULL,
    CONSTRAINT PK_banks PRIMARY KEY (id)
);

CREATE TABLE user_bank_accounts (
    account_number varchar(50) NOT NULL,
    user_id varchar(50) NOT NULL,
    bank_id varchar(50) NOT NULL,
    branch_number varchar(20) NOT NULL,
    holder_name varchar(200) NOT NULL,
    reference varchar(500),
    CONSTRAINT PK_user_bank_accounts PRIMARY KEY (account_number)
);

CREATE TABLE transactions (
    id varchar(50) NOT NULL,
    user_id varchar(50) NOT NULL,
    type DECIMAL  NOT NULL,
    amount DECIMAL NOT NULL,
    date TIMESTAMP DEFAULT NOW,
    CONSTRAINT PK_transactions PRIMARY KEY (id)
);

CREATE TABLE transaction_accounts (
    transaction_id varchar(50) NOT NULL,
    account_number varchar(50) NOT NULL,
    CONSTRAINT PK_transaction_accounts PRIMARY KEY (transaction_id)
);

CREATE TABLE transaction_recipients (
    transaction_id varchar(50) NOT NULL,
    recipient_id varchar(50) NOT NULL,
    CONSTRAINT PK_transaction_recipients PRIMARY KEY (transaction_id)
);

END_OF_BATCH


INSERT INTO users(id,
username,
email   ,
password,
balance ) values ('user1', 'ehab1', 'ehab1@flash.com', '$2a$10$rDeErP2BTdnDnYoPjbHr8e4KUvkf8OT.KrXP2LXyDzBI9k8xugqYG', 0);

INSERT INTO users(id,
username,
email   ,
password,
balance ) values ('user2', 'ehab2', 'ehab2@flash.com', '$2a$10$rDeErP2BTdnDnYoPjbHr8e4KUvkf8OT.KrXP2LXyDzBI9k8xugqYG', 0);



INSERT INTO banks(id,
                  name ) values ('bank1', 'My Bank');

INSERT INTO banks(id,
                  name ) values ('bank2', 'My 2nd Bank');
