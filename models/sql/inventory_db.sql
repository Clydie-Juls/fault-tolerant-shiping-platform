DROP TABLE IF EXISTS warehouse;
DROP TABLE IF EXISTS shipment;
DROP TABLE IF EXISTS person;
DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS location;


CREATE TABLE IF NOT EXISTS location (
  id SERIAL PRIMARY KEY,
  longitude DECIMAL NOT NULL,
  latitude DECIMAL NOT NULL
);
CREATE TABLE IF NOT EXISTS person (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  age INT,
  location_id INT NOT NULL,
  CONSTRAINT fk_person_location_id FOREIGN KEY(location_id) REFERENCES location(id)
);

CREATE TABLE IF NOT EXISTS inventory (
  id SERIAL PRIMARY KEY,
  product_name VARCHAR(255) NOT NULL,
  price DECIMAL NOT NULL
);

CREATE TABLE IF NOT EXISTS shipment (
  id SERIAL PRIMARY KEY,
  person_id INT NOT NULL,
  product_name VARCHAR(255) NOT NULL,
  sell_amount DECIMAL NOT NULL,
  inventory_id INT NOT NULL,
  CONSTRAINT fk_shipment_person_id FOREIGN KEY(person_id) REFERENCES person(id),
  CONSTRAINT fk_shipment_inventory_id FOREIGN KEY(inventory_id) REFERENCES inventory(id)
);

CREATE TABLE IF NOT EXISTS warehouse (
  id SERIAL PRIMARY KEY,
  location_id INT NOT NULL,
  inventory_id INT NOT NULL,
  CONSTRAINT fk_warehouse_location_id FOREIGN KEY(location_id) REFERENCES location(id),
  CONSTRAINT fk_warehouse_inventory_id FOREIGN KEY(inventory_id) REFERENCES inventory(id)
);
