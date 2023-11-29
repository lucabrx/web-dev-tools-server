CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT NOW(),
    name text,
    email text NOT NULL UNIQUE,
    image_url text,
    role text NOT NULL DEFAULT 'user',
    version integer NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp NOT NULL DEFAULT NOW(),
    scope text NOT NULL
);

CREATE INDEX IF NOT EXISTS users_email_idx ON users USING GIN (to_tsvector('simple', email));

CREATE TABLE IF NOT EXISTS tools (
    id bigserial PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    category text NOT NULL ,
    description text NOT NULL,
    website text NOT NULL,
    image_url text,
    published boolean NOT NULL DEFAULT false,
    version integer NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS categories (
    id bigserial PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    published boolean NOT NULL DEFAULT false,
    version integer NOT NULL DEFAULT 1
);

CREATE INDEX idx_tools_search ON tools USING gin(to_tsvector('simple', name || ' ' || category));

-- insert default categories
INSERT INTO categories (name, published) VALUES ('DB', true);
INSERT INTO categories (name, published) VALUES ('API', true);
INSERT INTO categories (name, published) VALUES ('UI', true);
INSERT INTO categories (name, published) VALUES ('Animations', true);
INSERT INTO categories (name, published) VALUES ('Testing', true);
INSERT INTO categories (name, published) VALUES ('Books', true);
INSERT INTO categories (name, published) VALUES ('AI', true);
INSERT INTO categories (name, published) VALUES ('Auth', true);
INSERT INTO categories (name, published) VALUES ('DevOps', true);
INSERT INTO categories (name, published) VALUES ('Headless', true);
INSERT INTO categories (name, published) VALUES ('CMS', true);
INSERT INTO categories (name, published) VALUES ('State', true);
INSERT INTO categories (name, published) VALUES ('CSS', true);
INSERT INTO categories (name, published) VALUES ('Tutorials', true);
INSERT INTO categories (name, published) VALUES ('Icons', true);
INSERT INTO categories (name, published) VALUES ('Typography', true);
INSERT INTO categories (name, published) VALUES ('Frameworks', true);
INSERT INTO categories (name, published) VALUES ('Forms', true);