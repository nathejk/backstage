CREATE TABLE IF NOT EXISTS payment (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    ts text NOT NULL,
    shopNumber text NOT NULL,
    amount integer NOT NULL,
    currency text NOT NULL,
    message text NOT NULL,
    userPhoneNumber text NOT NULL,
    userName text NOT NULL,
    version integer NOT NULL DEFAULT 1
);
