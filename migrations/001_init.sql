-- Simple IM Database Schema

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(100),
    avatar VARCHAR(500),
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Friends table
CREATE TABLE IF NOT EXISTS friends (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    friend_id BIGINT NOT NULL REFERENCES users(id),
    status SMALLINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, friend_id)
);

CREATE INDEX IF NOT EXISTS idx_friend_user ON friends(user_id, friend_id);

-- Groups table
CREATE TABLE IF NOT EXISTS groups (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    owner_id BIGINT NOT NULL REFERENCES users(id),
    avatar VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Group members table
CREATE TABLE IF NOT EXISTS group_members (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id),
    role SMALLINT DEFAULT 0,
    joined_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(group_id, user_id)
);

-- Messages table
CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    sender_id BIGINT NOT NULL REFERENCES users(id),
    receiver_id BIGINT,
    group_id BIGINT,
    msg_type SMALLINT NOT NULL,
    content TEXT,
    file_url VARCHAR(500),
    file_name VARCHAR(255),
    file_size BIGINT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_messages_sender ON messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_messages_receiver ON messages(receiver_id, created_at);
CREATE INDEX IF NOT EXISTS idx_messages_group ON messages(group_id, created_at);

-- Files table
CREATE TABLE IF NOT EXISTS files (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    filename VARCHAR(255) NOT NULL,
    filepath VARCHAR(500) NOT NULL,
    filesize BIGINT NOT NULL,
    mimetype VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_files_user ON files(user_id);
