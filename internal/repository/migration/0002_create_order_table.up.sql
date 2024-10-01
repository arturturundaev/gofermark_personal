CREATE TABLE IF NOT EXISTS orders
(
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    "number" character varying NOT NULL,
    status character varying NOT NULL,
    accrual float,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
)