CREATE TABLE IF NOT EXISTS user_balance
(
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    sum money,
    withDrawn money
)