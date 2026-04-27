-- up migration

-- category
CREATE TABLE category (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- weapon_item
CREATE TABLE weapon_item (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('range', 'melee')),
    firemode VARCHAR(50) NOT NULL,
    craftable BOOLEAN NOT NULL DEFAULT false,
    stacksize INTEGER NOT NULL DEFAULT 1,
    category_id INTEGER REFERENCES category(id)
);

-- clothing_item
CREATE TABLE clothing_item (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    equipmentslot VARCHAR(50) NOT NULL,
    protection INTEGER,
    craftable BOOLEAN NOT NULL DEFAULT false,
    stacksize INTEGER NOT NULL DEFAULT 1,
    category_id INTEGER REFERENCES category(id)
);

-- ammo
CREATE TABLE ammo (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    weapon_item_id INTEGER REFERENCES weapon_item(id)
);

-- mods
CREATE TABLE mods (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    weapon_item_id INTEGER REFERENCES weapon_item(id)
);

-- ingredients
CREATE TABLE ingredients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    weapon_item_id INTEGER REFERENCES weapon_item(id),
    clothing_item_id INTEGER REFERENCES clothing_item(id)
);

-- insert data
INSERT INTO category (id, name) VALUES 
    (1, 'weapons'), 
    (2, 'armor');

INSERT INTO weapon_item (name, type, firemode, craftable, stacksize, category_id) VALUES
    ('AK-47', 'range', 'automatic', true, 1, 1),
    ('M4A4', 'range', 'automatic', false, 1, 1),
    ('MP5A4', 'range', 'automatic', true, 1, 1),
    ('Tommy Gun', 'range', 'automatic', true, 1, 1),
    ('Pump Shotgun', 'range', 'semi', true, 1, 1),
    ('Double Barrel', 'range', 'double', true, 1, 1),
    ('Revolver', 'range', 'semi', true, 1, 1),
    ('Semi-Auto Pistol', 'range', 'semi', true, 1, 1),
    ('Machete', 'melee', 'none', true, 1, 1),
    ('Hunting Bow', 'range', 'none', true, 1, 1),
    ('Compound Bow', 'range', 'none', true, 1, 1);

INSERT INTO clothing_item (name, equipmentslot, protection, craftable, stacksize, category_id) VALUES
    ('Vest', 'body', 70, true, 1, 2),
    ('Hoodie', 'body', 35, true, 1, 2),
    ('Jacket', 'body', 45, true, 1, 2),
    ('T-Shirt', 'body', 0, true, 1, 2),
    ('Jeans', 'legs', 35, true, 1, 2),
    ('Boots', 'feet', 10, true, 1, 2),
    ('Baseball Cap', 'head', 0, true, 1, 2),
    ('Balaclava', 'head', 5, true, 1, 2),
    ('Bone Helmet', 'head', 15, true, 1, 2),
    ('Coffee Can Helmet', 'head', 20, true, 1, 2);