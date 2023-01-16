drop table if exists Subscribe;
create table Subscribe (
    SubscriptionID int auto_increment not null,
    SubscriberID int not null,
    SubscribeeID int not null,

    primary key (`SubscriptionID`),
    foreign key (`SubscriberID`) references User(`UserID`),
    foreign key (`SubscribeeID`) references User(`UserID`)
);
