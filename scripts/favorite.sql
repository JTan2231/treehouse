drop table if exists Favorite;
create table Favorite (
    FavoriteID int auto_increment not null,
    UserID int not null,
    ArticleID int not null,
    
    primary key (`FavoriteID`),
    foreign key (`UserID`) references User(`UserID`),
    foreign key (`ArticleID`) references Article(`ArticleID`)
);
