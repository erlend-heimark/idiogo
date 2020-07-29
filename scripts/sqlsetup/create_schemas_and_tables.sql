create schema dadjokes
go

create table dadjokes.dadjokes
(
    id   varchar(20)
        constraint PK_DADJOKES
            primary key,
    joke varchar(2000) not null,
);
go
