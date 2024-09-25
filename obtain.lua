-- obtain.lua: arguments => [value, tokenLen, ttl]

if redis.call("set", KEYS[1], ARGV[1], "NX", "PX", ARGV[3]) then
	return redis.status_reply("OK")
end

local offset = tonumber(ARGV[2])
if redis.call("getrange", KEYS[1], 0, offset-1) == string.sub(ARGV[1], 1, offset) then
	return redis.call("set", KEYS[1], ARGV[1], "PX", ARGV[3])
end
