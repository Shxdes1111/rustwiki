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
    icon VARCHAR(255)
);

CREATE TABLE weapon_ammo (
    weapon_item_id INTEGER REFERENCES weapon_item (id) ON DELETE CASCADE,
    ammo_id INTEGER REFERENCES ammo (id) ON DELETE CASCADE,
    PRIMARY KEY (weapon_item_id, ammo_id) -- Составной ключ защитит от дублей
);

-- mods
CREATE TABLE mods (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    icon VARCHAR(255)
);

CREATE TABLE weapon_mods (
    weapon_item_id INTEGER REFERENCES weapon_item(id) ON DELETE CASCADE,
    mod_id INTEGER REFERENCES mods(id) ON DELETE CASCADE,
    PRIMARY KEY (weapon_item_id, mod_id) -- Составной ключ, чтобы избежать дублей связей
);

-- ingredients
CREATE TABLE ingredients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    icon VARCHAR(255)
);

CREATE TABLE weapon_ingredients (
    weapon_item_id INTEGER REFERENCES weapon_item(id) ON DELETE CASCADE,
    ingredients_id INTEGER REFERENCES ingredients(id) ON DELETE CASCADE,
    amount INTEGER NOT NULL DEFAULT 1,
    PRIMARY KEY (weapon_item_id, ingredients_id)
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

INSERT INTO ammo (id, name, icon) VALUES
    (1, '5.56mm Rifle Ammo', '/icons/ammo/rifle.png'),
    (2, '9mm Pistol Ammo', '/icons/ammo/pistol.png'),
    (3, '12 Gauge Shells', '/icons/ammo/shells.png'),
    (4, 'Arrow', '/icons/ammo/arrow.png'),
    (5, 'Handmade Shell', '/icons/ammo/shells.png');

INSERT INTO weapon_ammo (weapon_item_id, ammo_id) VALUES
    (1, 1), -- AK-47 (id 1) использует 5.56mm (id 1)
    (2, 1), -- M4A4 (id 2) использует 5.56mm (id 1)
    (3, 2), -- MP5A4 (id 3) использует 9mm (id 2)
    (4, 2), -- Tommy Gun (id 4) использует 9mm (id 2)
    (5, 3), -- Pump Shotgun (id 5) использует 12 Gauge (id 3)
    (6, 5), -- Double Barrel (id 6) использует Handmade Shell (id 5)
    (7, 2), -- Revolver (id 7) использует 9mm (id 2)
    (8, 2), -- Semi-Auto Pistol (id 8) использует 9mm (id 2)
    (10, 4); -- Hunting Bow (id 10) использует Arrow (id 4)

INSERT INTO mods (name, icon) VALUES
    ('Silencer', '/icons/mods/silencer.png'),
    ('Red Dot Sight', '/icons/mods/red_dot.png'),
    ('Holographic Sight', '/icons/mods/holo.png'),
    ('8x Scope', '/icons/mods/scope.png'),
    ('Muzzle Boost', '/icons/mods/muzzle.png');

INSERT INTO weapon_mods (weapon_item_id, mod_id) VALUES
    (1, 1),   -- Silencer для AK-47 (id 1)
    (3, 2),   -- Red Dot Sight для MP5A4 (id 3)
    (5, 3),   -- Holographic Sight для Pump Shotgun (id 5)
    (10, 4),  -- 8x Scope для Hunting Bow (id 10)
    (7, 5),   -- Muzzle Boost для Revolver (id 7)
    (1, 4),   -- 8x Scope для AK-47
    (1, 3),   -- 8x Scope для AK-47
    (1, 2), 
    (1, 5); 

INSERT INTO ingredients (name, icon) VALUES
    ('Metal Fragments', '/icons/ingredients/metal_frags.png'),
    ('Wood', '/icons/ingredients/wood.png'),
    ('Cloth', '/icons/ingredients/cloth.png'),
    ('High Quality Metal', '/icons/ingredients/hqm.png'),
    ('Animal Fat', '/icons/ingredients/animal_fat.png');

INSERT INTO weapon_ingredients (weapon_item_id, ingredients_id, amount) VALUES
    (1,1, 250),
    (2,1, 150),
    (3,2, 200);