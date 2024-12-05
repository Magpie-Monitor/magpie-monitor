package tests

import (
	"github.com/Magpie-Monitor/magpie-monitor/agent/internal/database"
)

type MockRedis struct {
	storage map[string]string
}

func NewMockRedis() database.RedisDatabase {
	return &MockRedis{storage: make(map[string]string, 0)}
}

func (m *MockRedis) Set(key, value string, ttl int) error {
	m.storage[key] = value
	return nil
}

func (m *MockRedis) Get(key string) string {
	return m.storage[key]
}

var db database.RedisDatabase = &MockRedis{}
