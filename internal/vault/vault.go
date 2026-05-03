package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/wolfsouldev/ssh/internal/crypto"
)

// AuthType represents the type of authentication stored.
type AuthType string

const (
	AuthPassword   AuthType = "password"
	AuthPrivateKey AuthType = "privatekey"
)

// Credential represents a stored SSH credential.
type Credential struct {
	Alias      string   `json:"alias"`
	Host       string   `json:"host"`
	Port       int      `json:"port"`
	User       string   `json:"user"`
	AuthType   AuthType `json:"auth_type"`
	Password   string   `json:"password,omitempty"`
	PrivateKey string   `json:"private_key,omitempty"`
	KeyPass    string   `json:"key_passphrase,omitempty"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
}

// VaultData holds all credentials.
type VaultData struct {
	Version     int          `json:"version"`
	Credentials []Credential `json:"credentials"`
}

// Vault manages encrypted credential storage.
type Vault struct {
	path string
}

// New creates a new Vault instance pointing to the default vault file.
func New() (*Vault, error) {
	dir, err := vaultDir()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("creating vault directory: %w", err)
	}
	return &Vault{path: filepath.Join(dir, "vault.enc")}, nil
}

// Exists returns true if the vault file exists.
func (v *Vault) Exists() bool {
	_, err := os.Stat(v.path)
	return err == nil
}

// Init creates a new empty vault encrypted with the master password.
func (v *Vault) Init(masterPass []byte) error {
	data := VaultData{Version: 1, Credentials: []Credential{}}
	return v.save(data, masterPass)
}

// Load decrypts and loads the vault data.
func (v *Vault) Load(masterPass []byte) (*VaultData, error) {
	encData, err := os.ReadFile(v.path)
	if err != nil {
		return nil, fmt.Errorf("reading vault file: %w", err)
	}

	plaintext, err := crypto.Decrypt(encData, masterPass)
	if err != nil {
		return nil, err
	}

	var data VaultData
	if err := json.Unmarshal(plaintext, &data); err != nil {
		return nil, fmt.Errorf("parsing vault data: %w", err)
	}

	return &data, nil
}

// Save encrypts and saves the vault data.
func (v *Vault) save(data VaultData, masterPass []byte) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding vault data: %w", err)
	}

	encData, err := crypto.Encrypt(jsonData, masterPass)
	if err != nil {
		return err
	}

	if err := os.WriteFile(v.path, encData, 0600); err != nil {
		return fmt.Errorf("writing vault file: %w", err)
	}

	return nil
}

// AddCredential adds a new credential to the vault.
func (v *Vault) AddCredential(masterPass []byte, cred Credential) error {
	data, err := v.Load(masterPass)
	if err != nil {
		return err
	}

	now := time.Now().Format(time.RFC3339)
	cred.CreatedAt = now
	cred.UpdatedAt = now

	if cred.Alias == "" {
		cred.Alias = fmt.Sprintf("%s@%s", cred.User, cred.Host)
	}
	if cred.Port == 0 {
		cred.Port = 22
	}

	// Check for duplicates
	for _, c := range data.Credentials {
		if c.Alias == cred.Alias {
			return fmt.Errorf("credential with alias '%s' already exists", cred.Alias)
		}
	}

	data.Credentials = append(data.Credentials, cred)
	return v.save(*data, masterPass)
}

// UpdateCredential updates an existing credential by alias.
func (v *Vault) UpdateCredential(masterPass []byte, alias string, updater func(*Credential)) error {
	data, err := v.Load(masterPass)
	if err != nil {
		return err
	}

	found := false
	for i := range data.Credentials {
		if data.Credentials[i].Alias == alias {
			updater(&data.Credentials[i])
			data.Credentials[i].UpdatedAt = time.Now().Format(time.RFC3339)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("credential '%s' not found", alias)
	}

	return v.save(*data, masterPass)
}

// DeleteCredential removes a credential by alias.
func (v *Vault) DeleteCredential(masterPass []byte, alias string) error {
	data, err := v.Load(masterPass)
	if err != nil {
		return err
	}

	found := false
	filtered := make([]Credential, 0, len(data.Credentials))
	for _, c := range data.Credentials {
		if c.Alias == alias {
			found = true
			continue
		}
		filtered = append(filtered, c)
	}

	if !found {
		return fmt.Errorf("credential '%s' not found", alias)
	}

	data.Credentials = filtered
	return v.save(*data, masterPass)
}

// FindByHostUser finds a credential by host and user.
func (v *Vault) FindByHostUser(data *VaultData, host, user string) *Credential {
	for i := range data.Credentials {
		if data.Credentials[i].Host == host && data.Credentials[i].User == user {
			return &data.Credentials[i]
		}
	}
	return nil
}

// FindByAlias finds a credential by alias.
func (v *Vault) FindByAlias(data *VaultData, alias string) *Credential {
	for i := range data.Credentials {
		if data.Credentials[i].Alias == alias {
			return &data.Credentials[i]
		}
	}
	return nil
}

// ChangeMasterPassword re-encrypts the vault with a new master password.
func (v *Vault) ChangeMasterPassword(oldPass, newPass []byte) error {
	data, err := v.Load(oldPass)
	if err != nil {
		return errors.New("wrong current master password")
	}
	return v.save(*data, newPass)
}

func vaultDir() (string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("finding home directory: %w", err)
		}
		appData = filepath.Join(home, ".config")
	}
	return filepath.Join(appData, "sshh"), nil
}
