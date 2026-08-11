CREATE TABLE chat_settings (
    chat_id BIGINT PRIMARY KEY REFERENCES chats(id) ON DELETE CASCADE,
    is_private BOOLEAN DEFAULT false,
    allow_messages_from_all BOOLEAN DEFAULT true,
    allow_media BOOLEAN DEFAULT true,
    allow_reactions BOOLEAN DEFAULT true,
    message_ttl_seconds INTEGER DEFAULT 0,
    language VARCHAR(10) DEFAULT 'ru',
    is_announcement BOOLEAN DEFAULT false
);
