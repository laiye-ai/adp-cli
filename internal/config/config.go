package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	configDirName = ".adp"
	configFileName = "config.json"
	keyFileName = "key.enc"
)

func getHomeDir() string {
	// Try USERPROFILE first (Windows standard)
	if home := os.Getenv("USERPROFILE"); home != "" {
		return home
	}
	// Fallback to HOME (Unix-style)
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	// Use os/user for cross-platform reliability
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

var (
	configDir = filepath.Join(getHomeDir(), configDirName)
	configPath = filepath.Join(configDir, configFileName)
	keyPath = filepath.Join(configDir, keyFileName)
)

// Config holds all configuration
type Config struct {
	APIKey      string `json:"api_key,omitempty"`
	APIBaseURL  string `json:"api_base_url,omitempty"`
	TenantName  string `json:"tenant_name,omitempty"`
}

// Load reads configuration from file
func Load() (*Config, error) {
	ensureConfigDir()

	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")

	var cfg Config
	if err := viper.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	cfg.APIKey = viper.GetString("api_key")
	cfg.APIBaseURL = viper.GetString("api_base_url")
	cfg.TenantName = viper.GetString("tenant_name")
	if cfg.TenantName == "" {
		cfg.TenantName = "laiye"
	}

	return &cfg, nil
}

// Save writes configuration to file
func Save(cfg *Config) error {
	ensureConfigDir()

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0600)
}

// Clear removes all configuration
func Clear() error {
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.Remove(keyPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// IsConfigured checks if API key and base URL are set
func IsConfigured(cfg *Config) bool {
	return cfg.APIKey != "" && cfg.APIBaseURL != ""
}

// GetAPIKey gets the decrypted API key
func GetAPIKey(cfg *Config) (string, error) {
	if cfg.APIKey == "" {
		return "", nil
	}
	return decryptAPIKey(cfg.APIKey)
}

// SetAPIKey encrypts and stores the API key
func SetAPIKey(apiKey string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	encrypted, err := encryptAPIKey(apiKey)
	if err != nil {
		return err
	}

	cfg.APIKey = encrypted
	return Save(cfg)
}

// SetAPIBaseURL stores the API base URL
func SetAPIBaseURL(url string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	cfg.APIBaseURL = url
	return Save(cfg)
}

// GetAPIKeyMasked returns masked API key for display
func GetAPIKeyMasked(cfg *Config) string {
	apiKey, err := GetAPIKey(cfg)
	if err != nil || apiKey == "" {
		return ""
	}
	if len(apiKey) <= 8 {
		return "****"
	}
	return apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
}

func ensureConfigDir() {
	// Try user directory first
	if err := os.MkdirAll(configDir, 0700); err != nil {
		// Fallback to current directory
		configDir = ".adp"
		// Update paths since configDir changed
		configPath = filepath.Join(configDir, configFileName)
		keyPath = filepath.Join(configDir, keyFileName)
		log.Warn().Err(err).Msg("Failed to create config in user directory, using current directory")
	}
	if err := os.MkdirAll(configDir, 0700); err != nil {
		log.Fatal().Err(err).Msg("Failed to create config directory")
	}
}

func encryptAPIKey(apiKey string) (string, error) {
	key := getOrCreateKey()
	plaintext := []byte(apiKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

func decryptAPIKey(encrypted string) (string, error) {
	key := getOrCreateKey()
	ciphertext, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func getOrCreateKey() []byte {
	ensureConfigDir()

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		key := make([]byte, 32) // AES-256
		if _, err := rand.Read(key); err != nil {
			log.Fatal().Err(err).Msg("Failed to generate encryption key")
		}
		if err := os.WriteFile(keyPath, key, 0600); err != nil {
			log.Fatal().Err(err).Msg("Failed to save encryption key")
		}
		return key
	}

	key, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read encryption key")
	}
	return key
}

// GetConfigSummary returns a summary for display
func GetConfigSummary(cfg *Config) map[string]interface{} {
	return map[string]interface{}{
		"configured":    IsConfigured(cfg),
		"api_key_masked": GetAPIKeyMasked(cfg),
		"api_base_url":  cfg.APIBaseURL,
	}
}

// GetCachePath returns the path to the cache file
func GetCachePath() string {
	return filepath.Join(configDir, "app_cache.json")
}
