drop table if exists Comment;
drop table if exists Subscribe;
drop table if exists Article;
drop table if exists User;

create table User (
    UserID int auto_increment not null,
    Username varchar(64) not null,
    Email varchar(64) not null,
    Password varchar(256) not null,

    primary key (`UserID`)
);

create table Article (
    ArticleID int auto_increment not null,
    UserID int not null,
    Title varchar(64) not null,
    Slug varchar(96) not null,
    Content text not null,

    primary key (`ArticleID`),
    foreign key (`UserID`) references User(`UserID`)
);

create table Subscribe (
    SubscriptionID int auto_increment not null,       
    SubscriberID int not null,
    SubscribeeID int not null,

    primary key (`SubscriptionID`),
    foreign key (`SubscriberID`) references User(`UserID`),
    foreign key (`SubscribeeID`) references User(`UserID`)
);

create table Comment (
    CommentID int auto_increment not null,
    ArticleID int not null,
    ParentID int null references Comment(`CommentID`) on delete cascade,
    UserID int not null,
    Content text not null,

    primary key (`CommentID`),
    foreign key (`ArticleID`) references Article(`ArticleID`),
    foreign key (`UserID`) references User(`UserID`)
);
