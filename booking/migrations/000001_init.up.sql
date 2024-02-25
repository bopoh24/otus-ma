CREATE TABLE service (
    id serial primary key,
    parent_id int references service(id) on delete cascade,
    name TEXT NOT NULL
);

-- Services tree
INSERT INTO service (id, parent_id, name) VALUES (1, NULL, 'Красота и уход');
INSERT INTO service (id, parent_id, name) VALUES (2, NULL, 'Автосервис');
INSERT INTO service (id, parent_id, name) VALUES (3, NULL, 'Медицина');
INSERT INTO service (id, parent_id, name) VALUES (4, NULL, 'Ремонт электроники');
INSERT INTO service (id, parent_id, name) VALUES (5, NULL, 'Спорт');
INSERT INTO service (id, parent_id, name) VALUES (6, NULL, 'Развлечения');


INSERT INTO service (id, parent_id, name) VALUES (7, 1, 'Мужская стрижка/бритье');
INSERT INTO service (id, parent_id, name) VALUES (8, 1, 'Женская стрижка');
INSERT INTO service (id, parent_id, name) VALUES (9, 1, 'Маникюр');
INSERT INTO service (id, parent_id, name) VALUES (10, 1, 'Педикюр');

;
INSERT INTO service (id, parent_id, name) VALUES (11, 2, 'Шиномонтаж');
INSERT INTO service (id, parent_id, name) VALUES (12, 2, 'Диагностика');
INSERT INTO service (id, parent_id, name) VALUES (13, 2, 'Замена масла');

INSERT INTO service (id, parent_id, name) VALUES (14, 3, 'Терапевт');
INSERT INTO service (id, parent_id, name) VALUES (15, 3, 'Стоматолог');
INSERT INTO service (id, parent_id, name) VALUES (16, 3, 'Окулист');
INSERT INTO service (id, parent_id, name) VALUES (17, 3, 'Лор');
INSERT INTO service (id, parent_id, name) VALUES (18, 3, 'Хирург');
INSERT INTO service (id, parent_id, name) VALUES (19, 3, 'Педиатр');
INSERT INTO service (id, parent_id, name) VALUES (20, 3, 'Кардиолог');

INSERT INTO service (id, parent_id, name) VALUES (21, 4, 'Ремонт телефонов');
INSERT INTO service (id, parent_id, name) VALUES (22, 4, 'Ремонт ноутбуков/компьютеров');
INSERT INTO service (id, parent_id, name) VALUES (23, 4, 'Ремонт телевизоров');
INSERT INTO service (id, parent_id, name) VALUES (24, 4, 'Ремонт бытовой техники');

INSERT INTO service (id, parent_id, name) VALUES (25, 5, 'Фитнес');
INSERT INTO service (id, parent_id, name) VALUES (26, 5, 'Бассейн');
INSERT INTO service (id, parent_id, name) VALUES (27, 5, 'Теннис');
INSERT INTO service (id, parent_id, name) VALUES (28, 5, 'Баскетбол');
INSERT INTO service (id, parent_id, name) VALUES (29, 5, 'Футбол');
INSERT INTO service (id, parent_id, name) VALUES (30, 5, 'Волейбол');
INSERT INTO service (id, parent_id, name) VALUES (31, 5, 'Йога');
INSERT INTO service (id, parent_id, name) VALUES (32, 5, 'Пилатес');
INSERT INTO service (id, parent_id, name) VALUES (33, 5, 'Бокс');
INSERT INTO service (id, parent_id, name) VALUES (34, 5, 'Кроссфит');



CREATE TYPE offer_status AS ENUM ('open', 'reserved', 'paid', 'submitted', 'canceled_by_customer', 'canceled_by_company', 'completed');


CREATE TABLE offer (
    id serial primary key,
    service_id int references service(id) on delete cascade,
    customer varchar(255) not null default '',
    company_id int not null,
    company_name text not null,
    location point not null default point(0, 0),
    datetime TIMESTAMP NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    status offer_status not null default 'open',
    canceled_reason TEXT NOT NULL default '',
    created_by varchar(255) not null,
    updated_by varchar(255) not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

CREATE INDEX offer_customer_idx ON offer (customer);
CREATE INDEX offer_company_id_idx ON offer (company_id);
CREATE INDEX offer_datetime_idx ON offer (datetime);






