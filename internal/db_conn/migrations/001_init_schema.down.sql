
DROP INDEX IF EXISTS idx_matches_unique;
DROP INDEX IF EXISTS idx_matches_user_b_id;
DROP INDEX IF EXISTS idx_matches_user_a_id;
DROP TABLE IF EXISTS matches;

DROP INDEX IF EXISTS idx_swipes_unique;
DROP INDEX IF EXISTS idx_swipes_user_card;
DROP INDEX IF EXISTS idx_swipes_user_id;
DROP TABLE IF EXISTS swipes;

DROP INDEX IF EXISTS idx_game_cards_created;
DROP INDEX IF EXISTS idx_game_cards_game_id;
DROP INDEX IF EXISTS idx_game_cards_user_id;
DROP TABLE IF EXISTS game_cards;

DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;