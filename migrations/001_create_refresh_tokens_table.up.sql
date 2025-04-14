CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    token_hash TEXT NOT NULL,
    associated_jti TEXT NOT NULL,
    client_ip TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_id ON refresh_tokens(user_id);