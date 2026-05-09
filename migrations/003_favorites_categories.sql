INSERT OR IGNORE INTO categories (name, created_at, updated_at)
VALUES
    ('text', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('image', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('link', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('file', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT OR IGNORE INTO clipboard_item_categories (clipboard_item_id, category_id, created_at)
SELECT ci.id, c.id, CURRENT_TIMESTAMP
FROM clipboard_items AS ci
JOIN categories AS c ON c.name = ci.item_type
LEFT JOIN clipboard_item_categories AS cic ON cic.clipboard_item_id = ci.id
WHERE cic.clipboard_item_id IS NULL
  AND ci.item_type IN ('text', 'image', 'link', 'file');
