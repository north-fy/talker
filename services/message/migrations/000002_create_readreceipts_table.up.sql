CREATE TABLE read_receipts (
    chat_id UUID NOT NULL,
    user_id UUID NOT NULL,
    last_read_message_id UUID,
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (chat_id, user_id)
);