CREATE TABLE chats (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type INTEGER NOT NULL DEFAULT 2,          -- 1 private / 2 group / 3 channel
    created_by BIGINT NOT NULL DEFAULT 0,
    avatar_url TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE chat_members (
    chat_id BIGINT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,
    role INTEGER NOT NULL DEFAULT 1,          -- 1 member / 2 moderator / 3 admin / 4 owner
    joined_at TIMESTAMP DEFAULT NOW(),
    last_read_at TIMESTAMP,
    unread_count BIGINT DEFAULT 0,
    PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE invite_links (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    code VARCHAR(50) UNIQUE NOT NULL,
    max_uses INTEGER DEFAULT 0,
    used_count INTEGER DEFAULT 0,
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    created_by BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_chat_members_user ON chat_members(user_id);
CREATE INDEX idx_chat_members_role ON chat_members(role);
CREATE INDEX idx_invite_links_code ON invite_links(code) WHERE is_active = true;
CREATE INDEX idx_chats_name ON chats USING GIN (to_tsvector('russian', name));
