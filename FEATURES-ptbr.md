# Funcionalidades Avançadas do Migrador IMAP

## 📋 Lista Completa de Funcionalidades

### ✅ Funcionalidades Implementadas

#### 1. **Detecção de Duplicados**
- Verifica se mensagens já existem no destino antes de copiar
- Usa Message-ID como identificador único
- Fallback para hash MD5 (assunto + remetente + data + tamanho) quando Message-ID não disponível
- Configurável via `skip_duplicates` no config.json

#### 2. **Filtro de Pastas - Exclusão**
- Permite excluir pastas específicas da migração
- Útil para pular Drafts, Trash, Junk, etc.
- Configurável via `exclude_folders` no config.json

#### 3. **Filtro de Pastas - Inclusão (Whitelist)**
- Permite migrar APENAS pastas específicas
- Quando configurado, todas as outras pastas são ignoradas
- Configurável via `include_folders` no config.json

#### 4. **Filtro de Data - A Partir De**
- Migra apenas mensagens >= data especificada
- Formato: AAAA-MM-DD (ex: 2024-01-15)
- Configurável via `date_from` no config.json

#### 5. **Filtro de Data - Até**
- Migra apenas mensagens <= data especificada
- Formato: AAAA-MM-DD (ex: 2015-10-10)
- Configurável via `date_to` no config.json

#### 6. **Mapeamento de Nomes de Pastas**
- Permite renomear pastas durante a migração
- Útil para compatibilizar diferentes convenções de nomes
- Exemplo: "INBOX.Sent Messages" → "INBOX.Sent"
- Configurável via `folder_mapping` no config.json

#### 7. **Limite de Tamanho de Mensagem**
- Pula mensagens maiores que X MB
- Útil quando há limitações no servidor de destino
- 0 = sem limite
- Configurável via `max_message_size_mb` no config.json

#### 8. **Modo Dry-Run (Simulação)**
- Executa sem copiar realmente
- Mostra o que seria feito
- Útil para validar filtros antes de executar
- Configurável via `dry_run` no config.json

#### 9. **Retry Automático**
- Tenta novamente mensagens que falharam
- Número configurável de tentativas
- Útil para lidar com erros temporários de rede
- Configurável via `max_retries` no config.json

#### 10. **Achatar Hierarquia de Pastas**
- Converte "INBOX.Sent.2024" em "INBOX_Sent_2024"
- Útil quando servidor de destino tem limitações
- Configurável via `flatten_folders` no config.json

#### 11. **Pastas de Sistema Configuráveis**
- Define nomes alternativos para Drafts, Sent, Junk, Trash, Archive
- Suporta múltiplos nomes (Gmail, Outlook, etc.)
- Configurável via `system_folders` no config.json

---

## 🔧 Arquivo de Configuração (config.json)

```json
{
  "skip_duplicates": true,
  "dry_run": false,
  "max_retries": 3,
  "max_message_size_mb": 50,
  "flatten_folders": false,
  
  "exclude_folders": [
    "INBOX.Drafts",
    "INBOX.Trash",
    "INBOX.Junk"
  ],
  
  "include_folders": [],
  
  "date_from": "2024-01-01",
  "date_to": "2024-12-31",
  
  "folder_mapping": {
    "INBOX.Sent Messages": "INBOX.Sent",
    "INBOX.Deleted Items": "INBOX.Trash"
  },
  
  "system_folders": {
    "drafts": ["Drafts", "INBOX.Drafts", "[Gmail]/Drafts"],
    "sent": ["Sent", "Sent Messages", "INBOX.Sent", "[Gmail]/Sent Mail"],
    "junk": ["Junk", "Spam", "INBOX.Junk", "[Gmail]/Spam"],
    "trash": ["Trash", "Deleted Items", "INBOX.Trash", "[Gmail]/Trash"],
    "archive": ["Archive", "INBOX.Archive", "[Gmail]/All Mail"]
  }
}
```

---

## 📖 Exemplos de Uso

### Exemplo 1: Migração Simples (Sem Filtros)
```json
{
  "skip_duplicates": false,
  "dry_run": false,
  "max_retries": 3,
  "max_message_size_mb": 0,
  "flatten_folders": false,
  "exclude_folders": [],
  "include_folders": [],
  "date_from": "",
  "date_to": "",
  "folder_mapping": {},
  "system_folders": {}
}
```

### Exemplo 2: Migrar Apenas 2024, Sem Lixo
```json
{
  "skip_duplicates": true,
  "exclude_folders": ["INBOX.Trash", "INBOX.Junk", "INBOX.Drafts"],
  "date_from": "2024-01-01",
  "date_to": "2024-12-31"
}
```

### Exemplo 3: Migrar Apenas Pastas Específicas
```json
{
  "skip_duplicates": true,
  "include_folders": [
    "INBOX",
    "INBOX.Important",
    "INBOX.Projects"
  ]
}
```

### Exemplo 4: Teste (Dry-Run)
```json
{
  "dry_run": true,
  "skip_duplicates": true,
  "date_from": "2024-01-01"
}
```

### Exemplo 5: Migração com Limite de Tamanho
```json
{
  "skip_duplicates": true,
  "max_message_size_mb": 25,
  "max_retries": 5
}
```

---

## ⚙️ Regras Aditivas

Todas as regras são **aditivas** (AND lógico). Uma mensagem só é copiada se passar por TODOS os filtros:

1. ✅ Pasta está na whitelist (se configurada)
2. ✅ Pasta NÃO está na blacklist
3. ✅ Data >= date_from (se configurado)
4. ✅ Data <= date_to (se configurado)
5. ✅ Tamanho <= max_message_size_mb (se configurado)
6. ✅ Não é duplicado (se skip_duplicates = true)

**Exemplo:**
```json
{
  "include_folders": ["INBOX", "INBOX.Important"],
  "date_from": "2015-10-11",
  "date_to": "2024-01-14",
  "max_message_size_mb": 50
}
```

Resultado: Migra apenas mensagens de INBOX e INBOX.Important, entre 11/10/2015 e 14/01/2024, menores que 50MB.

---

## 📊 Relatório

O relatório agora inclui:
- Mensagens puladas por filtro de data
- Mensagens puladas por tamanho
- Mensagens puladas por duplicação
- Razão específica para cada mensagem pulada

---

## 🚀 Como Usar

1. Edite `config.json` com suas preferências
2. Execute: `go run *.go`
3. Verifique o relatório em `relatorios/`

---

## ⚠️ Notas Importantes

- **Validação de Pastas**: Se usar `include_folders` ou `exclude_folders`, o programa verifica se todas as pastas existem antes de começar
- **Dry-Run**: Sempre teste com `dry_run: true` primeiro
- **Duplicados**: A detecção usa Message-ID, que é confiável mas não 100% garantido
- **Performance**: Detecção de duplicados adiciona overhead (busca mensagens existentes)
