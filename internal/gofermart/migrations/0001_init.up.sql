BEGIN TRANSACTION;

DO $$
BEGIN
   
    CREATE TABLE IF NOT EXISTS users (
            id INT GENERATED ALWAYS AS IDENTITY,
            userlogin TEXT NOT NULL,
            hashpass TEXT NOT NULL,
            PRIMARY KEY(id),
            UNIQUE(userlogin, hashpass)
    );

    CREATE INDEX IF NOT EXISTS userlogin_id ON users USING hash(userlogin);
END $$;
--
--
COMMIT TRANSACTION;