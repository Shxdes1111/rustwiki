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
    description TEXT,
    shortname VARCHAR(255),
    capacity INTEGER,
    time_to_craft INTEGER,
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
    icon VARCHAR(255),
    weapon_item_id INTEGER REFERENCES weapon_item(id)
);

-- mods
CREATE TABLE mods (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    icon VARCHAR(255),
    weapon_item_id INTEGER REFERENCES weapon_item(id)
);

-- ingredients
CREATE TABLE ingredients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    amount INTEGER,
    icon VARCHAR(255),
    weapon_item_id INTEGER REFERENCES weapon_item(id),
    clothing_item_id INTEGER REFERENCES clothing_item(id)
);

-- insert data
INSERT INTO category (id, name) VALUES 
    (1, 'weapons'), 
    (2, 'armor');

INSERT INTO weapon_item (name, type, firemode, craftable, stacksize, description, shortname, capacity, time_to_craft, category_id) VALUES
    ('AK-47',                'range', 'automatic', true,  1, 'The AK-47 is a powerful automatic rifle, deadly at medium range. A favorite among Rust players.',           'ak47',              30, 15, 1),
    ('M4A4',                 'range', 'automatic', false, 1, 'The M4A4 is a high-accuracy automatic rifle. Cannot be crafted — only found in loot crates.',                 'm4a4',              30, NULL, 1),
    ('MP5A4',                'range', 'automatic', true,  1, 'The MP5A4 is a compact submachine gun with a high rate of fire. Effective in close quarters.',               'mp5a4',             30, 10, 1),
    ('Tommy Gun',            'range', 'automatic', true,  1, 'The Thompson submachine gun, also known as the "Tommy Gun", is a classic automatic weapon.',                'tommy',             20, 10, 1),
    ('Pump Shotgun',         'range', 'semi',      true,  1, 'A pump-action shotgun with a tight spread. Devastating at close range.',                                   'pump_shotgun',       6, 15, 1),
    ('Double Barrel',        'range', 'double',    true,  1, 'A double-barreled shotgun. Two shots, high damage, slow reload.',                                           'double_barrel',      2, 10, 1),
    ('Revolver',             'range', 'semi',      true,  1, 'A six-shot revolver. Slower fire rate but packs a punch.',                                                  'revolver',           6, 5, 1),
    ('Semi-Auto Pistol',     'range', 'semi',      true,  1, 'A semi-automatic pistol. Reliable and easy to craft.',                                                       'semi_auto_pistol',   8, 5, 1),
    ('Machete',              'melee', 'none',      true,  1, 'A sharp machete. No ammo needed, just swing.',                                                               'machete',            NULL, 5, 1),
    ('Hunting Bow',          'range', 'none',      true,  1, 'A simple hunting bow. Silent but deadly with practice.',                                                     'hunting_bow',        1, 5, 1),
    ('Compound Bow',         'range', 'none',      true,  1, 'A modern compound bow. Higher damage and faster arrow velocity than the standard bow.',                      'compound_bow',       1, 10, 1);

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

INSERT INTO ammo (name, icon, weapon_item_id) VALUES
    ('5.56mm Rifle Ammo', '/icons/ammo/rifle.png', 1),
    ('9mm Pistol Ammo', '/icons/ammo/pistol.png', 3),
    ('12 Gauge Shells', '/icons/ammo/shells.png', 5),
    ('Arrow', '/icons/ammo/arrow.png', 10),
    ('Handmade Shell', '/icons/ammo/shells.png', 6);

INSERT INTO mods (name, icon, weapon_item_id) VALUES
    ('Silencer', '/icons/mods/silencer.png', 1),
    ('Red Dot Sight', '/icons/mods/red_dot.png', 3),
    ('Holographic Sight', '/icons/mods/holo.png', 5),
    ('8x Scope', '/icons/mods/scope.png', 10),
    ('Muzzle Boost', '/icons/mods/muzzle.png', 7);

INSERT INTO ingredients (name, amount, icon, weapon_item_id, clothing_item_id) VALUES
    ('Metal Fragments',    250, '/icons/ingredients/metal_frags.png', 1, NULL),
    ('Wood',              200, '/icons/ingredients/wood.png',        10, NULL),
    ('Cloth',              15, '/icons/ingredients/cloth.png',       NULL, 2),
    ('High Quality Metal', 10, '/icons/ingredients/hqm.png',         NULL, 9),
    ('Animal Fat',          8, '/icons/ingredients/animal_fat.png',  NULL, 8);