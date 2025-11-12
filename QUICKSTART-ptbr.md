# 🚀 Início Rápido - Migrador IMAP v11

## 📦 Conteúdo do Pacote

Este pacote contém tudo que você precisa para migrar contas de email IMAP.

### Arquivos de Código
- `main.go` - Programa principal
- `config.go` - Sistema de configuração
- `duplicates.go` - Detecção de duplicados
- `report.go` - Geração de relatórios
- `go.mod` - Dependências do Go

### Arquivos de Configuração
- `config.json` - Configuração padrão (use este)
- `config.json.exemplo` - Exemplo com todas as opções
- `accounts.csv` - **VOCÊ PRECISA CRIAR ESTE ARQUIVO**
- `accounts.csv.exemplo` - Modelo para criar seu accounts.csv

### Documentação
- `INICIO_RAPIDO.md` - Este arquivo (comece aqui!)
- `FUNCIONALIDADES.md` - Guia completo de funcionalidades
- `README.md` - Documentação geral
- `GMAIL.md` - Guia específico para Gmail
- `CHANGELOG.md` - Histórico de mudanças
- `EXEMPLO_RELATORIO.txt` - Exemplo de relatório gerado

---

## ⚡ Passos para Começar

### 1. **Extrair o Pacote**
```bash
unzip imap-migrator-v11-completo.zip
cd imap-migrator-v11
```

### 2. **Criar o Arquivo de Contas**
Copie o exemplo e edite com suas contas:
```bash
cp accounts.csv.exemplo accounts.csv
nano accounts.csv  # ou use seu editor preferido
```

**Formato do accounts.csv:**
```csv
email_origem,conta_origem,senha_origem,servidor_origem,email_destino,conta_destino,senha_destino,servidor_destino
user@origem.com,user,senha123,imap.origem.com,user@destino.com,user,senha456,imap.destino.com
```

### 3. **Configurar Opções (Opcional)**
Se quiser usar filtros, edite `config.json`:
```bash
nano config.json
```

**Para começar simples, deixe assim:**
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

### 4. **Executar o Programa**
```bash
go run *.go
```

**Ou compile primeiro:**
```bash
go build -o migrador *.go
./migrador
```

### 5. **Ver os Relatórios**
Após a migração, os relatórios estarão em:
```bash
ls relatorios/
cat relatorios/migracao_*.txt
```

---

## 🧪 Teste Primeiro (Recomendado!)

Antes de migrar de verdade, faça um teste com **dry-run**:

1. Edite `config.json`:
```json
{
  "dry_run": true
}
```

2. Execute:
```bash
go run *.go
```

3. Veja o que seria feito sem copiar nada!

---

## 📋 Exemplos Rápidos

### Migração Simples (Tudo)
```json
{
  "skip_duplicates": false,
  "dry_run": false
}
```

### Migrar Apenas 2024
```json
{
  "skip_duplicates": true,
  "date_from": "2024-01-01",
  "date_to": "2024-12-31"
}
```

### Excluir Lixo e Spam
```json
{
  "exclude_folders": [
    "INBOX.Trash",
    "INBOX.Junk",
    "INBOX.Drafts"
  ]
}
```

### Migrar Apenas Pastas Importantes
```json
{
  "include_folders": [
    "INBOX",
    "INBOX.Important",
    "INBOX.Projects"
  ]
}
```

---

## ❓ Perguntas Frequentes

### Como executar?
```bash
go run *.go
```

### Preciso instalar algo?
Sim, apenas o Go (versão 1.19+). As dependências são baixadas automaticamente.

### Como instalar o Go?
- **Ubuntu/Debian**: `sudo apt install golang-go`
- **macOS**: `brew install go`
- **Windows**: Baixe de https://go.dev/dl/

### O programa é seguro?
Sim! Ele:
- ✅ Testa todas as conexões antes de começar
- ✅ Não apaga nada da origem
- ✅ Gera relatórios detalhados
- ✅ Suporta dry-run para testar

### Posso parar no meio?
Sim, use Ctrl+C. O programa para graciosamente. Você pode executar novamente e ele continuará (use `skip_duplicates: true` para evitar recopiar).

### Como migrar Gmail?
Veja o arquivo `GMAIL.md` para instruções específicas.

### Onde estão os relatórios?
No diretório `relatorios/` que é criado automaticamente.

---

## 🆘 Problemas Comuns

### "Quota exceeded"
A conta de destino está cheia. Aumente o limite de armazenamento.

### "Connection closed"
Normal! O programa reconecta automaticamente.

### "Invalid credentials"
Verifique usuário e senha no `accounts.csv`.

### Para Gmail: "Authentication failed"
Você precisa usar **Senha de Aplicação**, não a senha normal. Veja `GMAIL.md`.

---

## 📚 Próximos Passos

1. ✅ Leia `FUNCIONALIDADES.md` para ver todas as opções avançadas
2. ✅ Teste com `dry_run: true` primeiro
3. ✅ Execute a migração real
4. ✅ Verifique os relatórios em `relatorios/`

---

## 🎯 Comando Completo

```bash
# 1. Extrair
unzip imap-migrator-v11-completo.zip
cd imap-migrator-v11

# 2. Criar accounts.csv
cp accounts.csv.exemplo accounts.csv
nano accounts.csv

# 3. Testar (dry-run)
# Edite config.json e coloque "dry_run": true
go run *.go

# 4. Executar de verdade
# Edite config.json e coloque "dry_run": false
go run *.go

# 5. Ver relatórios
ls relatorios/
```

---

## 💡 Dica Final

**Sempre comece com dry-run!** Isso evita surpresas e permite validar suas configurações antes de copiar milhares de mensagens.

Boa migração! 🚀
