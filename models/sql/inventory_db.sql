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

CREATE TABLE IF NOT EXISTS warehouse (
  id SERIAL PRIMARY KEY,
  warehouse_state VARCHAR(255) NOT NULL,
  warehouse_city VARCHAR(255) NOT NULL,
  CONSTRAINT uq_warehouse_wswc UNIQUE (warehouse_state, warehouse_city)
);

CREATE TABLE IF NOT EXISTS inventory (
  id SERIAL PRIMARY KEY,
  product_name VARCHAR(255) NOT NULL,
  price DECIMAL NOT NULL,
  warehouse_id INT NOT NULL,
  CONSTRAINT fk_inventory_warehouse_id FOREIGN KEY (warehouse_id) REFERENCES warehouse(id)
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
