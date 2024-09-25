package redislock

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"strings"
)

// RedisClient is a minimal client interface.
type RedisClient interface {
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
	EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) (interface{}, error)
}

type Script struct {
	src, hash string
}

func NewScript(src string) *Script {
	h := sha1.New()
	_, _ = io.WriteString(h, src)
	return &Script{
		src:  src,
		hash: hex.EncodeToString(h.Sum(nil)),
	}
}

func (s *Script) Hash() string {
	return s.hash
}

func (s *Script) Eval(ctx context.Context, c RedisClient, keys []string, args ...interface{}) (interface{}, error) {
	return c.Eval(ctx, s.src, keys, args...)
}

func (s *Script) EvalSha(ctx context.Context, c RedisClient, keys []string, args ...interface{}) (interface{}, error) {
	return c.EvalSha(ctx, s.hash, keys, args...)
}

// Run optimistically uses EVALSHA to run the script. If script does not exist
// it is retried using EVAL.
func (s *Script) Run(ctx context.Context, c RedisClient, keys []string, args ...interface{}) (interface{}, error) {
	res, err := s.EvalSha(ctx, c, keys, args...)
	if IsRedisNoScript(err) {
		return s.Eval(ctx, c, keys, args...)
	}
	return res, err
}

func IsRedisNil(err error) bool {
	return err != nil && err.Error() == "redis: nil"
}

func IsRedisNoScript(err error) bool {
	return err != nil && strings.Contains(err.Error(), "NOSCRIPT")
}
