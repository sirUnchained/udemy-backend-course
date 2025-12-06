CREATE DATABASE social;

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password bytea NOT NULL, -- password is hashed so we use bytea
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP(0)
    WITH
        TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tags text[] NOT NULL DEFAULT '{}'::text[],
    version INT DEFAULT 0,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS comments (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    post_id bigint NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP(0)
    WITH
        TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS followers (
    user_id bigint NOT NULL,
    follower_id bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP(0)
    WITH
        TIME ZONE NOT NULL DEFAULT NOW(),
        PRIMARY KEY (user_id, follower_id), -- composite key: user A can follow B only one time not 2 times!
        Foreign Key (user_id) REFERENCES users (id) ON DELETE CASCADE,
        Foreign Key (follower_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_invitations (
    token bytea PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expiry TIMESTAMP(0)
    WITH
        TIME ZONE NOT NULL
);

CREATE Extention IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_comments_content ON comments USING gin (content gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_posts_title ON posts USING gin (title gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING gin (tags);

CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);

CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);