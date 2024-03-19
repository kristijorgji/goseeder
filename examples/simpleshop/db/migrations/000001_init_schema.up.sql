create schema if not exists goseeder_simpleshop collate utf8_general_ci;

create table categories
(
    id         smallint unsigned auto_increment
        primary key,
    name       json                                not null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    updated_at timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
);

create table products
(
    id         char(36)                            not null
        primary key,
    name       json                                not null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    updated_at timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
);
