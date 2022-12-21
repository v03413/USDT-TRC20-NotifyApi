CREATE TABLE mm_trade (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	state tinyint DEFAULT (-1) NOT NULL, 
	amount varchar(16) not null,
	address varchar(34) not null,
	trade_no varchar(20) NOT NULL, 
	trade_hash varchar(64),
	notify_url varchar(256) not null,
	notify_retry tinyint default(0) NOT NULL,
	notify_time datetime, 
	expire_time datetime not null,
	create_time datetime NOT NULL, 
	update_time datetime
);

CREATE UNIQUE INDEX mm_trade_trade_no_IDX ON mm_trade (trade_no);