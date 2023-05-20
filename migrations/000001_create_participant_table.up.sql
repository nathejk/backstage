CREATE TABLE IF NOT EXISTS participant (
    id bigserial PRIMARY KEY,
    uuid text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    address text NOT NULL,
    email text NOT NULL,
    phone text NOT NULL,
    team text NOT NULL,
    days text[] NOT NULL,
    transport text NOT NULL,
    seatCount text NOT NULL,
    info text NOT NULL,
    video text NOT NULL,
    version integer NOT NULL DEFAULT 1
);
