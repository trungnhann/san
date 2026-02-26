CREATE TABLE posts (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    abstract TEXT,
    body TEXT NOT NULL,
    published BOOLEAN NOT NULL DEFAULT FALSE,
    publish_date TIMESTAMP WITH TIME ZONE,
    location VARCHAR(255),
    lat FLOAT,
    lon FLOAT,
    locale VARCHAR(10) DEFAULT 'en-US',
    tags TEXT[],
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_slug ON posts(slug);
CREATE INDEX idx_posts_created_at ON posts(created_at);
