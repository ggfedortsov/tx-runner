create table person
(
    id            integer      not null
        constraint person_pk
            primary key,
    first_name    varchar(100) not null,
    last_name     varchar(100) not null,
    date_of_birth date         not null
);
