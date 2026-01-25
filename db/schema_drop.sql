-- This file contains all DROP statements for the unified schema.
-- It is executed first to ensure a clean slate before creating tables.
-- The order respects foreign key constraints.

DROP TABLE IF EXISTS cell_variations;
DROP TABLE IF EXISTS cell_outputs;
DROP TABLE IF EXISTS evolution_runs;
DROP TABLE IF EXISTS cells;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS notebooks;
DROP TABLE IF EXISTS problem_statements;
DROP TABLE IF EXISTS password_reset_otps;
DROP TABLE IF EXISTS teamMembers;
DROP TABLE IF EXISTS access;
DROP TABLE IF EXISTS run;
DROP TABLE IF EXISTS registerOtp;
DROP TABLE IF EXISTS team;
DROP TABLE IF EXISTS users;
