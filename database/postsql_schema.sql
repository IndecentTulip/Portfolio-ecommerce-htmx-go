CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name VARCHAR,
    profileImage BYTEA
);

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    name VARCHAR,
    price INTEGER,
    descript VARCHAR,
    quantity INTEGER,
    image BYTEA
);

CREATE TABLE IF NOT EXISTS tags (
    tagName VARCHAR PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS productTags (
    ProductId UUID NOT NULL,
    TagName VARCHAR NOT NULL,
    FOREIGN KEY (ProductId) REFERENCES products(id),
    FOREIGN KEY (TagName) REFERENCES tags(tagName),
    UNIQUE (TagName, ProductId)
);

CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY,
    UserID UUID NOT NULL,
    created_at TIMESTAMP,
    current_page INTEGER,
    current_page_search INTEGER,
    searching BOOLEAN,
    FOREIGN KEY (UserID) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS cart (
    CartId UUID PRIMARY KEY,
    SessionId UUID NOT NULL,
    ProductId UUID NOT NULL,
    Quantity INTEGER DEFAULT 1,
    FOREIGN KEY (SessionId) REFERENCES sessions(id),
    FOREIGN KEY (ProductId) REFERENCES products(id),
    UNIQUE (SessionId, ProductId)
);

