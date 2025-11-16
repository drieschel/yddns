package cache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const (
	ExpirySecondsIndefinite = 0
	ExpirySecondsDefault    = 600
)

type Cache interface {
	Delete(key string) error
	Get(key string) (*Item, error)
	Set(item *Item) error
}

type Item struct {
	Key           string      `json:"key"`
	Value         interface{} `json:"value"`
	ExpirySeconds int         `json:"expiry_seconds"`
	CreatedAt     *time.Time  `json:"created_at"`
	ModifiedAt    *time.Time  `json:"modified_at"`
}

func NewItem(key string, value interface{}) *Item {
	return NewItemWithCustomExpiry(key, value, -1)
}

func NewItemWithCustomExpiry(key string, value interface{}, expirySeconds int) *Item {
	now := time.Now()
	return &Item{CreatedAt: &now, ExpirySeconds: expirySeconds, Key: key, Value: value}
}

func (i *Item) IsValid() bool {
	if i.ModifiedAt != nil {
		return i.ExpirySeconds == ExpirySecondsIndefinite || time.Now().Sub(*i.ModifiedAt) < time.Duration(i.ExpirySeconds)*time.Second
	}

	return false
}

type FileCache struct {
	cacheDir      string
	expirySeconds int
}

func NewFileCacheWithDefaultExpiry(cacheDir string) *FileCache {
	return NewFileCache(cacheDir, ExpirySecondsDefault)
}

func NewFileCache(cacheDir string, expirySeconds int) *FileCache {
	return &FileCache{cacheDir: cacheDir, expirySeconds: expirySeconds}
}

func (c *FileCache) Delete(key string) error {
	fileName := c.createFilePath(key)
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	err = os.Remove(fileName)
	if err != nil {
		return err
	}

	return nil
}

func (c *FileCache) Get(key string) (*Item, error) {
	fileName := c.createFilePath(key)
	data, err := os.ReadFile(fileName)
	if err != nil {
		return NewItem(key, nil), nil
	}

	var item Item
	err = json.Unmarshal(data, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (c *FileCache) Set(item *Item) error {
	if item.ExpirySeconds < 0 {
		item.ExpirySeconds = c.expirySeconds
	}

	now := time.Now()
	item.ModifiedAt = &now

	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	filePath := c.createFilePath(item.Key)

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (c *FileCache) createFilePath(key string) string {
	return filepath.Join(c.cacheDir, CreateHashedKey(key))
}

func CreateHashedKey(key string) string {
	hash := md5.New()
	hash.Write([]byte(key))

	return hex.EncodeToString(hash.Sum(nil))
}
