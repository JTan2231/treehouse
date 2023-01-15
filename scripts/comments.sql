drop table if exists Comment;
create table Comment (
    CommentID int auto_increment not null,
    ArticleID int not null,
    ParentID int null,
    UserID int not null,
    Content text not null,

    primary key (`CommentID`),
    foreign key (`ArticleID`) references Article(`ArticleID`),
    foreign key (`ParentID`) references Comment(`CommentID`),
    foreign key (`UserID`) references User(`UserID`)
);
