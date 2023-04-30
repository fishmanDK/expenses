CREATE TABLE users_list
(
	id SERIAL NOT NULL,
	firstName_lastNAme varchar(255) NOT NULL,
	chatID INT NOT NULL PRIMARY KEY
);


CREATE TABLE category
(
	id SERIAL NOT NULL PRIMARY KEY,
	user_id INT REFERENCES users_list(chatID) ON DELETE CASCADE,
	category_name varchar(255) NOT NULL 
);

CREATE TABLE product
(
	id SERIAL NOT NULL,
	category_id INT REFERENCES category(id) ON DELETE CASCADE,
	user_id INT REFERENCES users_list(chatID) ON DELETE CASCADE,
	product_name varchar(255) NOT NULL,
	price FLOAT NOT NULL,
	count INT
)
