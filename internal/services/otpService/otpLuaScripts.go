package otpservice

var verifyOtpLuaScript = `
-- otpVerifyLua.lua
-- KEYS[1] = Redis key for OTP
-- ARGV[1] = submitted OTP

local otpKey = KEYS[1]
local submittedOTP = ARGV[1]

-- Fetch all fields
local data = redis.call("HGETALL", otpKey)
if not data or #data == 0 then
    return -2  -- OTP expired or does not exist
end

local otp = nil
local attempts = 0
local maxVerifyAttempts = 0

-- Convert HGETALL array to map
for i = 1, #data, 2 do
    local field = data[i]
    local value = data[i+1]

    if field == "otp" then
        otp = value
    elseif field == "verify_attempts" then
        attempts = tonumber(value) or 0
    elseif field == "max_verify_attempts" then
        maxVerifyAttempts = tonumber(value) or 0
    end
end

-- Safety check (in case data is corrupted)
if not otp then
    return -3  -- Invalid state (OTP missing)
end

-- Check attempt limit
if maxVerifyAttempts > 0 and attempts >= maxVerifyAttempts then
    return -1  -- Too many attempts
end

-- Verify OTP
if submittedOTP == otp then
    redis.call("DEL", otpKey)
    return 1  -- Success
else
    redis.call("HINCRBY", otpKey, "verify_attempts", 1)
    return 0  -- Invalid OTP
end
`

var resendOtpLuaScript = `
-- KEYS[1] = Redis key for OTP
-- ARGV[1] = current timestamp (in seconds)

local key = KEYS[1]
local now = tonumber(ARGV[1])

-- Fetch all fields from hash
local data = redis.call("HGETALL", key)
if not data or #data == 0 then
    return -2  -- OTP expired or does not exist
end

-- Convert HGETALL array to table
local otpData = {}
for i = 1, #data, 2 do
    otpData[data[i]] = data[i+1]
end

-- Parse important fields with safe defaults
local resendAllowed      = tonumber(otpData["resend_allowed"] or 0)
local resendAttempts     = tonumber(otpData["resend_attempts"] or 0)
local maxResendAttempts  = tonumber(otpData["max_resend_attempts"] or 0)
local lastSent           = tonumber(otpData["last_sent"] or 0)
local cooldown           = tonumber(otpData["cooldown_seconds"] or 0)
local reuseOtp           = tonumber(otpData["reuse_otp_on_resend"] or 1)

-- Check if resends are allowed
if resendAllowed == 0 then
    return -3  -- Resend not allowed
end

-- Check max resend attempts
if maxResendAttempts ~= -1 and resendAttempts >= maxResendAttempts then
    return -1  -- Too many resends
end

-- Check cooldown
if cooldown > 0 and (now - lastSent) < cooldown then
    return -4  -- Still in cooldown
end

-- Determine OTP value
if reuseOtp == 1 then
    return "__REUSE__"
else
    return "__GENERATE__"
end

-- Increment resend attempts and update last_sent atomically
redis.call("HINCRBY", key, "resend_attempts", 1)
redis.call("HSET", key, "last_sent", now)

-- Return the OTP to caller
return otp
`

var resendPostLuaScript = `
-- KEYS[1] = otp key
-- ARGV[1] = otp to write
-- ARGV[2] = ttl in seconds
-- ARGV[3] = current timestamp

local key = KEYS[1]
local otp = ARGV[1]
local ttl = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
-- increment resend attempts
local resend_attempts = redis.call("HINCRBY", key, "resend_attempts", 1)

-- update OTP, last_sent, reset verify_attempts
redis.call("HSET", key,
  "otp", otp,
  "verify_attempts", 0,
  "resend_attempts", resend_attempts,
  "last_sent", now
)

-- reset TTL
redis.call("EXPIRE", key, ttl)

return 1
`
