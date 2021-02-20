-- -------
-- BMS
-- -------
CREATE TABLE book(
    id bigint(20) auto_increment PRIMARY KEY,
    title VARCHAR(20) not null,
    price double(16,2) not null #double对应go的float64
)engine=InnoDB default charset =utf8mb4;