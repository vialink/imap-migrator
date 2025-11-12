package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

// MigrationAccount armazena as informações para uma única migração.
type MigrationAccount struct {
	LineNumber       int
	SourceEmail      string
	SourceUser       string
	SourcePass       string
	SourceHost       string
	DestinationEmail string
	DestinationUser  string
	DestinationPass  string
	DestinationHost  string
}

// FolderStats armazena estatísticas de uma pasta.
type FolderStats struct {
	Name            string
	SourceMessages  uint32
	CopiedMessages  int
	FailedMessages  int
	SkippedMessages int
}

// MigrationReport armazena o relatório completo de uma migração.
type MigrationReport struct {
	SourceEmail      string
	DestinationEmail string
	StartTime        time.Time
	EndTime          time.Time
	Duration         time.Duration
	Folders          []FolderStats
	Errors           []string
	Success          bool
	TotalFolders     int
	TotalSourceMsgs  uint32
	TotalCopied      int
	TotalFailed      int
	TotalSkipped     int
}

// readCSV lê o ficheiro de contas e retorna uma lista de MigrationAccount.
func readCSV(filePath string) ([]MigrationAccount, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir o ficheiro CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("erro ao ler o ficheiro CSV: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("ficheiro CSV está vazio")
	}

	var accounts []MigrationAccount
	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) < 8 {
			log.Printf("AVISO: Linha %d tem menos de 8 colunas, ignorando.", i+1)
			continue
		}

		acc := MigrationAccount{
			LineNumber:       i + 1,
			SourceEmail:      strings.TrimSpace(record[0]),
			SourceUser:       strings.TrimSpace(record[1]),
			SourcePass:       strings.TrimSpace(record[2]),
			SourceHost:       strings.TrimSpace(record[3]),
			DestinationEmail: strings.TrimSpace(record[4]),
			DestinationUser:  strings.TrimSpace(record[5]),
			DestinationPass:  strings.TrimSpace(record[6]),
			DestinationHost:  strings.TrimSpace(record[7]),
		}

		accounts = append(accounts, acc)
	}

	return accounts, nil
}

// connectClient estabelece conexão TLS com um servidor IMAP.
func connectClient(host, user, pass string) (*imapclient.Client, error) {
	addr := net.JoinHostPort(host, "993")
	c, err := imapclient.DialTLS(addr, nil)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar via TLS: %w", err)
	}

	if err := c.Login(user, pass).Wait(); err != nil {
		c.Close()
		return nil, fmt.Errorf("falha ao fazer login: %w", err)
	}

	return c, nil
}

// testConnection testa a conexão com um servidor IMAP.
func testConnection(host, user, pass string) error {
	client, err := connectClient(host, user, pass)
	if err != nil {
		return err
	}
	defer client.Logout()
	return nil
}

// isConnectionClosed verifica se um erro indica conexão fechada.
func isConnectionClosed(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "closed network connection") ||
		strings.Contains(errStr, "EOF") ||
		strings.Contains(errStr, "broken pipe") ||
		strings.Contains(errStr, "connection reset")
}

// reconnectIfNeeded reconecta se o erro indicar conexão fechada.
func reconnectIfNeeded(client **imapclient.Client, host, user, pass string, err error) error {
	if !isConnectionClosed(err) {
		return err
	}

	log.Printf("Conexão fechada detectada. Tentando reconectar a %s...", host)

	if *client != nil {
		(*client).Close()
	}

	newClient, err := connectClient(host, user, pass)
	if err != nil {
		return fmt.Errorf("falha ao reconectar: %w", err)
	}

	*client = newClient
	log.Printf("Reconexão bem-sucedida a %s", host)
	return nil
}

// filterValidFlags remove flags que podem causar problemas.
func filterValidFlags(flags []imap.Flag) []imap.Flag {
	var validFlags []imap.Flag
	for _, flag := range flags {
		if flag != "\\Recent" {
			validFlags = append(validFlags, flag)
		}
	}
	return validFlags
}

// migrateAccount executa a migração para uma única conta.
func migrateAccount(acc MigrationAccount, config MigrationConfig) error {
	log.Printf("[ÍNÍCIO MIGRAÇÃO] %s -> %s", acc.SourceEmail, acc.DestinationEmail)

	// Inicializar relatório
	report := MigrationReport{
		SourceEmail:      acc.SourceEmail,
		DestinationEmail: acc.DestinationEmail,
		StartTime:        time.Now(),
		Folders:          []FolderStats{},
		Errors:           []string{},
		Success:          false,
	}
	defer func() {
		report.EndTime = time.Now()
		report.Duration = report.EndTime.Sub(report.StartTime)
		if err := saveReport(report); err != nil {
			log.Printf("[%s] AVISO: não foi possível guardar o relatório: %v", acc.SourceEmail, err)
		}
	}()

	sourceClient, err := connectClient(acc.SourceHost, acc.SourceUser, acc.SourcePass)
	if err != nil {
		return fmt.Errorf("erro ao conectar à origem: %w", err)
	}
	defer sourceClient.Logout()

	destClient, err := connectClient(acc.DestinationHost, acc.DestinationUser, acc.DestinationPass)
	if err != nil {
		return fmt.Errorf("erro ao conectar ao destino: %w", err)
	}
	defer destClient.Logout()

	mailboxes, err := sourceClient.List("", "*", nil).Collect()
	if err != nil {
		return fmt.Errorf("falha ao listar pastas na origem: %w", err)
	}

	log.Printf("[%s] Encontradas %d pastas para migrar.", acc.SourceEmail, len(mailboxes))

	// Inicializar rastreador de duplicados se necessário
	var dupTracker *DuplicateTracker
	if config.SkipDuplicates {
		dupTracker = NewDuplicateTracker()
		log.Printf("[%s] Detecção de duplicados ativada", acc.SourceEmail)
	}

	for _, mb := range mailboxes {
		if slices.Contains(mb.Attrs, imap.MailboxAttrNoSelect) {
			log.Printf("[%s] Ignorando pasta não selecionável: %s", acc.SourceEmail, mb.Mailbox)
			continue
		}

		folderName := mb.Mailbox

		// Filtro de pastas
		if !config.ShouldIncludeFolder(folderName) {
			log.Printf("[%s] Pasta '%s' excluída por filtro de configuração", acc.SourceEmail, folderName)
			continue
		}

		log.Printf("[%s] Processando pasta: %s", acc.SourceEmail, folderName)

		folderStats := FolderStats{
			Name:            folderName,
			SourceMessages:  0,
			CopiedMessages:  0,
			FailedMessages:  0,
			SkippedMessages: 0,
		}

		// Mapeamento e flatten de nome
		destFolderName := config.GetMappedFolderName(folderName)
		destFolderName = config.FlattenFolderName(destFolderName)

		if destFolderName != folderName {
			log.Printf("[%s] Pasta '%s' será criada como '%s' no destino", acc.SourceEmail, folderName, destFolderName)
		}

		// Criar pasta no destino
		if !config.DryRun {
			err := destClient.Create(destFolderName, nil).Wait()
			if err != nil {
				if reconnectErr := reconnectIfNeeded(&destClient, acc.DestinationHost, acc.DestinationUser, acc.DestinationPass, err); reconnectErr == nil {
					err = destClient.Create(destFolderName, nil).Wait()
					if err != nil {
						log.Printf("[%s] Aviso: não foi possível criar a pasta '%s' no destino (pode já existir): %v", acc.DestinationEmail, destFolderName, err)
					}
				} else {
					log.Printf("[%s] Aviso: não foi possível criar a pasta '%s' no destino (pode já existir): %v", acc.DestinationEmail, destFolderName, err)
				}
			}

			// Construir índice de duplicados se necessário
			if config.SkipDuplicates {
				log.Printf("[%s] Construindo índice de mensagens existentes na pasta '%s'...", acc.DestinationEmail, destFolderName)
				if err := dupTracker.BuildExistingMessagesIndex(destClient, destFolderName); err != nil {
					log.Printf("[%s] AVISO: não foi possível construir índice de duplicados para '%s': %v", acc.DestinationEmail, destFolderName, err)
				}
			}
		} else {
			log.Printf("[%s] [DRY-RUN] Pasta '%s' seria criada como '%s'", acc.SourceEmail, folderName, destFolderName)
		}

		// Selecionar pasta de origem
		sourceData, err := sourceClient.Select(folderName, nil).Wait()
		if err != nil {
			if reconnectErr := reconnectIfNeeded(&sourceClient, acc.SourceHost, acc.SourceUser, acc.SourcePass, err); reconnectErr == nil {
				sourceData, err = sourceClient.Select(folderName, nil).Wait()
				if err != nil {
					log.Printf("[%s] ERRO: não foi possível selecionar a pasta '%s' na origem: %v", acc.SourceEmail, folderName, err)
					continue
				}
			} else {
				log.Printf("[%s] ERRO: não foi possível selecionar a pasta '%s' na origem: %v", acc.SourceEmail, folderName, err)
				continue
			}
		}

		log.Printf("[%s] Pasta '%s': servidor reporta %d mensagens (UIDNext: %d, UIDValidity: %d)",
			acc.SourceEmail, folderName, sourceData.NumMessages, sourceData.UIDNext, sourceData.UIDValidity)

		folderStats.SourceMessages = sourceData.NumMessages

		if sourceData.NumMessages == 0 {
			log.Printf("[%s] Pasta '%s' está vazia, passando para a próxima.", acc.SourceEmail, folderName)
			report.Folders = append(report.Folders, folderStats)
			continue
		}

		// Buscar mensagens
		uidSet := imap.UIDSet{}
		uidSet.AddRange(1, sourceData.UIDNext-1)
		fetchOptions := &imap.FetchOptions{
			BodySection: []*imap.FetchItemBodySection{{}},
			Flags:       true,
			Envelope:    true,
		}

		log.Printf("[%s] Fazendo fetch de mensagens da pasta '%s' usando UIDs...", acc.SourceEmail, folderName)

		messages, err := sourceClient.Fetch(uidSet, fetchOptions).Collect()
		if err != nil {
			if reconnectErr := reconnectIfNeeded(&sourceClient, acc.SourceHost, acc.SourceUser, acc.SourcePass, err); reconnectErr == nil {
				sourceData, err = sourceClient.Select(folderName, nil).Wait()
				if err != nil {
					log.Printf("[%s] ERRO: não foi possível reselecionar pasta após reconexão: %v", acc.SourceEmail, err)
					continue
				}
				messages, err = sourceClient.Fetch(uidSet, fetchOptions).Collect()
				if err != nil {
					log.Printf("[%s] ERRO: falha ao obter mensagens após reconexão: %v", acc.SourceEmail, err)
					continue
				}
			} else {
				log.Printf("[%s] ERRO: falha ao obter mensagens: %v", acc.SourceEmail, err)
				continue
			}
		}

		log.Printf("[%s] Pasta '%s' tem %d mensagens para processar.", acc.SourceEmail, folderName, len(messages))

		// Selecionar pasta de destino
		if !config.DryRun {
			_, err = destClient.Select(destFolderName, nil).Wait()
			if err != nil {
				if reconnectErr := reconnectIfNeeded(&destClient, acc.DestinationHost, acc.DestinationUser, acc.DestinationPass, err); reconnectErr == nil {
					_, err = destClient.Select(destFolderName, nil).Wait()
					if err != nil {
						log.Printf("[%s] ERRO: não foi possível selecionar a pasta '%s' no destino após reconexão: %v", acc.DestinationEmail, destFolderName, err)
						continue
					}
				} else {
					log.Printf("[%s] ERRO: não foi possível selecionar a pasta '%s' no destino: %v", acc.DestinationEmail, destFolderName, err)
					continue
				}
			}
		}

		copiedCount := 0
		for i, msg := range messages {
			// Verificar corpo da mensagem
			if len(msg.BodySection) == 0 {
				log.Printf("[%s] AVISO: mensagem %d/%d da pasta '%s' tem corpo vazio, pulando.", acc.SourceEmail, i+1, len(messages), folderName)
				folderStats.SkippedMessages++
				continue
			}
			bodyBytes := msg.BodySection[0].Bytes

			// Filtro de tamanho
			if shouldInclude, reason := config.ShouldIncludeMessage(msg.Envelope.Date, len(bodyBytes)); !shouldInclude {
				log.Printf("[%s] Mensagem %d/%d pulada: %s", acc.SourceEmail, i+1, len(messages), reason)
				folderStats.SkippedMessages++
				continue
			}

			// Verificar duplicados
			if config.SkipDuplicates {
				messageID := msg.Envelope.MessageID
				if messageID == "" {
					messageID = GenerateMessageHash(msg.Envelope, len(bodyBytes))
				}
				if dupTracker.IsDuplicate(messageID) {
					log.Printf("[%s] Mensagem %d/%d pulada: duplicada (Message-ID: %s)", acc.SourceEmail, i+1, len(messages), messageID)
					folderStats.SkippedMessages++
					continue
				}
				dupTracker.MarkAsCopied(messageID)
			}

			validFlags := filterValidFlags(msg.Flags)

			log.Printf("[%s] Copiando mensagem %d/%d da pasta '%s' (tamanho: %d bytes)...", acc.SourceEmail, i+1, len(messages), folderName, len(bodyBytes))

			if config.DryRun {
				log.Printf("[%s] [DRY-RUN] Mensagem %d/%d seria copiada", acc.SourceEmail, i+1, len(messages))
				folderStats.CopiedMessages++
				copiedCount++
				continue
			}

			// Tentar copiar com retry
			var copyErr error
			for attempt := 0; attempt <= config.MaxRetries; attempt++ {
				if attempt > 0 {
					log.Printf("[%s] Tentativa %d/%d para mensagem %d/%d", acc.SourceEmail, attempt, config.MaxRetries, i+1, len(messages))
				}

				appendCmd := destClient.Append(destFolderName, int64(len(bodyBytes)), &imap.AppendOptions{
					Flags: validFlags,
					Time:  msg.Envelope.Date,
				})

				var writeErr error
				_, writeErr = io.Copy(appendCmd, bytes.NewReader(bodyBytes))

				closeErr := appendCmd.Close()

				if writeErr != nil {
					if strings.Contains(writeErr.Error(), "OVERQUOTA") || strings.Contains(writeErr.Error(), "Quota exceeded") {
						errMsg := fmt.Sprintf("Quota excedida no destino ao copiar mensagem %d/%d da pasta '%s'", i+1, len(messages), folderName)
						report.Errors = append(report.Errors, errMsg)
						log.Printf("[%s] ERRO CRÍTICO: Quota excedida no destino!", acc.DestinationEmail)
						return fmt.Errorf("quota excedida no destino: %w", writeErr)
					}
					copyErr = writeErr
					continue
				}

				if closeErr != nil {
					if strings.Contains(closeErr.Error(), "OVERQUOTA") || strings.Contains(closeErr.Error(), "Quota exceeded") {
						errMsg := fmt.Sprintf("Quota excedida no destino ao finalizar mensagem %d/%d da pasta '%s'", i+1, len(messages), folderName)
						report.Errors = append(report.Errors, errMsg)
						log.Printf("[%s] ERRO CRÍTICO: Quota excedida no destino!", acc.DestinationEmail)
						return fmt.Errorf("quota excedida no destino: %w", closeErr)
					}
					copyErr = closeErr
					continue
				}

				// Sucesso
				copyErr = nil
				break
			}

			if copyErr != nil {
				errMsg := fmt.Sprintf("Falha ao copiar mensagem %d/%d da pasta '%s' após %d tentativas: %v", i+1, len(messages), folderName, config.MaxRetries+1, copyErr)
				report.Errors = append(report.Errors, errMsg)
				folderStats.FailedMessages++
				log.Printf("[%s] ERRO: %s", acc.SourceEmail, errMsg)
				continue
			}

			copiedCount++
			folderStats.CopiedMessages++
			log.Printf("[%s] Mensagem %d/%d copiada com sucesso para '%s'", acc.SourceEmail, i+1, len(messages), destFolderName)
		}

		log.Printf("[%s] Pasta '%s': %d/%d mensagens copiadas com sucesso.", acc.SourceEmail, folderName, copiedCount, len(messages))

		report.Folders = append(report.Folders, folderStats)
	}

	// Calcular totais
	report.TotalFolders = len(report.Folders)
	for _, folder := range report.Folders {
		report.TotalSourceMsgs += folder.SourceMessages
		report.TotalCopied += folder.CopiedMessages
		report.TotalFailed += folder.FailedMessages
		report.TotalSkipped += folder.SkippedMessages
	}
	report.Success = true

	log.Printf("[FIM MIGRAÇÃO] %s -> %s", acc.SourceEmail, acc.DestinationEmail)
	log.Printf("[RESUMO] Total: %d mensagens na origem, %d copiadas, %d falhadas, %d puladas",
		report.TotalSourceMsgs, report.TotalCopied, report.TotalFailed, report.TotalSkipped)
	return nil
}

func main() {
	log.Println("Iniciando migrador IMAP...")

	// Carregar configuração
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatalf("ERRO FATAL ao carregar configuração: %v", err)
	}

	if config.DryRun {
		log.Println("*** MODO DRY-RUN ATIVADO - Nenhuma mensagem será realmente copiada ***")
	}

	log.Println("Iniciando a verificação de conexões...")

	accounts, err := readCSV(config.AccountsFile)
	if err != nil {
		log.Fatalf("ERRO FATAL: %v", err)
	}

	if len(accounts) == 0 {
		log.Printf("Nenhuma conta encontrada no ficheiro %s.", config.AccountsFile)
		return
	}

	// FASE 1: Verificação
	var wgCheck sync.WaitGroup
	results := make(chan string, len(accounts)*2)
	allConnectionsOK := true
	var mu sync.Mutex

	for _, acc := range accounts {
		wgCheck.Add(2)
		go func(a MigrationAccount) {
			defer wgCheck.Done()
			err := testConnection(a.SourceHost, a.SourceUser, a.SourcePass)
			mu.Lock()
			if err != nil {
				results <- fmt.Sprintf("❌ [Linha %d] Origem %s (%s): FALHOU - %v", a.LineNumber, a.SourceEmail, a.SourceHost, err)
				allConnectionsOK = false
			} else {
				results <- fmt.Sprintf("✅ [Linha %d] Origem %s (%s): OK", a.LineNumber, a.SourceEmail, a.SourceHost)
			}
			mu.Unlock()
		}(acc)

		go func(a MigrationAccount) {
			defer wgCheck.Done()
			err := testConnection(a.DestinationHost, a.DestinationUser, a.DestinationPass)
			mu.Lock()
			if err != nil {
				results <- fmt.Sprintf("❌ [Linha %d] Destino %s (%s): FALHOU - %v", a.LineNumber, a.DestinationEmail, a.DestinationHost, err)
				allConnectionsOK = false
			} else {
				results <- fmt.Sprintf("✅ [Linha %d] Destino %s (%s): OK", a.LineNumber, a.DestinationEmail, a.DestinationHost)
			}
			mu.Unlock()
		}(acc)
	}

	wgCheck.Wait()
	close(results)

	fmt.Println("\n--- Relatório de Verificação de Conexões ---")
	var sortedResults []string
	for res := range results {
		sortedResults = append(sortedResults, res)
	}
	slices.Sort(sortedResults)
	for _, res := range sortedResults {
		fmt.Println(res)
	}
	fmt.Println("-------------------------------------------")

	// FASE 2: Migração
	if allConnectionsOK {
		log.Println("\nTodas as conexões foram verificadas com sucesso. Iniciando a migração...")
		log.Printf("Máximo de migrações simultâneas: %d\n", config.MaxConcurrentMigrations)

		semaphore := make(chan struct{}, config.MaxConcurrentMigrations)
		var wgMigrate sync.WaitGroup

		for _, acc := range accounts {
			wgMigrate.Add(1)
			semaphore <- struct{}{}

			go func(a MigrationAccount) {
				defer wgMigrate.Done()
				if err := migrateAccount(a, config); err != nil {
					log.Printf("ERRO NA MIGRAÇÃO de %s: %v", a.SourceEmail, err)
				}
				<-semaphore
			}(acc)
		}

		wgMigrate.Wait()
		log.Println("\nProcesso de migração concluído.")

	} else {
		log.Printf("\nForam encontrados erros em uma ou mais conexões. Corrija o ficheiro '%s' e tente novamente.", config.AccountsFile)
	}
}
