drop table if exists User;
create table User (
    ID int auto_increment not null,
    Username varchar(64) not null,
    FirstName varchar(32) not null,
    LastName varchar(32) not null,
    Email varchar(64) not null,
    Password varchar(256) not null,
    primary key (`ID`)
);
