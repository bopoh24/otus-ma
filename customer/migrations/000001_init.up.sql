CREATE TABLE customer (
    id varchar(255) not null primary key,
    email varchar(100) not null default '' unique,
    first_name varchar(50) not null default '',
    last_name varchar(50) not null default '',
    photo varchar(255) not null default '',
    phone varchar(50) not null default '',
    location point not null default point(0, 0),
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);
