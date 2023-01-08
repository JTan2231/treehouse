drop table if exists User;
create table User (
    UserID int auto_increment not null,
    Username varchar(64) not null,
    Email varchar(64) not null,
    Password varchar(256) not null,

    primary key (`UserID`)
);

drop table if exists Article;
create table Article (
    ArticleID int auto_increment not null,
    UserID int not null,
    Title varchar(64) not null,
    Content text not null,

    primary key (`ArticleID`),
    foreign key (`UserID`) references User(`UserID`)
);
