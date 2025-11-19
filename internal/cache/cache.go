package cache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	ExpirySecondsIndefinite      = -1
	ExpirySecondsUseGlobal       = -2
	CreatedExpirySecondsDefault  = 86400
	ModifiedExpirySecondsDefault = 600
)

type Cache interface {
	Get(key string) (*Item, error)
	Set(item *Item) error
	Delete(key string) error
	DeleteExpired() error
	IsValid(item Item) bool
}

type Item struct {
	Key            string      `json:"key"`
	Value          interface{} `json:"value"`
	Created        *time.Time  `json:"created"`
	CreatedExpiry  int         `json:"created_expiry"`
	Modified       *time.Time  `json:"modified"`
	ModifiedExpiry int         `json:"modified_expiry"`
}

func NewItem(key string, value interface{}) *Item {
	return NewItemWithCustomExpiry(key, value, ExpirySecondsUseGlobal, ExpirySecondsUseGlobal)
}

func NewItemWithCustomExpiry(key string, value interface{}, createdExpiry int, modifiedExpiry int) *Item {
	now := time.Now()
	return &Item{Key: key, Value: value, Created: &now, CreatedExpiry: createdExpiry, ModifiedExpiry: modifiedExpiry}
}

type FileCache struct {
	cacheDir       string
	createdExpiry  int
	modifiedExpiry int
}

func NewFileCacheWithDefaultValues(cacheDir string) *FileCache {
	return NewFileCache(cacheDir, CreatedExpirySecondsDefault, ModifiedExpirySecondsDefault)
}

func NewFileCache(cacheDir string, createdExpiry int, modifiedExpiry int) *FileCache {
	return &FileCache{cacheDir: cacheDir, createdExpiry: createdExpiry, modifiedExpiry: modifiedExpiry}
}

func (c *FileCache) Get(key string) (*Item, error) {
	fileName := c.createFilePath(key)
	item, err := c.GetByFileName(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return NewItem(key, nil), nil
		}

		return nil, err
	}

	return item, nil
}

func (c *FileCache) GetByFileName(fileName string) (*Item, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var item Item
	err = json.Unmarshal(data, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (c *FileCache) Set(item *Item) error {
	now := time.Now()
	item.Modified = &now

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

func (c *FileCache) DeleteExpired() error {
	files, _ := filepath.Glob(filepath.Join(c.cacheDir, "*.json"))
	var errs []error
	for _, fileName := range files {
		item, err := c.GetByFileName(fileName)
		if err != nil {
			errs = append(errs, err)
		} else {
			if item.CreatedExpiry == ExpirySecondsUseGlobal {
				item.CreatedExpiry = c.createdExpiry
			}

			now := time.Now()
			if item.CreatedExpiry > ExpirySecondsIndefinite && now.Sub(*item.Created) > time.Duration(item.CreatedExpiry)*time.Second {
				err = os.Remove(fileName)
				if err != nil {
					errs = append(errs, err)
				}
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (c *FileCache) IsValid(item Item) bool {
	if item.Modified != nil {
		if item.ModifiedExpiry == ExpirySecondsUseGlobal {
			item.ModifiedExpiry = c.modifiedExpiry
		}

		return item.ModifiedExpiry == ExpirySecondsIndefinite || time.Now().Sub(*item.Modified) < time.Duration(item.ModifiedExpiry)*time.Second
	}

	return false
}

func (c *FileCache) createFilePath(key string) string {
	return fmt.Sprintf("%s.json", filepath.Join(c.cacheDir, CreateHashedKey(key)))
}

func CreateHashedKey(key string) string {
	hash := md5.New()
	hash.Write([]byte(key))

	return hex.EncodeToString(hash.Sum(nil))
}
