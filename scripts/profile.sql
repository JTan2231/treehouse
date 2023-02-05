drop table if exists Profile;
create table Profile (
    ProfileID int auto_increment not null,
    UserID int not null,
    Bio varchar(500) null, 
    TwitterURL varchar(200) null,
    ProfilePicture varchar(1000) null,
    
    primary key (`ProfileID`),
    foreign key (`UserID`) references User(`UserID`)
);
