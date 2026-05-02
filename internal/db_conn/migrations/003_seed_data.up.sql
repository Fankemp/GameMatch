-- Seed test users (password for all: password123)
-- bcrypt hash of "password123"
INSERT INTO users (username, email, password_hash, age, language, discord, telegram, region) VALUES
('KazFrag',     'kaz@test.com',     '$2a$10$3VDDZiwLs3tpsuHLK04CuuFmaJ94bj818PlLOVqkioKRNePDiQ4ZO', 21, 'RU', 'KazFrag#1234',     '@kazfrag',     'CIS'),
('ShadowBlade', 'shadow@test.com',  '$2a$10$3VDDZiwLs3tpsuHLK04CuuFmaJ94bj818PlLOVqkioKRNePDiQ4ZO', 19, 'EN', 'ShadowBlade#5678', '@shadowblade', 'EU Central'),
('NeonRush',    'neon@test.com',    '$2a$10$3VDDZiwLs3tpsuHLK04CuuFmaJ94bj818PlLOVqkioKRNePDiQ4ZO', 24, 'RU', 'NeonRush#9012',    '@neonrush',    'CIS'),
('IceQueen',    'ice@test.com',     '$2a$10$3VDDZiwLs3tpsuHLK04CuuFmaJ94bj818PlLOVqkioKRNePDiQ4ZO', 22, 'EN', 'IceQueen#3456',    '@icequeen',    'EU West'),
('PhoenixKZ',   'phoenix@test.com', '$2a$10$3VDDZiwLs3tpsuHLK04CuuFmaJ94bj818PlLOVqkioKRNePDiQ4ZO', 20, 'KZ', 'PhoenixKZ#7890',   '@phoenixkz',   'CIS'),
('VoltStrike',  'volt@test.com',    '$2a$10$3VDDZiwLs3tpsuHLK04CuuFmaJ94bj818PlLOVqkioKRNePDiQ4ZO', 25, 'EN', 'VoltStrike#1111',  '@voltstrike',  'NA East'),
('AstraMain',   'astra@test.com',   '$2a$10$3VDDZiwLs3tpsuHLK04CuuFmaJ94bj818PlLOVqkioKRNePDiQ4ZO', 18, 'RU', 'AstraMain#2222',   '@astramain',   'CIS'),
('TurkishAce',  'turk@test.com',    '$2a$10$3VDDZiwLs3tpsuHLK04CuuFmaJ94bj818PlLOVqkioKRNePDiQ4ZO', 23, 'TR', 'TurkishAce#3333',  '@turkishace',  'EU Central'),
('PolandSniper','poland@test.com',  '$2a$10$3VDDZiwLs3tpsuHLK04CuuFmaJ94bj818PlLOVqkioKRNePDiQ4ZO', 17, 'PL', 'PolandSniper#4444','@polandsniper','EU Central'),
('BerlinWall',  'berlin@test.com',  '$2a$10$3VDDZiwLs3tpsuHLK04CuuFmaJ94bj818PlLOVqkioKRNePDiQ4ZO', 26, 'DE', 'BerlinWall#5555',  '@berlinwall',  'EU Central');

-- Seed profiles
INSERT INTO profiles (user_id, bio, avatar_url) VALUES
(1,  'Ищу тиму на ранкед, играю каждый вечер',  ''),
(2,  'Competitive Valorant player, LFT',         ''),
(3,  'Контроллер мейн, 2000+ часов',             ''),
(4,  'Support main, looking for duo',             ''),
(5,  'Астана, играю по вечерам',                  ''),
(6,  'Ex-CSGO player, grinding Valorant',         ''),
(7,  'Astra/Omen main, smoke god',                ''),
(8,  'Turkish player looking for EU team',        ''),
(9,  'Young and hungry, aiming for Immortal',     ''),
(10, 'Casual player, just vibing',                '');

-- Seed game cards (Valorant)
INSERT INTO game_cards (user_id, game_id, rank, role, description) VALUES
(1,  'valorant', 'Diamond',   'Duelist',     'Jett/Raze мейн, ищу сентинела для дуо'),
(1,  'valorant', 'Diamond',   'Initiator',   'Могу играть Sova/Fade на подмене'),
(2,  'valorant', 'Immortal',  'Sentinel',    'Cypher main, 300+ ADR, looking for aggressive duelist'),
(3,  'valorant', 'Platinum',  'Controller',  'Omen/Astra, курю как бог, ищу дуэлиста'),
(4,  'valorant', 'Diamond',   'Sentinel',    'Killjoy main, solid anchor. Need a good IGL'),
(5,  'valorant', 'Gold',      'Duelist',     'Фоникс мейн, учусь, ищу терпеливых тиммейтов'),
(6,  'valorant', 'Ascendant', 'Initiator',   'Sova/Fade/KayO, great util usage'),
(7,  'valorant', 'Platinum',  'Controller',  'Астра мейн, 1500 часов, играю CIS серверы'),
(8,  'valorant', 'Diamond',   'Duelist',     'Reyna/Jett player, 250+ ADR, need smokes'),
(9,  'valorant', 'Silver',    'Sentinel',    'Learning Cypher, want to improve'),
(10, 'valorant', 'Gold',      'Initiator',   'Casual Sova, just for fun');

-- Seed game cards (CS2)
INSERT INTO game_cards (user_id, game_id, rank, role, description) VALUES
(1,  'cs2', 'Diamond',   'Duelist',    'AWPer, 3000 часов в CS, ищу тиму'),
(2,  'cs2', 'Immortal',  'Sentinel',   'Anchor player, great crosshair placement'),
(3,  'cs2', 'Platinum',  'Controller', 'IGL, smoke lineups for all maps'),
(6,  'cs2', 'Ascendant', 'Duelist',    'Entry fragger, 10+ years of CS');

-- Seed some swipes (mutual likes for matches)
INSERT INTO swipes (user_id, target_card_id, action) VALUES
(1, 3, 'like'),   -- KazFrag likes ShadowBlade's Immortal Sentinel card
(2, 1, 'like'),   -- ShadowBlade likes KazFrag's Diamond Duelist card (MUTUAL → match)
(1, 4, 'like'),   -- KazFrag likes NeonRush's Platinum Controller card
(3, 1, 'like'),   -- NeonRush likes KazFrag's Diamond Duelist card (MUTUAL → match)
(4, 1, 'like'),   -- IceQueen likes KazFrag's Diamond Duelist card
(5, 3, 'like'),   -- PhoenixKZ likes ShadowBlade's Immortal Sentinel card
(2, 4, 'dislike'),-- ShadowBlade dislikes NeonRush
(6, 1, 'like'),   -- VoltStrike likes KazFrag's Diamond Duelist card
(7, 8, 'like'),   -- AstraMain likes TurkishAce's Diamond Duelist card
(8, 7, 'like');   -- TurkishAce likes AstraMain's Platinum Controller (MUTUAL → match)

-- Seed matches (from mutual likes above)
INSERT INTO matches (user_a_id, user_b_id, game_id) VALUES
(1, 2, 'valorant'),  -- KazFrag & ShadowBlade
(1, 3, 'valorant'),  -- KazFrag & NeonRush
(7, 8, 'valorant');  -- AstraMain & TurkishAce
