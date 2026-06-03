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
    ('Assault Rifle',        'range', 'automatic', true,  1, 'The AK-47 is a powerful automatic rifle, deadly at medium range. A favorite among Rust players.',           'ak47',              30, 15, 1),
    ('LR-300',               'range', 'automatic', false, 1, 'The M4A4 is a high-accuracy automatic rifle. Cannot be crafted — only found in loot crates.',                 'm4a4',              30, NULL, 1),
    ('MP5A4',                'range', 'automatic', true,  1, 'The MP5A4 is a compact submachine gun with a high rate of fire. Effective in close quarters.',               'mp5a4',             30, 10, 1),
    ('Thompson',             'range', 'automatic', true,  1, 'The Thompson submachine gun, also known as the "Tommy Gun", is a classic automatic weapon.',                'tommy',             20, 10, 1),
    ('Pump Shotgun',         'range', 'semi',      true,  1, 'A pump-action shotgun with a tight spread. Devastating at close range.',                                   'pump_shotgun',       6, 15, 1),
    ('Double Barrel Shotgun','range', 'double',    true,  1, 'A double-barreled shotgun. Two shots, high damage, slow reload.',                                           'double_barrel',      2, 10, 1),
    ('Revolver',             'range', 'semi',      true,  1, 'A six-shot revolver. Slower fire rate but packs a punch.',                                                  'revolver',           6, 5, 1),
    ('Semi-Automatic Rifle', 'range', 'semi',      true,  1, 'A semi-automatic pistol. Reliable and easy to craft.',                                                       'semi_auto_pistol',   8, 5, 1),
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
    (2, 'HV 5.56mm Rifle Ammo', '/icons/ammo/rifle.png'),
    (3, 'Incendiary 5.56mm Rifle Ammo', '/icons/ammo/rifle.png'),
    (4, 'Explosive 5.56mm Rifle Ammo', '/icons/ammo/rifle.png'),
    (5, '9mm Pistol Ammo', '/icons/ammo/pistol.png'),
    (6, 'HV 9mm Pistol Ammo', '/icons/ammo/pistol.png'),
    (7, 'Incendiary 9mm Pistol Ammo', '/icons/ammo/pistol.png'),
    (8, 'Handmade Shell', '/icons/ammo/shells.png'),
    (9, '12 Gauge Buckshot', '/icons/ammo/shells.png'),
    (10,'12 Gauge Slug', '/icons/ammo/shells.png'),
    (11,'12 Gauge Incendiary Shell', '/icons/ammo/shells.png'),
    (12,'Wooden Arrow', '/icons/ammo/arrow.png'),
    (13,'High Velocity Arrow', '/icons/ammo/arrow.png'),
    (14,'Bone Arrow', '/icons/ammo/arrow.png'),
    (15,'Fire Arrow', '/icons/ammo/arrow.png');

INSERT INTO weapon_ammo (weapon_item_id, ammo_id) VALUES
    (1, 1),(1, 2),(1, 3),(1, 4),        -- AK-47 (id 1) 
    (2, 1),(2, 2),(2, 3),(2, 4),        -- M4A4 (id 2) 
    (3, 5),(3, 6),(3, 7),               -- MP5A4 (id 3) 
    (4, 5),(4, 6),(4, 7),               -- Tommy Gun (id 4) 
    (5, 8),(5, 9),(5, 10),(5, 11),      -- Pump Shotgun (id 5) 
    (6, 8),(6, 9),(6, 10),(6, 11),      -- Double Barrel (id 6)
    (7, 5),(7, 6),(7, 7),               -- Revolver (id 7) 
    (8, 1),(8, 2),(8, 3),(8, 4),        -- Semi-Auto Pistol (id 8) 
    (10, 12),(10, 13),(10, 14),(10, 15),-- Hunting Bow (id 10)
    (11, 12),(11, 13),(11, 14),(11, 15);-- Hunting Bow (id 10)

INSERT INTO mods (name, icon) VALUES
    ('Weapon flashlight', 'Weapon flashlight'),                     --(id 1)  свет
    ('Weapon Lasersight', '/icons/mods/Weapon Lasersight.png'),     --(id 2)  лазер
    ('Soda Can Silencer', '/icons/mods/Soda Can Silencer.png'),     --(id 3)  глуш     
    ('Oil Filter Silencer', '/icons/mods/Oil Filter Silencer.png'), --(id 4)  глуш     
    ('Military Silencer', '/icons/mods/silencer.png'),              --(id 5)  глуш
    ('Simple Handmade Sight', 'Simple Handmade Sight'),             --(id 6)  соло прицел
    ('Holosight', '/icons/mods/holo.png'),                          --(id 7)  холик
    ('8x Zoom Scope', '/icons/mods/scope.png'),                     --(id 8)  8х 
    ('Variable Zoom Scope', 'Variable Zoom Scope'),                 --(id 9)  16х
    ('Gas Compression Overdrive', 'Gas Compression Overdrive'),     --(id 10) газовый компрессор
    ('Muzzle Boost', '/icons/mods/muzzle.png'),                     --(id 11) ускоритель
    ('Muzzle Brake', '/icons/mods/muzzle.png'),                     --(id 12) тормоз
    ('Burst Module', 'Burst Module'),                               --(id 13) трипл
    ('Extended Magazine', 'Extended Magazine'),                     --(id 14) увелмаг
    ('Targeting Attachment', 'Targeting Attachment');               --(id 15) метки

INSERT INTO weapon_mods (weapon_item_id, mod_id) VALUES
    (1,1),(1,2),(1,3),(1,4),(1,5),(1,6),(1,7),(1,8),(1,9),(1,11),(1,12),(1,14),(1,15),   
    (2,1),(2,2),(2,3),(2,4),(2,5),(2,6),(2,7),(2,8),(2,9),(2,11),(2,12),(2,14),(2,15),
    (3,1),(3,2),(3,3),(3,4),(3,5),(3,6),(3,7),(3,8),(3,9),(3,11),(3,12),(3,14), 
    (4,1),(4,2),(4,3),(4,4),(4,5),(4,6),(4,7),(4,8),(4,9),(4,10),(4,11),(4,12),(4,13),(4,14),(4,15), 
    (5,1),(5,2),(5,3),(5,4),(5,5),(5,6),(5,7),(5,8),(5,9),(5,11),(5,12),      
    (6,1),(6,2),(6,3),(6,4),(6,5),(6,6),(6,7),(6,8),(6,9),(6,11),(6,12),        
    (7,3),(7,4),(7,5),(7,11),(7,12),
    (8,1),(8,2),(8,3),(8,4),(8,5),(8,6),(8,7),(8,8),(8,9),(8,10),(8,11),(8,12),(8,14),(8,15);    

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