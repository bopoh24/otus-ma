CREATE TABLE account (
    customer varchar(255) not null primary key,
    balance decimal(10,2) not null,
    created_at timestamp not null default current_timestamp
);

CREATE TYPE transaction_type AS ENUM ('top-up','payment','refund');


CREATE TABLE transaction (
    id serial not null primary key,
    type transaction_type not null,
    customer varchar(255) not null,
    offer_id int default null,
    amount decimal(10,2) not null,
    created_at timestamp not null default current_timestamp
);


CREATE INDEX transaction_customer_idx ON transaction (customer);
CREATE INDEX transaction_offer_id_idx ON transaction (offer_id);

