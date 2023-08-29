BEGIN TRANSACTION;

DO $$
BEGIN
       
    CREATE TABLE IF NOT EXISTS ordersaccrual (
            id INT GENERATED ALWAYS AS IDENTITY,
            ordernumber BIGINT NOT NULL,
            statusorder TEXT NOT NULL,
            accrual BIGINT,
            goods JSONB,
            PRIMARY KEY(id),
            UNIQUE(ordernumber)
    );

    CREATE INDEX IF NOT EXISTS ordernumber_id ON ordersaccrual USING hash(ordernumber);
END $$;

--
--
COMMIT TRANSACTION;