ALTER TABLE weapon_item ADD COLUMN IF NOT EXISTS description TEXT;
ALTER TABLE weapon_item ADD COLUMN IF NOT EXISTS shortname VARCHAR(255);
ALTER TABLE weapon_item ADD COLUMN IF NOT EXISTS capacity INTEGER;
ALTER TABLE weapon_item ADD COLUMN IF NOT EXISTS time_to_craft INTEGER;

ALTER TABLE ammo ADD COLUMN IF NOT EXISTS icon VARCHAR(255);
ALTER TABLE mods ADD COLUMN IF NOT EXISTS icon VARCHAR(255);
ALTER TABLE ingredients ADD COLUMN IF NOT EXISTS amount INTEGER;
ALTER TABLE ingredients ADD COLUMN IF NOT EXISTS icon VARCHAR(255);

UPDATE weapon_item SET
  shortname = 'ak47',
  description = 'The AK-47 is a powerful automatic rifle, deadly at medium range. A favorite among Rust players.',
  capacity = 30,
  time_to_craft = 15
WHERE name = 'Assault Rifle';

UPDATE weapon_item SET
  shortname = 'm4a4',
  description = 'The M4A4 is a high-accuracy automatic rifle. Cannot be crafted — only found in loot crates.',
  capacity = 30
WHERE name = 'LR-300';

UPDATE weapon_item SET
  shortname = 'mp5a4',
  description = 'The MP5A4 is a compact submachine gun with a high rate of fire. Effective in close quarters.',
  capacity = 30,
  time_to_craft = 10
WHERE name = 'MP5A4';

UPDATE weapon_item SET
  shortname = 'tommy',
  description = 'The Thompson submachine gun, also known as the "Tommy Gun", is a classic automatic weapon.',
  capacity = 20,
  time_to_craft = 10
WHERE name = 'Thompson';

UPDATE weapon_item SET
  shortname = 'pump_shotgun',
  description = 'A pump-action shotgun with a tight spread. Devastating at close range.',
  capacity = 6,
  time_to_craft = 15
WHERE name = 'Pump Shotgun';

UPDATE weapon_item SET
  shortname = 'double_barrel',
  description = 'A double-barreled shotgun. Two shots, high damage, slow reload.',
  capacity = 2,
  time_to_craft = 10
WHERE name = 'Double Barrel Shotgun';

UPDATE weapon_item SET
  shortname = 'revolver',
  description = 'A six-shot revolver. Slower fire rate but packs a punch.',
  capacity = 6,
  time_to_craft = 5
WHERE name = 'Revolver';

UPDATE weapon_item SET
  shortname = 'semi_auto_pistol',
  description = 'A semi-automatic pistol. Reliable and easy to craft.',
  capacity = 8,
  time_to_craft = 5
WHERE name = 'Semi-Automatic Rifle';

UPDATE weapon_item SET
  shortname = 'machete',
  description = 'A sharp machete. No ammo needed, just swing.',
  time_to_craft = 5
WHERE name = 'Machete';

UPDATE weapon_item SET
  shortname = 'hunting_bow',
  description = 'A simple hunting bow. Silent but deadly with practice.',
  capacity = 1,
  time_to_craft = 5
WHERE name = 'Hunting Bow';

UPDATE weapon_item SET
  shortname = 'compound_bow',
  description = 'A modern compound bow. Higher damage and faster arrow velocity than the standard bow.',
  capacity = 1,
  time_to_craft = 10
WHERE name = 'Compound Bow';

UPDATE ammo SET icon = '/icons/ammo/rifle.png' WHERE name = '5.56mm Rifle Ammo';
UPDATE ammo SET icon = '/icons/ammo/rifle.png' WHERE name = 'HV 5.56mm Rifle Ammo';
UPDATE ammo SET icon = '/icons/ammo/rifle.png' WHERE name = 'Incendiary 5.56mm Rifle Ammo';
UPDATE ammo SET icon = '/icons/ammo/rifle.png' WHERE name = 'Explosive 5.56mm Rifle Ammo';
UPDATE ammo SET icon = '/icons/ammo/pistol.png' WHERE name = '9mm Pistol Ammo';
UPDATE ammo SET icon = '/icons/ammo/pistol.png' WHERE name = 'HV 9mm Pistol Ammo';
UPDATE ammo SET icon = '/icons/ammo/pistol.png' WHERE name = 'Incendiary 9mm Pistol Ammo';
UPDATE ammo SET icon = '/icons/ammo/shells.png' WHERE name = 'Handmade Shell';
UPDATE ammo SET icon = '/icons/ammo/shells.png' WHERE name = '12 Gauge Buckshot';
UPDATE ammo SET icon = '/icons/ammo/shells.png' WHERE name = '12 Gauge Slug';
UPDATE ammo SET icon = '/icons/ammo/shells.png' WHERE name = '12 Gauge Incendiary Shell';
UPDATE ammo SET icon = '/icons/ammo/arrow.png' WHERE name = 'Wooden Arrow';
UPDATE ammo SET icon = '/icons/ammo/arrow.png' WHERE name = 'High Velocity Arrow';
UPDATE ammo SET icon = '/icons/ammo/arrow.png' WHERE name = 'Bone Arrow';
UPDATE ammo SET icon = '/icons/ammo/arrow.png' WHERE name = 'Fire Arrow';

UPDATE mods SET icon = 'Weapon flashlight' WHERE name = 'Weapon flashlight';
UPDATE mods SET icon = '/icons/mods/Weapon Lasersight.png' WHERE name = 'Weapon Lasersight';
UPDATE mods SET icon = '/icons/mods/Soda Can Silencer.png' WHERE name = 'Soda Can Silencer';
UPDATE mods SET icon = '/icons/mods/Oil Filter Silencer.png' WHERE name = 'Oil Filter Silencer';
UPDATE mods SET icon = '/icons/mods/silencer.png' WHERE name = 'Military Silencer';
UPDATE mods SET icon = 'Simple Handmade Sight' WHERE name = 'Simple Handmade Sight';
UPDATE mods SET icon = '/icons/mods/holo.png' WHERE name = 'Holosight';
UPDATE mods SET icon = '/icons/mods/scope.png' WHERE name = '8x Zoom Scope';
UPDATE mods SET icon = 'Variable Zoom Scope' WHERE name = 'Variable Zoom Scope';
UPDATE mods SET icon = 'Gas Compression Overdrive' WHERE name = 'Gas Compression Overdrive';
UPDATE mods SET icon = '/icons/mods/muzzle.png' WHERE name = 'Muzzle Boost';
UPDATE mods SET icon = '/icons/mods/muzzle.png' WHERE name = 'Muzzle Brake';
UPDATE mods SET icon = 'Burst Module' WHERE name = 'Burst Module';
UPDATE mods SET icon = 'Extended Magazine' WHERE name = 'Extended Magazine';
UPDATE mods SET icon = 'Targeting Attachment' WHERE name = 'Targeting Attachment';