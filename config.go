package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// MigrationConfig armazena todas as opções de configuração da migração.
type MigrationConfig struct {
	// Opções gerais
	AccountsFile            string `json:"accounts_file"`
	MaxConcurrentMigrations int    `json:"max_concurrent_migrations"`
	SkipDuplicates          bool   `json:"skip_duplicates"`
	DryRun                  bool   `json:"dry_run"`
	MaxRetries              int    `json:"max_retries"`
	MaxMessageSizeMB        int    `json:"max_message_size_mb"`
	FlattenFolders          bool   `json:"flatten_folders"`
	
	// Filtros de pastas
	ExcludeFolders     []string          `json:"exclude_folders"`
	IncludeFolders     []string          `json:"include_folders"`
	
	// Filtros de data
	DateFrom           string            `json:"date_from"` // Formato: 2006-01-02
	DateTo             string            `json:"date_to"`   // Formato: 2006-01-02
	
	// Mapeamento de pastas
	FolderMapping      map[string]string `json:"folder_mapping"`
	SystemFolders      SystemFolders     `json:"system_folders"`
	
	// Campos internos (parseados)
	dateFromParsed     *time.Time
	dateToParsed       *time.Time
}

// SystemFolders define nomes alternativos para pastas de sistema.
type SystemFolders struct {
	Drafts  []string `json:"drafts"`
	Sent    []string `json:"sent"`
	Junk    []string `json:"junk"`
	Trash   []string `json:"trash"`
	Archive []string `json:"archive"`
}

// DefaultConfig retorna uma configuração padrão.
func DefaultConfig() MigrationConfig {
	return MigrationConfig{
		AccountsFile:            "accounts.csv",
		MaxConcurrentMigrations: 5,
		SkipDuplicates:          false,
		DryRun:                  false,
		MaxRetries:              3,
		MaxMessageSizeMB:        0, // 0 = sem limite
		FlattenFolders:          false,
		ExcludeFolders:   []string{},
		IncludeFolders:   []string{},
		DateFrom:         "",
		DateTo:           "",
		FolderMapping:    make(map[string]string),
		SystemFolders: SystemFolders{
			Drafts:  []string{"Drafts", "INBOX.Drafts", "[Gmail]/Drafts"},
			Sent:    []string{"Sent", "Sent Messages", "INBOX.Sent", "[Gmail]/Sent Mail"},
			Junk:    []string{"Junk", "Spam", "INBOX.Junk", "[Gmail]/Spam"},
			Trash:   []string{"Trash", "Deleted Items", "INBOX.Trash", "[Gmail]/Trash"},
			Archive: []string{"Archive", "INBOX.Archive", "[Gmail]/All Mail"},
		},
	}
}

// LoadConfig carrega a configuração de um ficheiro JSON.
func LoadConfig(filePath string) (MigrationConfig, error) {
	// Se o ficheiro não existir, usar configuração padrão
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}
	
	file, err := os.Open(filePath)
	if err != nil {
		return MigrationConfig{}, fmt.Errorf("erro ao abrir ficheiro de configuração: %w", err)
	}
	defer file.Close()
	
	var config MigrationConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return MigrationConfig{}, fmt.Errorf("erro ao parsear configuração JSON: %w", err)
	}
	
	// Parsear datas
	if config.DateFrom != "" {
		t, err := time.Parse("2006-01-02", config.DateFrom)
		if err != nil {
			return MigrationConfig{}, fmt.Errorf("formato de data inválido em date_from: %w", err)
		}
		config.dateFromParsed = &t
	}
	
	if config.DateTo != "" {
		t, err := time.Parse("2006-01-02", config.DateTo)
		if err != nil {
			return MigrationConfig{}, fmt.Errorf("formato de data inválido em date_to: %w", err)
		}
		// Ajustar para o final do dia
		endOfDay := t.Add(24*time.Hour - time.Second)
		config.dateToParsed = &endOfDay
	}
	
	// Se AccountsFile não foi especificado, usar padrão
	if config.AccountsFile == "" {
		config.AccountsFile = "accounts.csv"
	}
	
	// Se MaxConcurrentMigrations não foi especificado ou é inválido, usar padrão
	if config.MaxConcurrentMigrations <= 0 {
		config.MaxConcurrentMigrations = 5
	}
	
	return config, nil
}

// ShouldIncludeFolder verifica se uma pasta deve ser incluída na migração.
func (c *MigrationConfig) ShouldIncludeFolder(folderName string) bool {
	// Se há whitelist, apenas pastas nela são incluídas
	if len(c.IncludeFolders) > 0 {
		found := false
		for _, f := range c.IncludeFolders {
			if f == folderName {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	// Verificar blacklist
	for _, f := range c.ExcludeFolders {
		if f == folderName {
			return false
		}
	}
	
	return true
}

// ShouldIncludeMessage verifica se uma mensagem deve ser incluída baseado em filtros de data e tamanho.
func (c *MigrationConfig) ShouldIncludeMessage(messageDate time.Time, messageSize int) (bool, string) {
	// Verificar data mínima
	if c.dateFromParsed != nil && messageDate.Before(*c.dateFromParsed) {
		return false, fmt.Sprintf("data anterior a %s", c.DateFrom)
	}
	
	// Verificar data máxima
	if c.dateToParsed != nil && messageDate.After(*c.dateToParsed) {
		return false, fmt.Sprintf("data posterior a %s", c.DateTo)
	}
	
	// Verificar tamanho máximo
	if c.MaxMessageSizeMB > 0 {
		maxBytes := c.MaxMessageSizeMB * 1024 * 1024
		if messageSize > maxBytes {
			return false, fmt.Sprintf("tamanho %d bytes excede limite de %d MB", messageSize, c.MaxMessageSizeMB)
		}
	}
	
	return true, ""
}

// GetMappedFolderName retorna o nome mapeado de uma pasta, se houver.
func (c *MigrationConfig) GetMappedFolderName(originalName string) string {
	if mapped, ok := c.FolderMapping[originalName]; ok {
		return mapped
	}
	return originalName
}

// FlattenFolderName converte hierarquia de pastas em nome plano.
func (c *MigrationConfig) FlattenFolderName(folderName string) string {
	if !c.FlattenFolders {
		return folderName
	}
	// Substituir separadores por underscore
	flattened := folderName
	flattened = replaceAll(flattened, ".", "_")
	flattened = replaceAll(flattened, "/", "_")
	return flattened
}

// replaceAll substitui todas as ocorrências de old por new em s.
func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); i++ {
		if i <= len(s)-len(old) && s[i:i+len(old)] == old {
			result += new
			i += len(old) - 1
		} else {
			result += string(s[i])
		}
	}
	return result
}
