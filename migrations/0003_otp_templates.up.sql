CREATE TABLE IF NOT EXISTS otp_templates (
    id UUID PRIMARY KEY,
    client_id UUID NOT NULL, -- ideally should reference a clients table if it exists
    
    -- Array of supported channels: sms, email, whatsapp
    channel VARCHAR(50) NOT NULL,
    
    template_name VARCHAR(100) NOT NULL,  -- e.g., login_otp, signup_otp
    template_body TEXT NOT NULL,
    
    vendor VARCHAR(100),  -- e.g., msg91, twilio, sendgrid
    email_subject VARCHAR(255), -- Useful for email templates
    
    -- Vendor-specific template identifier (e.g., MSG91, WhatsApp)
    vendor_template_id VARCHAR(255),
    
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO otp_templates (
    id,
    client_id,
    channel,
    template_name,
    template_body,
    vendor_template_id,
    vendor,
    email_subject,
    is_active
) VALUES (
    '660e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440000',
    'email',
    'login_otp',
    'Hello,\n
    Your {{.otp_type}} OTP is {{.otp}}. It is valid for {{.expiry_minutes}} minutes.\n
    If you did not request this OTP, please ignore this email.\n
    Thank you,\n
    Team',
    NULL,
    'sendgrid',
    'Your Login OTP',
    TRUE
) On CONFLICT (id) DO NOTHING;

INSERT INTO otp_templates (
    id,
    client_id,
    channel,
    template_name,
    template_body,
    vendor_template_id,
    vendor,
    is_active
) VALUES (
    '39466665-bcc5-4054-936b-f6c171464de6',
    '550e8400-e29b-41d4-a716-446655440000',
    'sms',
    'login_otp',
    'Your {{.otp_type}} OTP is {{.otp}}. Valid for {{.expiry_minutes}} minutes. Do not share it with anyone.',
    NULL,
    'msg91',
    TRUE
) On CONFLICT (id) DO NOTHING;


INSERT INTO otp_templates (
    id,
    client_id,
    channel,
    template_name,
    template_body,
    vendor_template_id,
    vendor,
    is_active
) VALUES (
    'ab12c868-6fe4-4b13-b02c-c44d168fee61',
    '770e8400-e29b-41d4-a716-446655440003',
    'sms',
    'transaction_otp',
    'Your {{.otp_type}} OTP is {{.otp}}. It is valid for {{.expiry_minutes}} minutes. Please do not share this code with anyone.',
    NULL,
    'twilio',
    TRUE
) On CONFLICT (id) DO NOTHING;


INSERT INTO otp_templates (
    id,
    client_id,
    channel,
    template_name,
    template_body,
    vendor_template_id,
    vendor,
    email_subject,
    is_active
) VALUES (
    '9d6b6bd1-0f51-43a7-8048-81f8de0a31c5',
    '550e8400-e29b-41d4-a716-446655440000',
    'email',
    'signup_otp',
    'Hello,\n
    Your {{.otp_type}} OTP is {{.otp}}. It is valid for {{.expiry_minutes}} minutes.\n
    If you did not request this OTP, please ignore this email.\n
    Thank you,\n
    Team',
    NULL,
    'sendgrid',
    'Your Login OTP',
    TRUE
) On CONFLICT (id) DO NOTHING;

