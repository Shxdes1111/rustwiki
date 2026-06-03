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