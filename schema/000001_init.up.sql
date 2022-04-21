create table links
(
    id           serial
        constraint table_name_pk
            primary key,
    original_url varchar,
    short_url    varchar
);

create unique index table_name_id_uindex
    on links (id);

create unique index table_name_original_url_uindex
    on links (original_url);