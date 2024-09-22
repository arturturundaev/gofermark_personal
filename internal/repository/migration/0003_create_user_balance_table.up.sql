CREATE TABLE IF NOT EXISTS user_balance
(
    id uuid NOT NULL,
    user_id uuid  UNIQUE NOT NULL,
    sum float,
    with_drawn float
)