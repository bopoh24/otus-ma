CREATE TABLE company (
    id serial primary key,
    logo varchar(255) not null default '',
    name varchar(100) not null,
    description text not null default '',
    address varchar(255) not null default '',
    phone varchar(50) not null default '',
    email varchar(100) not null default '',
    location point not null default point(0, 0),
    active boolean not null default false,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);


CREATE TYPE owner_role AS ENUM ('admin','manager');

CREATE TABLE company_manager (
    id serial primary key,
    company_id int not null references company(id) on delete cascade,
    user_id varchar(255) not null,
    email varchar(100) not null,
    first_name varchar(100) not null,
    last_name varchar(100) not null,
    role owner_role NOT NULL DEFAULT 'manager',
    active boolean not null default false,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

CREATE INDEX company_owners_user_id_index ON company_manager (user_id);
CREATE INDEX company_owners_email_index ON company_manager (email);


