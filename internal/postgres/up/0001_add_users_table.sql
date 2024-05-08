BEGIN;
DO $$ 
BEGIN 
    IF EXISTS(SELECT 1 FROM migrations WHERE version = 1) THEN
        RAISE NOTICE 'migration add_users_table already applied, skipping';
        RETURN;
    END IF;

    CREATE TABLE IF NOT EXISTS users(
        id INTEGER NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
        name TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );

    CREATE UNIQUE INDEX IF NOT EXISTS index_users_name ON users(name);

    INSERT INTO migrations (version, name) VALUES (1, 'add_users_table');
END $$;
COMMIT;
