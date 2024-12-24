-- Users Table
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    full_name VARCHAR(100),
    email VARCHAR(100) UNIQUE NOT NULL,
    phone_number VARCHAR(15) UNIQUE,
    password_hash TEXT NOT NULL,
    profile_picture TEXT,
    status_message VARCHAR(255) DEFAULT 'Hey there! I am using the app.',
    online_status BOOLEAN DEFAULT FALSE,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Privacy and Notification Settings Table
CREATE TABLE user_settings (
    setting_id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    privacy_status VARCHAR(20) DEFAULT 'Everyone',
    -- Options: 'Everyone', 'Contacts', 'Nobody'
    notification_preference JSONB DEFAULT '{"push": true, "email": true, "sms": false}',
    language_preference VARCHAR(20) DEFAULT 'en',
    theme_preference VARCHAR(10) DEFAULT 'light',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Contacts Table
CREATE TABLE contacts (
    contact_id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    contact_user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    is_blocked BOOLEAN DEFAULT FALSE,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Messages Table
CREATE TABLE messages (
    message_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    receiver_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_read BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE
);
-- Chat Threads Table
CREATE TABLE chat_threads (
    thread_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user1_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    user2_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    last_message_id UUID REFERENCES messages(message_id),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Devices Table
CREATE TABLE user_devices (
    device_id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    device_name VARCHAR(100) NOT NULL,
    device_token TEXT NOT NULL,
    -- For notifications
    last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Activity Logs Table
CREATE TABLE activity_logs (
    log_id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    activity_type VARCHAR(50),
    -- e.g., 'login', 'logout', 'password_change'
    activity_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address INET
);
-- Encryption Keys Table (Optional for End-to-End Encryption)
CREATE TABLE encryption_keys (
    key_id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    public_key TEXT NOT NULL,
    private_key TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Posts Table
CREATE TABLE shares (
    share_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    -- The main content of the post
    media_url TEXT,
    -- Optional: URL for any attached media (image, video, etc.)
    visibility VARCHAR(20) DEFAULT 'Friends',
    -- Options: 'Public', 'Friends', 'Private'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- share_reactions Table
CREATE TABLE share_reactions (
    reaction_id SERIAL PRIMARY KEY,
    share_id UUID REFERENCES shares(share_id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    reaction_type VARCHAR(20) DEFAULT 'like',
    -- Options: 'like', 'love', 'laugh', etc.
    reacted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- share_comments Table
CREATE TABLE share_comments (
    comment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    share_id UUID REFERENCES shares(share_id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    comment_text TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- share_visibility Table
CREATE TABLE share_visibility (
    visibility_id SERIAL PRIMARY KEY,
    share_id UUID REFERENCES shares(share_id) ON DELETE CASCADE,
    visible_to_user_id UUID REFERENCES users(user_id) ON DELETE CASCADE
);
-- Unique indexes
CREATE UNIQUE INDEX idx_users_email ON users (email);
CREATE UNIQUE INDEX idx_users_username ON users (username);
-- Common query optimization
CREATE INDEX idx_users_last_seen ON users (last_seen);
CREATE INDEX idx_users_online_status ON users (online_status);
-- Foreign key index
CREATE INDEX idx_user_settings_user_id ON user_settings (user_id);
-- Foreign key indexes
CREATE INDEX idx_contacts_user_id ON contacts (user_id);
CREATE INDEX idx_contacts_contact_user_id ON contacts (contact_user_id);
-- Common query optimization
CREATE INDEX idx_contacts_is_blocked ON contacts (is_blocked);
-- Foreign key indexes
CREATE INDEX idx_messages_sender_id ON messages (sender_id);
CREATE INDEX idx_messages_receiver_id ON messages (receiver_id);
-- Common query optimization
CREATE INDEX idx_messages_created_at ON messages (created_at);
CREATE INDEX idx_messages_is_read ON messages (is_read);
-- Foreign key indexes
CREATE INDEX idx_chat_threads_user1_id ON chat_threads (user1_id);
CREATE INDEX idx_chat_threads_user2_id ON chat_threads (user2_id);
-- Optimize queries for last updated threads
CREATE INDEX idx_chat_threads_last_updated ON chat_threads (last_updated);
-- Foreign key index
CREATE INDEX idx_user_devices_user_id ON user_devices (user_id);
-- Optimize login-related queries
CREATE INDEX idx_user_devices_last_login ON user_devices (last_login);
-- Foreign key index
CREATE INDEX idx_activity_logs_user_id ON activity_logs (user_id);
-- Optimize queries for recent activity
CREATE INDEX idx_activity_logs_activity_timestamp ON activity_logs (activity_timestamp);
-- Foreign key index
CREATE INDEX idx_shares_user_id ON shares (user_id);
-- Common query optimization
CREATE INDEX idx_shares_created_at ON shares (created_at);
CREATE INDEX idx_shares_visibility ON shares (visibility);
-- Foreign key indexes
CREATE INDEX idx_share_reactions_share_id ON share_reactions (share_id);
CREATE INDEX idx_share_reactions_user_id ON share_reactions (user_id);
-- Common query optimization
CREATE INDEX idx_share_reactions_reaction_type ON share_reactions (reaction_type);
-- Foreign key indexes
CREATE INDEX idx_share_comments_share_id ON share_comments (share_id);
CREATE INDEX idx_share_comments_user_id ON share_comments (user_id);
-- Optimize queries for recent comments
CREATE INDEX idx_share_comments_created_at ON share_comments (created_at);
-- Foreign key indexes
CREATE INDEX idx_share_visibility_share_id ON share_visibility (share_id);
CREATE INDEX idx_share_visibility_visible_to_user_id ON share_visibility (visible_to_user_id);