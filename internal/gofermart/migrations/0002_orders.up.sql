BEGIN TRANSACTION;

DO $$
BEGIN
       
    CREATE TABLE IF NOT EXISTS orders (
            id INT GENERATED ALWAYS AS IDENTITY,
            ordernumber BIGINT NOT NULL,
            userlogin TEXT NOT NULL,
            orderdate TIMESTAMP NOT NULL,  
            statusorder TEXT NOT NULL,
            accrual BIGINT,
            withdraw BIGINT,
            PRIMARY KEY(id),
            UNIQUE(ordernumber)
    );

    CREATE INDEX IF NOT EXISTS orderdate_id ON orders USING hash(orderdate);
END $$;

--
--
COMMIT TRANSACTION;