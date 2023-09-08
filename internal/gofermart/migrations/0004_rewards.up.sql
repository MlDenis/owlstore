BEGIN TRANSACTION;

DO $$
BEGIN
       
    CREATE TABLE IF NOT EXISTS rewards (
            id INT GENERATED ALWAYS AS IDENTITY,
            match TEXT,
            reward BIGINT,
            reward_type TEXT,
            PRIMARY KEY(id),
            UNIQUE(match)
    );

    CREATE INDEX IF NOT EXISTS match_id ON rewards USING hash(match);
END $$;

--
--
COMMIT TRANSACTION;