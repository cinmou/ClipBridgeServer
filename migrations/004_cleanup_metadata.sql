ALTER TABLE clipboard_items ADD COLUMN size_bytes INTEGER NOT NULL DEFAULT 0;
ALTER TABLE clipboard_items ADD COLUMN expires_at TEXT;

UPDATE clipboard_items
SET size_bytes = LENGTH(COALESCE(text_content, ''))
WHERE size_bytes = 0;
