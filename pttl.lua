-- pttl.lua: => Arguments: [value]

if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("pttl", KEYS[1])
else
	return -3
end
