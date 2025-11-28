-- Drop child tables first
DROP TABLE IF EXISTS cell_outputs CASCADE;
DROP TABLE IF EXISTS evolution_runs CASCADE;
DROP TABLE IF EXISTS cell_variations CASCADE;
DROP TABLE IF EXISTS cells CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS notebooks CASCADE;

-- Finally, drop parent table
DROP TABLE IF EXISTS problem_statements CASCADE;
DROP TABLE IF EXISTS users CASCADE;
