CREATE TABLE IF NOT EXISTS otp_configs (
    client_id UUID NOT NULL,
    otp_type VARCHAR(50) NOT NULL,            -- e.g. login, signup, transaction
    channel VARCHAR(50) NOT NULL,             -- sms, email, whatsapp

    reuse_otp_on_resend BOOLEAN NOT NULL DEFAULT FALSE,

    template_id UUID,

    otp_length INT NOT NULL DEFAULT 6,
    expiry_seconds INT NOT NULL DEFAULT 300,

    max_verify_attempts INT,                 

    resend_allowed BOOLEAN NOT NULL DEFAULT TRUE,
    max_resend_attempts INT,

    cooldown_seconds INT,

    vendor VARCHAR(100) NOT NULL,             -- e.g. msg91, twilio

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (client_id, otp_type, channel)
);


-- sample queries to insert data into otp_configs table
INSERT INTO otp_configs (client_id, otp_type, reuse_otp_on_resend, otp_length, expiry_seconds, max_verify_attempts, 
resend_allowed, max_resend_attempts, cooldown_seconds, template_name) VALUES 
('550e8400-e29b-41d4-a716-446655440000', 'login', TRUE, 6, 300, 3, TRUE, 2, 60, 'login_otp') ON CONFLICT DO NOTHING;

INSERT INTO otp_configs (client_id, otp_type, reuse_otp_on_resend, otp_length, expiry_seconds, max_verify_attempts, 
resend_allowed, max_resend_attempts, cooldown_seconds, template_name) VALUES 
('550e8400-e29b-41d4-a716-446655440000', 'signup', FALSE, 4, 600, 5, TRUE, 3, NULL, 'signup_otp') ON CONFLICT DO NOTHING;

INSERT INTO otp_configs (client_id, otp_type, reuse_otp_on_resend, otp_length, expiry_seconds, max_verify_attempts, 
resend_allowed, max_resend_attempts, cooldown_seconds, template_name) VALUES 
('770e8400-e29b-41d4-a716-446655440003', 'transaction', TRUE, 4, 200, 2, FALSE, 0, NULL, 'transaction_otp') ON CONFLICT DO NOTHING;



