CREATE TABLE users_list
(
	id SERIAL NOT NULL PRIMARY KEY,
	firstName_lastNAme varchar(255) NOT NULL,
	chatID INT NOT NULL
);

CREATE TABLE product
(
	id SERIAL NOT NULL,
	user_id INT REFERENCES users_list(id) ON DELETE CASCADE,
	product_name varchar(255) NOT NULL UNIQUE,
	price FLOAT NOT NULL
)