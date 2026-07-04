DROP TRIGGER IF EXISTS messages_content_tsvector_trigger ON messages;
DROP FUNCTION IF EXISTS messages_content_tsvector_update();
DROP INDEX IF EXISTS idx_messages_content_fts;
ALTER TABLE messages DROP COLUMN IF EXISTS content_tsvector;
