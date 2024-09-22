CREATE TABLE IF NOT EXISTS user_withdrawals
(
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    sum float,
    "number" character varying NOT NULL,
    created_at timestamp
)
