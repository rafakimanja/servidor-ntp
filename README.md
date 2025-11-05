# Servidor Web com Sincronização NTP

Servidor web simples em Go com funcionalidades de sincronização de horário via NTP (Network Time Protocol).

## 🚀 Funcionalidades

- **Servidor HTTP** na porta 8080
- **Sincronização automática** com servidores NTP brasileiros
- **Logging detalhado** de todas as requisições
- **API REST** para gerenciamento de horário
- **Múltiplos servidores NTP** com fallback automático

## 📋 Pré-requisitos

- Go 1.16 ou superior
- Acesso à internet para sincronização NTP

## 🔧 Instalação

```bash
# Clonar ou criar o projeto
cd ntp-server

# Inicializar o módulo Go (se ainda não foi feito)
go mod init teste-ntp

# Instalar dependências
go get github.com/beevik/ntp

# Baixar todas as dependências
go mod tidy
```

## ▶️ Execução

```bash
# Executar diretamente
go run main.go ntp_config.go

# Ou compilar e executar
go build -o servidor
./servidor
```

## 🌐 Rotas Disponíveis

### Rotas Gerais

- **GET /** - Página inicial com informações de horário
- **GET /status** - Status do servidor (JSON)
- **GET /info** - Informações da requisição

### Rotas NTP

- **GET /ntp/status** - Status completo da sincronização NTP
- **GET /ntp/time** - Horário atual corrigido pelo NTP
- **POST /ntp/sync** - Força sincronização manual com servidores NTP

## 📡 Servidores NTP Configurados

O sistema utiliza os seguintes servidores NTP brasileiros:
- 0.br.pool.ntp.org
- 1.br.pool.ntp.org
- 2.br.pool.ntp.org
- a.st1.ntp.br
- b.st1.ntp.br

## 🔄 Sincronização Automática

O servidor sincroniza automaticamente a cada **10 minutos** com os servidores NTP.
A primeira sincronização ocorre imediatamente ao iniciar o servidor.

## 📊 Exemplos de Uso

### Verificar status do NTP

```bash
curl http://localhost:8080/ntp/status
```

Resposta:
```json
{
  "available_servers": [
    "0.br.pool.ntp.org",
    "1.br.pool.ntp.org",
    "2.br.pool.ntp.org",
    "a.st1.ntp.br",
    "b.st1.ntp.br"
  ],
  "corrected_time": "10:23:45 05/11/2025",
  "current_server": "0.br.pool.ntp.org",
  "is_running": true,
  "last_sync": "2025-11-05T10:23:30-03:00",
  "offset": "123.456ms",
  "sync_interval": "10m0s",
  "system_time": "10:23:45 05/11/2025",
  "time_since_sync": "15.234s"
}
```

### Obter horário corrigido

```bash
curl http://localhost:8080/ntp/time
```

Resposta:
```json
{
  "ntp_time": "2025-11-05T10:23:45-03:00",
  "system_time": "2025-11-05T10:23:45-03:00",
  "offset": "123.456ms",
  "unix_timestamp": 1730811825
}
```

### Forçar sincronização manual

```bash
curl -X POST http://localhost:8080/ntp/sync
```

Resposta:
```json
{
  "status": "success",
  "message": "Sincronização realizada com sucesso",
  "corrected_time": "2025-11-05T10:23:45-03:00",
  "offset": "123.456ms"
}
```

## 📝 Logs

O servidor registra automaticamente:
- Inicialização do servidor
- Sincronizações NTP (automáticas e manuais)
- Todas as requisições HTTP com timestamp
- Tempo de processamento de cada requisição
- Erros de sincronização

Exemplo de logs:
```
2025/11/05 10:23:30 =================================
2025/11/05 10:23:30 Inicializando configuração NTP...
2025/11/05 10:23:30 =================================
2025/11/05 10:23:30 Iniciando sincronização automática (intervalo: 10m0s)
2025/11/05 10:23:30 ✓ Sincronizado com 0.br.pool.ntp.org
2025/11/05 10:23:30   Offset: 123.456ms
2025/11/05 10:23:30   Stratum: 2
2025/11/05 10:23:30   Precisão: 1ms
2025/11/05 10:23:30   RTT: 45.234ms
2025/11/05 10:23:45 [GET] /ntp/status 127.0.0.1:42292
2025/11/05 10:23:45 Requisição processada em 234.567µs
```

## 🏗️ Estrutura do Projeto

```
ntp-server/
├── main.go          # Servidor HTTP e rotas
├── ntp_config.go    # Módulo de configuração NTP
├── go.mod           # Dependências do projeto
├── go.sum           # Checksums das dependências
└── README.md        # Este arquivo
```

## 🔒 Observação sobre Permissões

**Nota**: Este servidor **não altera** o horário do sistema operacional. Ele apenas:
- Obtém o horário correto dos servidores NTP
- Calcula o offset (diferença) entre o horário do sistema e o horário NTP
- Fornece o horário corrigido através da API

Para alterar o horário do sistema Linux, você precisaria de permissões root e usar comandos como `timedatectl` ou `date`.

## 🛠️ Desenvolvimento

### Adicionar novo servidor NTP

Edite o arquivo `ntp_config.go` e adicione o servidor no array `servers` da função `NewNTPConfig()`.

### Alterar intervalo de sincronização

Modifique o valor de `syncInterval` em `NewNTPConfig()` no arquivo `ntp_config.go`.

## 📄 Licença

Este projeto é de uso livre para fins educacionais e comerciais.

