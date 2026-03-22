CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       username VARCHAR(255) NOT NULL UNIQUE,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       password_hash VARCHAR(255) NOT NULL,
                       age INT NOT NULL,
                       language VARCHAR(10) NOT NULL,
                       discord VARCHAR(100),
                       telegram VARCHAR(100),
                       region VARCHAR(50) NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);


CREATE TABLE game_cards(
            id BIGSERIAL PRIMARY KEY,
            user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            game_id  VARCHAR(50) NOT NULL,
            rank    VARCHAR(50) NOT NULL,
            role    VARCHAR(50) NOT NULL,
            description TEXT NOT NULL,
            is_active BOOLEAN NOT NULL DEFAULT true,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_game_cards_user_id ON game_cards(user_id);
CREATE INDEX idx_game_cards_game_id ON game_cards(game_id) WHERE is_active = true;
CREATE INDEX idx_game_cards_created ON game_cards(created_at DESC);

CREATE TABLE swipes (
    id BIGSERIAL PRIMARY KEY,
        user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        target_card_id  BIGINT NOT NULL REFERENCES game_cards(id) ON DELETE CASCADE,
    action VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_swipes_user_id ON swipes(user_id);
CREATE UNIQUE INDEX idx_swipes_unique ON swipes(user_id, target_card_id);

CREATE TABLE matches(
    id BIGSERIAL PRIMARY KEY,
    user_a_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_b_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    game_id VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_matches_user_a_id ON matches(user_a_id);
CREATE INDEX idx_matches_user_b_id ON matches(user_b_id);
CREATE UNIQUE INDEX idx_matches_unique ON matches(
                                                LEAST(user_a_id, user_b_id),
                                                GREATEST(user_a_id, user_b_id),
                                                game_id

);
