drop table if exists User;
create table User (
    UserID int auto_increment not null,
    Username varchar(16) not null, 
    Email varchar(64) not null,    
    Password varchar(256) not null,
    
    primary key (`UserID`)
);
