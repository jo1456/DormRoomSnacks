use heroku_5873df879639de6;

create table Persons(
	netID varchar(6),
    firstName varchar(10) not null,
    lastName varchar(10) not null,
    email varchar(50),
    phone int,
    student bool, # 1 if student, 0 if staff 
    primary key(netID)
);

create table DiningHalls(
	id varchar(10),
	name varchar(255) not null,
    address varchar(255) not null,
    phone int,
    #operatingHours??
    primary key(id)
);

create table Foods(
    name varchar(255),
    description varchar(255) not null default '',
    price decimal(6,2) not null,
    availability bool not null,
    #nutritionFacts
    primary key(name)
);

create table Menu(
	diningHallID varchar(10),
    foodID varchar(255),
    servingTime char(1), # b-breakfast, l-lunch, d-dinner
    primary key (diningHallID),
    foreign key (diningHallID) references DiningHalls(id),
    foreign key (foodID) references Foods(name)
);

create table Orders(
	personID varchar(6),
    diningHallID varchar(10),
    foodID varchar(255),
    primary key(personID, diningHallID, foodID),
    foreign key (personID) references Persons(netID),
    foreign key (diningHallID) references DiningHalls(id),
    foreign key (foodID) references Foods(name)
);
