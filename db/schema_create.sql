-- #############################################################################
-- ### UNIFIED DATABASE SCHEMA (CREATE STATEMENTS)                             ###
-- #############################################################################
-- This file contains the complete and authoritative CREATE statements for the
-- entire application, merging 'auth_microservice' and 'controller_microservice_v2'.
--
-- Conventions:
--   - 'users' table columns follow 'auth_microservice' conventions for compatibility.
--   - Other tables use snake_case.
-- #############################################################################


-- =============================================================================
-- CORE TABLES (SHARED)
-- =============================================================================

-- Merged from both services, represents the authoritative user model.
-- Column names like 'userName', 'fullName', 'accountStatus' retained for
-- backward compatibility with 'auth_microservice'.
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    userName TEXT UNIQUE NOT NULL,
    fullName TEXT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role TEXT DEFAULT 'user',
    accountStatus TEXT DEFAULT 'active',
    createdAt TIMESTAMPTZ DEFAULT now() NOT NULL,
    updatedAt TIMESTAMPTZ DEFAULT now() NOT NULL
);


-- =============================================================================
-- AUTH MICROSERVICE & V1 TABLES
-- =============================================================================

CREATE TABLE IF NOT EXISTS registerOtp (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    otp TEXT NOT NULL,
    createdAt TIMESTAMPTZ DEFAULT now() NOT NULL,
    updatedAt TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS run (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    status TEXT DEFAULT 'scheduled', -- 'scheduled', 'running', 'completed', 'failed'
    type TEXT NOT NULL, -- 'ea', 'gp', 'ml', 'pso'
    command TEXT NOT NULL,
    createdBy UUID REFERENCES users(id),
    createdAt TIMESTAMPTZ DEFAULT now() NOT NULL,
    updatedAt TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS access (
    runID UUID REFERENCES run(id),
    userID UUID REFERENCES users(id),
    mode TEXT DEFAULT 'read', -- 'read', 'write'
    PRIMARY KEY (runID, userID),
    createdAt TIMESTAMPTZ DEFAULT now() NOT NULL,
    updatedAt TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS team (
    teamID UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    teamName TEXT UNIQUE NOT NULL,
    teamDesc TEXT,
    createdBy TEXT NOT NULL,
    createdAt TIMESTAMPTZ DEFAULT now() NOT NULL,
    updatedAt TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS teamMembers (
    memberId UUID REFERENCES users(id),
    teamID UUID REFERENCES team(teamID),
    role TEXT,
    PRIMARY KEY (memberId, teamID),
    createdAt TIMESTAMPTZ DEFAULT now() NOT NULL,
    updatedAt TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS password_reset_otps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    otp_code TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    is_used BOOLEAN DEFAULT FALSE
);


-- =============================================================================
-- CONTROLLER MICROSERVICE V2 TABLES
-- =============================================================================

CREATE TABLE IF NOT EXISTS problem_statements (
  id UUID PRIMARY KEY,
  title TEXT NOT NULL,
  description_json JSONB NOT NULL,
  created_by UUID REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS notebooks (
  id UUID PRIMARY KEY,
  title TEXT NOT NULL,
  context_minio_url TEXT,
  problem_statement_id UUID REFERENCES problem_statements(id) ON DELETE CASCADE,
  requirements TEXT,
  created_at TIMESTAMPTZ NOT NULL,
  last_modified_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
  id UUID PRIMARY KEY,
  notebook_id UUID REFERENCES notebooks(id) ON DELETE CASCADE,
  current_kernel_id UUID,
  status TEXT NOT NULL,
  last_active_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS cells (
  id UUID PRIMARY KEY,
  notebook_id UUID REFERENCES notebooks(id) ON DELETE CASCADE,
  cell_index INT NOT NULL,
  cell_name TEXT,
  cell_type TEXT NOT NULL CHECK (cell_type IN ('code', 'markdown', 'raw')),
  source TEXT NOT NULL,
  execution_count INT
);

CREATE TABLE IF NOT EXISTS cell_outputs (
  id UUID PRIMARY KEY,
  cell_id UUID REFERENCES cells(id) ON DELETE CASCADE,
  output_index INT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('stream', 'display_data', 'execute_result', 'error')),
  data_json JSONB,
  minio_url TEXT,
  execution_count INT
);

CREATE TABLE IF NOT EXISTS evolution_runs (
  id UUID PRIMARY KEY,
  source_cell_id UUID REFERENCES cells(id) ON DELETE CASCADE,
  start_time TIMESTAMPTZ NOT NULL,
  end_time TIMESTAMPTZ,
  status TEXT NOT NULL
);

CREATE TABLE cell_variations (
  id UUID PRIMARY KEY,
  evolution_run_id UUID REFERENCES evolution_runs(id) ON DELETE CASCADE,
  code TEXT NOT NULL,
  metric FLOAT NOT NULL,
  is_best BOOLEAN NOT NULL,
  generation INT NOT NULL,
  parent_variant_id UUID REFERENCES cell_variations(id) ON DELETE SET NULL
);


-- =============================================================================
-- INDEXES
-- =============================================================================

CREATE INDEX IF NOT EXISTS idx_password_reset_user_id ON password_reset_otps(user_id);
