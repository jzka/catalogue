CREATE USER IF NOT EXISTS 'catalogue_user' IDENTIFIED BY 'default_password';

GRANT ALL ON cataloguedb.* TO 'catalogue_user';

CREATE TABLE IF NOT EXISTS product (
	product_id varchar(40) NOT NULL, 
	name varchar(20), 
	description varchar(200), 
	price float, 
	stock int, 
	image_url varchar(40), 
	PRIMARY KEY(product_id)
);

CREATE TABLE IF NOT EXISTS type (
	type_id MEDIUMINT NOT NULL AUTO_INCREMENT, 
	name varchar(20), 
	PRIMARY KEY(type_id)
);

CREATE TABLE IF NOT EXISTS product_type (
	product_id varchar(40), 
	type_id MEDIUMINT NOT NULL, 
	FOREIGN KEY (product_id) 
		REFERENCES product(product_id), 
	FOREIGN KEY(type_id)
		REFERENCES type(type_id)
);

INSERT INTO product(product_id, name, description, price, stock, image_url) VALUES ("a0a4f044-b040-410d-8ead-4de0446aec7e", "Nerd leg", "For all those leg lovers out there. A perfect example of a swivel chair trained calf. Meticulously trained on a diet of sitting and Pina Coladas. Phwarr...", 7.99, 115, "asdf");
INSERT INTO product VALUES ("808a2de1-1aaa-4c25-a9b9-6612e8f29a38", "Crossed", "A mature sock, crossed, with an air of nonchalance.",  17.32, 738, "/catalogue/images/cross_1.jpeg");
INSERT INTO product VALUES ("510a0d7e-8e83-4193-b483-e27e09ddc34d", "SuperSport XL", "Ready for action. Engineers: be ready to smash that next bug! Be ready, with these super-action-sport-masterpieces. This particular engineer was chased away from the office with a stick.",  15.00, 820, "/catalogue/images/puma_1.jpeg");
INSERT INTO product VALUES ("03fef6ac-1896-4ce8-bd69-b798f85c6e0b", "Holy", "Socks fit for a Messiah. You too can experience walking in water with these special edition beauties. Each hole is lovingly proggled to leave smooth edges. The only sock approved by a higher power.",  99.99, 1, "/catalogue/images/holy_1.jpeg");
INSERT INTO product VALUES ("d3588630-ad8e-49df-bbd7-3167f7efb246", "YouTube.sock", "We were not paid to sell this sock. It's just a bit geeky.",  10.99, 801, "/catalogue/images/youtube_1.jpeg");
INSERT INTO product VALUES ("819e1fbf-8b7e-4f6d-811f-693534916a8b", "Figueroa", "enim officia aliqua excepteur esse deserunt quis aliquip nostrud anim",  14, 808, "/catalogue/images/WAT.jpg");
INSERT INTO product VALUES ("zzz4f044-b040-410d-8ead-4de0446aec7e", "Classic", "Keep it simple.",  12, 127, "/catalogue/images/classic.jpg");
INSERT INTO product VALUES ("3395a43e-2d88-40de-b95f-e00e1502085b", "Colourful", "proident occaecat irure et excepteur labore minim nisi amet irure",  18, 438, "/catalogue/images/colourful_socks.jpg");
INSERT INTO product VALUES ("837ab141-399e-4c1f-9abc-bace40296bac", "Cat socks", "consequat amet cupidatat minim laborum tempor elit ex consequat in",  15, 175, "/catalogue/images/catsocks.jpg");

INSERT INTO type (name) VALUES ("brown");
INSERT INTO type (name) VALUES ("geek");
INSERT INTO type (name) VALUES ("formal");
INSERT INTO type (name) VALUES ("blue");
INSERT INTO type (name) VALUES ("skin");
INSERT INTO type (name) VALUES ("red");
INSERT INTO type (name) VALUES ("action");
INSERT INTO type (name) VALUES ("sport");
INSERT INTO type (name) VALUES ("black");
INSERT INTO type (name) VALUES ("magic");
INSERT INTO type (name) VALUES ("green");

INSERT INTO product_type VALUES ("a0a4f044-b040-410d-8ead-4de0446aec7e", "4");
INSERT INTO product_type VALUES ("a0a4f044-b040-410d-8ead-4de0446aec7e", "5");
INSERT INTO product_type VALUES ("808a2de1-1aaa-4c25-a9b9-6612e8f29a38", "4");
INSERT INTO product_type VALUES ("808a2de1-1aaa-4c25-a9b9-6612e8f29a38", "6");
INSERT INTO product_type VALUES ("808a2de1-1aaa-4c25-a9b9-6612e8f29a38", "7");
INSERT INTO product_type VALUES ("808a2de1-1aaa-4c25-a9b9-6612e8f29a38", "3");
INSERT INTO product_type VALUES ("510a0d7e-8e83-4193-b483-e27e09ddc34d", "8");
INSERT INTO product_type VALUES ("510a0d7e-8e83-4193-b483-e27e09ddc34d", "9");
INSERT INTO product_type VALUES ("510a0d7e-8e83-4193-b483-e27e09ddc34d", "3");
INSERT INTO product_type VALUES ("03fef6ac-1896-4ce8-bd69-b798f85c6e0b", "10");
INSERT INTO product_type VALUES ("03fef6ac-1896-4ce8-bd69-b798f85c6e0b", "7");
INSERT INTO product_type VALUES ("d3588630-ad8e-49df-bbd7-3167f7efb246", "2");
INSERT INTO product_type VALUES ("d3588630-ad8e-49df-bbd7-3167f7efb246", "3");
INSERT INTO product_type VALUES ("819e1fbf-8b7e-4f6d-811f-693534916a8b", "3");
INSERT INTO product_type VALUES ("819e1fbf-8b7e-4f6d-811f-693534916a8b", "11");
INSERT INTO product_type VALUES ("819e1fbf-8b7e-4f6d-811f-693534916a8b", "4");
INSERT INTO product_type VALUES ("zzz4f044-b040-410d-8ead-4de0446aec7e", "1");
INSERT INTO product_type VALUES ("zzz4f044-b040-410d-8ead-4de0446aec7e", "11");
INSERT INTO product_type VALUES ("3395a43e-2d88-40de-b95f-e00e1502085b", "1");
INSERT INTO product_type VALUES ("3395a43e-2d88-40de-b95f-e00e1502085b", "4");
INSERT INTO product_type VALUES ("837ab141-399e-4c1f-9abc-bace40296bac", "1");
INSERT INTO product_type VALUES ("837ab141-399e-4c1f-9abc-bace40296bac", "11");
INSERT INTO product_type VALUES ("837ab141-399e-4c1f-9abc-bace40296bac", "3");


