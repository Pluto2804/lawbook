-- Drop tables in reverse order due to foreign key constraints

DROP TABLE IF EXISTS performance_evaluations;
DROP TABLE IF EXISTS session_participants;
DROP TABLE IF EXISTS moot_sessions;
DROP TABLE IF EXISTS recruiter_profiles;
DROP TABLE IF EXISTS lawyer_profiles;
DROP TABLE IF EXISTS student_profiles;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

DROP DATABASE IF EXISTS lawbookauth;
