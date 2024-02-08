CREATE TABLE company (
    id serial primary key,
    owner varchar(255) not null,
    logo varchar(255) not null default '',
    name varchar(100) not null,
    description text not null default '',
    address varchar(255) not null default '',
    phone varchar(50) not null default '',
    location point not null default point(0, 0),
    active boolean not null default false,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);
