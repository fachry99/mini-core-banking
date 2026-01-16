CREATE TABLE IF NOT EXISTS transfer_limits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL,
    daily_limit BIGINT NOT NULL,
    used_today BIGINT DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_account
        FOREIGN KEY(account_id)
        REFERENCES accounts(id)
);
