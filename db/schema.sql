

-- Drop tables in reverse order of creation to handle dependencies
DROP TABLE IF EXISTS cell_variations CASCADE;
DROP TABLE IF EXISTS evolution_runs CASCADE;
DROP TABLE IF EXISTS cell_outputs CASCADE;
DROP TABLE IF EXISTS cells CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS notebooks CASCADE;
DROP TABLE IF EXISTS problem_statements CASCADE;
DROP TABLE IF EXISTS users CASCADE;


CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  username TEXT NOT NULL,
  email TEXT NOT NULL,
  password TEXT NOT NULL,
  role TEXT NOT NULL,
  acc_status TEXT NOT NULL
);


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

CREATE TABLE IF NOT EXISTS cell_variations (
  id UUID PRIMARY KEY,
  evolution_run_id UUID REFERENCES evolution_runs(id) ON DELETE CASCADE,
  code TEXT NOT NULL,
  metric FLOAT NOT NULL,
  is_best BOOLEAN NOT NULL,
  generation INT NOT NULL,
  parent_variant_id UUID REFERENCES cell_variations(id) ON DELETE SET NULL
);

-- Remove this in productions, just a dummy data for testing before integrating the auth service
DELETE FROM users;

INSERT INTO users (
  id, username, email, password, role, acc_status
) VALUES (
  '123e4567-e89b-12d3-a456-426614174000',
  'Tharun Kumarr A',
  'tharunkumarra@gmail.com',
  'HelloThere',
  'Student',
  'valid'
);
