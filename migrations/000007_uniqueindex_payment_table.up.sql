ALTER TABLE "payment"
    DROP ts,
    ADD ts timestamp(0) with time zone NOT NULL,
    ADD CONSTRAINT "payment_shopnumber_ts" UNIQUE ("shopnumber", "ts");

