drop table if exists Article;
create table Article (
    ArticleID int auto_increment not null,
    UserID int not null,
    Title varchar(64) not null,
    Slug varchar(96) not null,
    Content text not null,
    
    primary key (`ArticleID`),
    foreign key (`UserID`) references User(`UserID`)
);
