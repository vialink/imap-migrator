# Migrador de Contas IMAP

Programa em Go para migrar contas de email IMAP de um servidor para outro, preservando todas as pastas, mensagens e respectivos status.

## Características

- ✅ **Verificação prévia de conexões**: Testa todas as conexões antes de iniciar a migração
- ✅ **Processamento paralelo**: Até 5 migrações simultâneas
- ✅ **Preservação completa**: Mantém estrutura de pastas, mensagens, flags e datas
- ✅ **Logs detalhados**: Acompanhamento completo do processo
- ✅ **Tratamento de erros**: Continua a migração mesmo se uma conta falhar

## Requisitos

- Go 1.22 ou superior
- Acesso aos servidores IMAP de origem e destino (porta 993 - IMAPS/TLS)

## Instalação

1. Clone ou copie os ficheiros para um diretório
2. Instale as dependências:

```bash
go mod tidy
```

## Configuração

Crie um ficheiro `contas.csv` no mesmo diretório do programa com o seguinte formato:

```csv
email_origem,conta_origem,senha_origem,servidor_origem,email_destino,conta_destino,senha_destino,servidor_destino
user1@origem.com,user1,senha123,imap.origem.com,user1@destino.com,user1,senha456,imap.destino.com
user2@origem.com,user2,senha789,imap.origem.com,user2@destino.com,user2,senha012,imap.destino.com
```

### Formato do CSV

- **email_origem**: Endereço de email de origem (usado apenas para logs)
- **conta_origem**: Nome de utilizador para login no servidor de origem
- **senha_origem**: Senha da conta de origem
- **servidor_origem**: Endereço do servidor IMAP de origem (sem porta, usa 993 automaticamente)
- **email_destino**: Endereço de email de destino (usado apenas para logs)
- **conta_destino**: Nome de utilizador para login no servidor de destino
- **senha_destino**: Senha da conta de destino
- **servidor_destino**: Endereço do servidor IMAP de destino (sem porta, usa 993 automaticamente)

## Uso

### Compilar o programa

```bash
go build -o imap-migrator main.go
```

### Executar

```bash
./imap-migrator
```

Ou diretamente sem compilar:

```bash
go run main.go
```

## Processo de Migração

O programa executa em duas fases:

### Fase 1: Verificação de Conexões

- Lê o ficheiro `contas.csv`
- Testa a conexão e autenticação com **todos** os servidores (origem e destino)
- Apresenta um relatório detalhado
- **Só avança para a Fase 2 se todas as conexões forem bem-sucedidas**

### Fase 2: Migração

Para cada conta:
1. Conecta aos servidores de origem e destino
2. Lista todas as pastas da conta de origem
3. Para cada pasta:
   - Cria a pasta no destino (se não existir)
   - Copia todas as mensagens, preservando:
     - Conteúdo completo
     - Flags (lida, não lida, marcada, etc.)
     - Data original da mensagem
4. Regista o progresso e eventuais erros

**Até 5 contas são migradas em paralelo** para acelerar o processo.

## Logs

O programa gera logs detalhados na consola, incluindo:
- Progresso de cada migração
- Número de pastas e mensagens processadas
- Erros encontrados (sem interromper o processo)

Exemplo de log:
```
[INÍCIO MIGRAÇÃO] user1@origem.com -> user1@destino.com
[user1@origem.com] Encontradas 8 pastas para migrar.
[user1@origem.com] Processando pasta: INBOX
[user1@origem.com] Pasta 'INBOX' tem 150 mensagens para copiar.
[user1@origem.com] Processando pasta: Sent
[user1@origem.com] Pasta 'Sent' tem 75 mensagens para copiar.
[FIM MIGRAÇÃO] user1@origem.com -> user1@destino.com
```## 📝 Notas Importantes

- **Duplicados**: A versão atual **não** verifica duplicados. Se executar o programa múltiplas vezes na mesma conta, as mensagens serão copiadas novamente.
- **Conexão segura**: O programa usa apenas conexões TLS (porta 993). Não suporta conexões não encriptadas.
- **Timeout**: Conexões que não respondem em 10 segundos são consideradas falhadas.
- **Pastas especiais**: Pastas marcadas como "não selecionáveis" são ignoradas automaticamente.
- **Quota**: Se a conta de destino ficar cheia, o programa para automaticamente com uma mensagem clara.
- **Gmail**: Totalmente compatível! Veja `GMAIL.md` para instruções específicas.
- **Relatórios**: O programa gera automaticamente um relatório detalhado para cada conta migrada no diretório `relatorios/`.

## 📋 Relatórios de Auditoria

Após cada migração, o programa gera automaticamente um relatório detalhado contendo:

- **Informações gerais**: origem, destino, duração, status
- **Resumo geral**: total de pastas, mensagens copiadas, falhadas e puladas
- **Detalhes por pasta**: estatísticas individuais de cada pasta
- **Lista de erros**: descrição detalhada de todos os erros ocorridos

Os relatórios são guardados no diretório `relatorios/` com o formato:
```
migracao_<email>_<timestamp>.txt
```

**Exemplo:**
```
relatorios/migracao_user_at_origem_com_20251111_200015.txt
```

Veja `EXEMPLO_RELATORIO.txt` para um exemplo completo de relatório.

## Resolução de Problemas

### Erro: "falha ao conectar ao servidor"
- Verifique se o endereço do servidor está correto
- Confirme que a porta 993 está acessível
- Verifique se o servidor suporta TLS/SSL

### Erro: "falha na autenticação"
- Confirme que o utilizador e senha estão corretos
- Alguns servidores exigem "senhas de aplicação" em vez da senha normal (ex: Gmail, Outlook)
- Verifique se a conta tem IMAP ativado

### Erro: "linha X ignorada por não ter 8 colunas"
- Verifique se o ficheiro CSV tem exatamente 8 colunas em cada linha
- Certifique-se de que não há vírgulas extras nos campos

## Melhorias Futuras

- [ ] Prevenção de duplicados usando Message-ID
- [ ] Suporte para conexões não encriptadas (porta 143)
- [ ] Modo de "dry run" (simulação sem copiar)
- [ ] Filtros por data ou pasta específica
- [ ] Retomar migrações interrompidas
- [ ] Interface web para configuração

## Licença

Este programa é fornecido como está, sem garantias. Use por sua conta e risco.
