CREATE TABLE IF NOT EXISTS otp_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL,
    user_id UUID NOT NULL,
    otp_type VARCHAR(50) NOT NULL,
    identifier VARCHAR(255) NOT NULL,          -- phone number or email
    channel VARCHAR(20) NOT NULL,              -- sms / email / whatsapp
    event_type VARCHAR(20) NOT NULL,           -- GENERATED / SENT / VERIFIED / FAILED / RESENT
    status VARCHAR(10) NOT NULL,               -- success / failure
    attempt_number INT DEFAULT 0,              -- increments with each verification attempt
    resend_number INT DEFAULT 0,
    vendor VARCHAR(50),                        -- e.g., MSG91, Twilio
    otp VARCHAR(10),                           -- Store OTP 
    submitted_otp VARCHAR(10),                 -- OTP submitted by user during verification attempts
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE otp_events ADD COLUMN error TEXT;
-- CAN BE ADDED LATER IF NEEDED

-- CONSTRAINT chk_event_type CHECK (
--         event_type IN ('GENERATED', 'SENT', 'VERIFIED', 'FAILED', 'RESENT')
--     ),
--     CONSTRAINT chk_status CHECK (
--         status IN ('success', 'failure')
--     ),
--     CONSTRAINT chk_channel CHECK (
--         channel IN ('sms', 'email', 'whatsapp')
--     )