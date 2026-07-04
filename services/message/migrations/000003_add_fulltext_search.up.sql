ALTER TABLE messages ADD COLUMN content_tsvector tsvector;

CREATE INDEX idx_messages_content_fts ON messages USING GIN (content_tsvector);

CREATE OR REPLACE FUNCTION messages_content_tsvector_update() RETURNS trigger AS $$
BEGIN
    NEW.content_tsvector := to_tsvector('russian', COALESCE(NEW.content, ''));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER messages_content_tsvector_trigger
    BEFORE INSERT OR UPDATE OF content ON messages
    FOR EACH ROW
    EXECUTE FUNCTION messages_content_tsvector_update();

UPDATE messages SET content_tsvector = to_tsvector('russian', COALESCE(content, ''));
