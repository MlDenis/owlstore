BEGIN TRANSACTION;

DO $$
BEGIN
   
    CREATE TABLE IF NOT EXISTS balance (
            id INT GENERATED ALWAYS AS IDENTITY,
            userlogin TEXT NOT NULL,
            sumaccrual BIGINT,
            sumwithdraw BIGINT,
            PRIMARY KEY(id),
            UNIQUE(userlogin)
    );

    CREATE INDEX IF NOT EXISTS userlogin_id ON users USING hash(userlogin);
END $$;
--
--
COMMIT TRANSACTION;