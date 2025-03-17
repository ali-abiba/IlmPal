-- Schema for IlmPal Reading Application

-- Delete tables if they exist (for clean recreation)
DROP TABLE IF EXISTS reading;
DROP TABLE IF EXISTS top_categories;
DROP TABLE IF EXISTS book_category;
DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS book;
DROP TABLE IF EXISTS "user";

-- Create book table
CREATE TABLE book (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    author VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create category table
CREATE TABLE category (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create book_category junction table
CREATE TABLE book_category (
    id SERIAL PRIMARY KEY,
    book_id INTEGER NOT NULL REFERENCES book(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES category(id) ON DELETE CASCADE,
    weight FLOAT DEFAULT 1.0, -- Weight to indicate relevance of category to book
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(book_id, category_id) -- Prevent duplicate book-category associations
);

-- Create user table
CREATE TABLE "user" (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create top_categories table for user category preferences
CREATE TABLE top_categories (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES category(id) ON DELETE CASCADE,
    weight FLOAT DEFAULT 1.0, -- Weight to indicate user's interest in category
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, category_id) -- Prevent duplicate user-category associations
);

-- Create reading table for tracking user reading progress
CREATE TABLE reading (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    book_id INTEGER NOT NULL REFERENCES book(id) ON DELETE CASCADE,
    notes TEXT,
    started_reading TIMESTAMP WITH TIME ZONE,
    last_reading TIMESTAMP WITH TIME ZONE,
    bookmark_page INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, book_id) -- A user can have only one reading record per book
);

-- Create indexes for performance
CREATE INDEX idx_book_category_book_id ON book_category(book_id);
CREATE INDEX idx_book_category_category_id ON book_category(category_id);
CREATE INDEX idx_top_categories_user_id ON top_categories(user_id);
CREATE INDEX idx_top_categories_category_id ON top_categories(category_id);
CREATE INDEX idx_reading_user_id ON reading(user_id);
CREATE INDEX idx_reading_book_id ON reading(book_id);

-- Trigger function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers to automatically update updated_at columns
CREATE TRIGGER update_book_modtime
    BEFORE UPDATE ON book
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_user_modtime
    BEFORE UPDATE ON "user"
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_top_categories_modtime
    BEFORE UPDATE ON top_categories
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_reading_modtime
    BEFORE UPDATE ON reading
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();
