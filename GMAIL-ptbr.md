# Compatibilidade com Gmail

Este documento explica como usar o migrador IMAP com contas do Gmail, tanto como origem quanto como destino.

## 🔐 Requisitos para Gmail

### 1. Ativar IMAP no Gmail

Antes de usar o migrador, é necessário ativar o acesso IMAP:

1. Aceda às **Configurações** do Gmail
2. Vá para a aba **Encaminhamento e POP/IMAP**
3. Na seção **Acesso IMAP**, selecione **Ativar IMAP**
4. Clique em **Guardar alterações**

### 2. Senhas de Aplicação (App Passwords)

O Gmail **não permite** usar a senha normal da conta para acesso IMAP. É necessário criar uma **Senha de Aplicação**:

#### Passo a Passo:

1. Aceda a [myaccount.google.com](https://myaccount.google.com)
2. Vá para **Segurança**
3. Em "Como fazer login no Google", selecione **Verificação em duas etapas**
   - **Importante:** A verificação em duas etapas DEVE estar ativada
4. Role até o final e clique em **Senhas de aplicação**
5. Selecione:
   - **App:** Correio
   - **Dispositivo:** Outro (personalizado)
   - Digite: "Migrador IMAP"
6. Clique em **Gerar**
7. **Copie a senha gerada** (16 caracteres sem espaços)

⚠️ **Esta senha só será mostrada uma vez!** Guarde-a em local seguro.

## 📋 Configuração no CSV

### Exemplo com Gmail como Origem

```csv
email_origem,conta_origem,senha_origem,servidor_origem,email_destino,conta_destino,senha_destino,servidor_destino
user@gmail.com,user@gmail.com,abcdefghijklmnop,imap.gmail.com,user@destino.com,user,senha123,imap.destino.com
```

**Notas:**
- `conta_origem`: Use o **email completo** (user@gmail.com)
- `senha_origem`: Use a **Senha de Aplicação** de 16 caracteres (sem espaços)
- `servidor_origem`: Use `imap.gmail.com`

### Exemplo com Gmail como Destino

```csv
email_origem,conta_origem,senha_origem,servidor_origem,email_destino,conta_destino,senha_destino,servidor_destino
user@origem.com,user,senha123,imap.origem.com,user@gmail.com,user@gmail.com,abcdefghijklmnop,imap.gmail.com
```

**Notas:**
- `conta_destino`: Use o **email completo** (user@gmail.com)
- `senha_destino`: Use a **Senha de Aplicação** de 16 caracteres (sem espaços)
- `servidor_destino`: Use `imap.gmail.com`

### Exemplo com Gmail em Ambos os Lados

```csv
email_origem,conta_origem,senha_origem,servidor_origem,email_destino,conta_destino,senha_destino,servidor_destino
user1@gmail.com,user1@gmail.com,abcdefghijklmnop,imap.gmail.com,user2@gmail.com,user2@gmail.com,qrstuvwxyzabcdef,imap.gmail.com
```

## 📁 Estrutura de Pastas do Gmail

O Gmail usa uma estrutura de pastas diferente dos servidores IMAP tradicionais:

### Pastas Padrão do Gmail

| Nome IMAP | Descrição |
|-----------|-----------|
| `INBOX` | Caixa de entrada |
| `[Gmail]/All Mail` | Todos os emails |
| `[Gmail]/Drafts` | Rascunhos |
| `[Gmail]/Sent Mail` | Enviados |
| `[Gmail]/Spam` | Spam |
| `[Gmail]/Starred` | Com estrela |
| `[Gmail]/Trash` | Lixeira |

### ⚠️ Comportamento Especial do Gmail

1. **All Mail contém TUDO**: A pasta `[Gmail]/All Mail` contém todas as mensagens, incluindo as que estão em outras pastas. Isto pode causar duplicação se não tiver cuidado.

2. **Labels vs Pastas**: O Gmail usa "labels" (etiquetas) em vez de pastas tradicionais. Uma mensagem pode ter múltiplas labels, mas no IMAP aparece como se estivesse em múltiplas pastas.

3. **Arquivamento**: Quando arquiva um email no Gmail, ele sai da INBOX mas permanece em `[Gmail]/All Mail`.

## 🎯 Recomendações para Migração

### Gmail como Origem

**Opção 1: Migrar Tudo (incluindo All Mail)**
- Migra todas as pastas, incluindo `[Gmail]/All Mail`
- ✅ Garante que nada seja perdido
- ❌ Pode criar duplicados no destino

**Opção 2: Excluir All Mail**
- Migra apenas as pastas específicas (INBOX, Sent, etc.)
- ✅ Evita duplicados
- ❌ Pode perder emails arquivados

### Gmail como Destino

**Cuidados:**
1. **Quota**: Contas gratuitas do Gmail têm 15 GB compartilhados (Gmail + Drive + Fotos)
2. **Limites de taxa**: O Gmail pode limitar o número de operações IMAP por segundo
3. **Estrutura de pastas**: Pastas personalizadas serão criadas como labels

## 🔧 Configurações Recomendadas

### Para Migração Rápida

Se estiver a migrar de/para Gmail e tiver muitas mensagens:

1. **Desative a verificação em duas etapas temporariamente** (opcional, mas pode ajudar)
2. **Use uma conexão de internet estável e rápida**
3. **Execute o programa em horários de baixo tráfego**

### Limites do Gmail

| Limite | Valor |
|--------|-------|
| Tamanho máximo de mensagem | 25 MB (com anexos) |
| Quota total (conta gratuita) | 15 GB |
| Quota total (Google Workspace) | 30 GB - ilimitado (depende do plano) |

## 🐛 Resolução de Problemas

### Erro: "Authentication failed"

**Causa:** Senha incorreta ou senha de aplicação não configurada.

**Solução:**
1. Verifique se a verificação em duas etapas está ativada
2. Gere uma nova senha de aplicação
3. Use o email completo como nome de utilizador

### Erro: "IMAP access is disabled"

**Causa:** IMAP não está ativado na conta.

**Solução:**
1. Aceda às configurações do Gmail
2. Ative o acesso IMAP
3. Aguarde alguns minutos e tente novamente

### Erro: "Quota exceeded"

**Causa:** A conta do Gmail está cheia.

**Solução:**
1. Liberte espaço na conta (apague emails grandes, esvazie o lixo)
2. Ou faça upgrade para um plano com mais espaço

### Migração Lenta

**Causa:** Limites de taxa do Gmail.

**Solução:**
1. O programa já tem reconexão automática
2. Seja paciente - migrações grandes podem demorar horas
3. O Gmail pode desacelerar temporariamente se detetar muita atividade

## ✅ Checklist Antes de Migrar

- [ ] IMAP ativado na conta Gmail
- [ ] Verificação em duas etapas ativada
- [ ] Senha de aplicação gerada e guardada
- [ ] Espaço suficiente na conta de destino
- [ ] Ficheiro `contas.csv` configurado corretamente
- [ ] Teste com uma conta pequena primeiro

## 📞 Suporte

Se encontrar problemas específicos do Gmail:
1. Verifique os [logs de atividade da conta](https://myaccount.google.com/notifications)
2. Consulte a [documentação oficial do Gmail IMAP](https://support.google.com/mail/answer/7126229)
3. Verifique se há alertas de segurança na sua conta Google
