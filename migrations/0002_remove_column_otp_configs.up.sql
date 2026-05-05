ALTER TABLE otp_configs DROP CONSTRAINT IF EXISTS otp_configs_client_id_otp_type_channel_key;

ALTER TABLE otp_configs
DROP COLUMN channel;

ALTER TABLE otp_configs
DROP COLUMN vendor;

ALTER TABLE otp_configs DROP COLUMN template_id;

ALTER TABLE otp_configs ADD COLUMN template_name VARCHAR(50);

-- update otp_configs set template_name = 'login_otp' where client_id = '550e8400-e29b-41d4-a716-446655440000' and otp_type = 'login';
-- update otp_configs set template_name = 'signup_otp' where client_id = '550e8400-e29b-41d4-a716-446655440000' and otp_type = 'signup';
-- update otp_configs set template_name = 'transaction_otp' where client_id = '770e8400-e29b-41d4-a716-446655440003' and otp_type = 'transaction';