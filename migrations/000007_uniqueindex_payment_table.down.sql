ALTER TABLE "payment"
    DROP ts,
    ADD ts text NOT NULL,
    DROP CONSTRAINT "payment_shopnumber_ts";
