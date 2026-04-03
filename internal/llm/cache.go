package llm

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

// Cache provides a simple file-based cache for LLM responses.
type Cache struct {
	dir string
}

func NewCache(dir string) *Cache {
	return &Cache{dir: dir}
}

func (c *Cache) Get(prompt string) (string, bool) {
	data, err := os.ReadFile(c.keyPath(prompt))
	if err != nil {
		return "", false
	}
	return string(data), true
}

func (c *Cache) Set(prompt, response string) error {
	if err := os.MkdirAll(c.dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(c.keyPath(prompt), []byte(response), 0o644)
}

func (c *Cache) keyPath(prompt string) string {
	hash := sha256.Sum256([]byte(prompt))
	return filepath.Join(c.dir, hex.EncodeToString(hash[:])+".txt")
}
