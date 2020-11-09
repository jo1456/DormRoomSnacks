#create schema DormRoomSnacks;

use DormRoomSnacks;

#create schema DormRoomSnacks;

use DormRoomSnacks;

create table Persons(
	id int NOT NULL AUTO_INCREMENT,
	netID varchar(6),
    firstName varchar(10) not null,
    lastName varchar(10) not null,
    email varchar(50),
    phone int,
    student bool, # 1 if student, 0 if staff
	dollarBalance float,
	mealSwipeBalance int,
    primary key(id)
);

create table DiningHalls(
	id int NOT NULL AUTO_INCREMENT,
	name varchar(255) not null,
  address varchar(255) not null,
  phone varchar(255),
	menuID int,
  hours varchar(255),
    primary key(id)
);

create table Foods(
	id int NOT NULL AUTO_INCREMENT,
	menuID int,
    name varchar(255),
    description varchar(255) not null default '',
    price int not null,
    availability bool not null,
    nutritionFacts varchar(255),
    primary key(id)
);

create table Menu(
		id int NOT NULL AUTO_INCREMENT,
    #servingTime char(1), # b-breakfast, l-lunch, d-dinner
		name varchar(255),
        DiningHallId int,
    primary key (id),
    foreign key (id) references DiningHalls(id)
);

create table Orders(
	id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
	personID int,
	diningHallID int,
	status varchar(255),
	submitTime varchar(255),
	lastStatusChange varchar(255),
    foreign key (personID) references Persons(id),
    foreign key (diningHallID) references DiningHalls(id)
);

create table OrderItem(
	id int NOT NULL AUTO_INCREMENT,
	foodID int,
	orderID int,
	Customization varchar(255),
  primary key(id)
);
