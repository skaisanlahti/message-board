BEGIN;
DO $$ 
BEGIN 
    IF NOT EXISTS(SELECT 1 FROM migrations WHERE version = 1) THEN 
        RAISE NOTICE 'migration add_users_table not applied, skipping';
        RETURN;
    END IF;

    DROP INDEX IF EXISTS index_user_name;
    DROP TABLE IF EXISTS users;

    DELETE FROM migrations WHERE version = 1;
END $$;
COMMIT;
