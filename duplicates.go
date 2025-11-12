package main

import (
	"crypto/md5"
	"fmt"
	"sync"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

// DuplicateTracker rastreia mensagens já copiadas para evitar duplicados.
type DuplicateTracker struct {
	mu      sync.Mutex
	hashes  map[string]bool // hash -> já copiado
}

// NewDuplicateTracker cria um novo rastreador de duplicados.
func NewDuplicateTracker() *DuplicateTracker {
	return &DuplicateTracker{
		hashes: make(map[string]bool),
	}
}

// BuildExistingMessagesIndex constrói um índice das mensagens já existentes no destino.
func (dt *DuplicateTracker) BuildExistingMessagesIndex(client *imapclient.Client, folderName string) error {
	// Selecionar a pasta
	selectData, err := client.Select(folderName, &imap.SelectOptions{ReadOnly: true}).Wait()
	if err != nil {
		return fmt.Errorf("erro ao selecionar pasta para indexação: %w", err)
	}
	
	if selectData.NumMessages == 0 {
		return nil // Pasta vazia, nada a indexar
	}
	
	// Buscar Message-IDs de todas as mensagens
	uidSet := imap.UIDSet{}
	uidSet.AddRange(1, selectData.UIDNext-1)
	
	fetchOptions := &imap.FetchOptions{
		Envelope: true,
	}
	
	messages, err := client.Fetch(uidSet, fetchOptions).Collect()
	if err != nil {
		return fmt.Errorf("erro ao buscar mensagens para indexação: %w", err)
	}
	
	dt.mu.Lock()
	defer dt.mu.Unlock()
	
	for _, msg := range messages {
		if msg.Envelope != nil && msg.Envelope.MessageID != "" {
			// Usar Message-ID como identificador único
			dt.hashes[msg.Envelope.MessageID] = true
		}
	}
	
	return nil
}

// IsDuplicate verifica se uma mensagem já foi copiada.
func (dt *DuplicateTracker) IsDuplicate(messageID string) bool {
	if messageID == "" {
		// Se não há Message-ID, considerar como não duplicado
		return false
	}
	
	dt.mu.Lock()
	defer dt.mu.Unlock()
	
	return dt.hashes[messageID]
}

// MarkAsCopied marca uma mensagem como copiada.
func (dt *DuplicateTracker) MarkAsCopied(messageID string) {
	if messageID == "" {
		return
	}
	
	dt.mu.Lock()
	defer dt.mu.Unlock()
	
	dt.hashes[messageID] = true
}

// GenerateMessageHash gera um hash único para uma mensagem baseado em múltiplos campos.
// Usado como fallback quando Message-ID não está disponível.
func GenerateMessageHash(envelope *imap.Envelope, bodySize int) string {
	if envelope == nil {
		return ""
	}
	
	// Combinar vários campos para criar um identificador único
	data := fmt.Sprintf("%s|%s|%v|%d",
		envelope.Subject,
		envelope.From,
		envelope.Date,
		bodySize,
	)
	
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}
