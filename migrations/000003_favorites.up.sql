CREATE TABLE IF NOT EXISTS favorites(
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    tool_id bigint NOT NULL REFERENCES tools ON DELETE CASCADE
);