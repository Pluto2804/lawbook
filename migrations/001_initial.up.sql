-- Create database
CREATE DATABASE IF NOT EXISTS lawbookauth CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE lawbookauth;

-- Users table with role-based authentication
CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    hashed_password CHAR(60) NOT NULL,
    role ENUM('student', 'lawyer', 'recruiter') NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    INDEX idx_email (email),
    INDEX idx_role (role)
);

-- Sessions table for managing user sessions
CREATE TABLE sessions (
    token CHAR(43) NOT NULL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expiry TIMESTAMP(6) NOT NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_expiry (expiry),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Optional: Role-specific profile tables

-- Student profiles
CREATE TABLE student_profiles (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_id INTEGER NOT NULL UNIQUE,
    university VARCHAR(255),
    year_of_study INTEGER,
    specialization VARCHAR(255),
    moot_court_experience TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Lawyer profiles
CREATE TABLE lawyer_profiles (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_id INTEGER NOT NULL UNIQUE,
    bar_registration_number VARCHAR(100),
    years_of_experience INTEGER,
    specialization VARCHAR(255),
    firm_name VARCHAR(255),
    bio TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Recruiter profiles
CREATE TABLE recruiter_profiles (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_id INTEGER NOT NULL UNIQUE,
    company_name VARCHAR(255) NOT NULL,
    position VARCHAR(255),
    company_website VARCHAR(255),
    bio TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Moot court sessions (for your virtual court feature)
CREATE TABLE moot_sessions (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    session_type ENUM('single_player', 'dual_player', 'trio') NOT NULL,
    case_type VARCHAR(100),
    difficulty_level ENUM('easy', 'medium', 'hard') NOT NULL,
    created_by INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    status ENUM('setup', 'in_progress', 'completed') NOT NULL DEFAULT 'setup',
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_created_by (created_by),
    INDEX idx_status (status)
);

-- Session participants
CREATE TABLE session_participants (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    session_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    role ENUM('judge', 'appellant_counsel', 'respondent_counsel') NOT NULL,
    is_ai BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (session_id) REFERENCES moot_sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_session_user (session_id, user_id)
);

-- Performance evaluations (for recruiter dashboard)
CREATE TABLE performance_evaluations (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    session_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    overall_score DECIMAL(5,2),
    legal_knowledge_score DECIMAL(5,2),
    argumentation_score DECIMAL(5,2),
    presentation_score DECIMAL(5,2),
    response_quality_score DECIMAL(5,2),
    ai_feedback TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES moot_sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_session_id (session_id)
);

-- Add some indexes for better query performance
CREATE INDEX idx_sessions_expiry ON sessions(expiry);
CREATE INDEX idx_users_created_at ON users(created_at);
