CREATE DATABASE social;

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password bytea NOT NULL, -- password is hashed so we use bytea
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