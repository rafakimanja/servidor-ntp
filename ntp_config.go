package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/beevik/ntp"
)

// NTPConfig gerencia a configuração e sincronização NTP
type NTPConfig struct {
	mu            sync.RWMutex
	servers       []string
	currentServer string
	lastSync      time.Time
	offset        time.Duration
	syncInterval  time.Duration
	isRunning     bool
	stopChan      chan struct{}
}

// NewNTPConfig cria uma nova instância de configuração NTP
func NewNTPConfig() *NTPConfig {
	return &NTPConfig{
		servers: []string{
			"0.br.pool.ntp.org",
			"1.br.pool.ntp.org",
			"2.br.pool.ntp.org",
			"a.st1.ntp.br",
			"b.st1.ntp.br",
		},
		currentServer: "0.br.pool.ntp.org",
		syncInterval:  time.Minute * 10, // Sincroniza a cada 10 minutos
		stopChan:      make(chan struct{}),
	}
}

// GetNTPTime obtém o horário de um servidor NTP
func (n *NTPConfig) GetNTPTime(server string) (time.Time, error) {
	response, err := ntp.Query(server)
	if err != nil {
		return time.Time{}, fmt.Errorf("erro ao consultar NTP %s: %w", server, err)
	}

	ntpTime := time.Now().Add(response.ClockOffset)
	return ntpTime, nil
}

// SyncTime sincroniza o horário com o servidor NTP
func (n *NTPConfig) SyncTime() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	var lastErr error

	// Tenta cada servidor até conseguir
	for _, server := range n.servers {
		response, err := ntp.Query(server)
		if err != nil {
			lastErr = err
			log.Printf("Falha ao consultar %s: %v", server, err)
			continue
		}

		n.currentServer = server
		n.lastSync = time.Now()
		n.offset = response.ClockOffset

		log.Printf("✓ Sincronizado com %s", server)
		log.Printf("  Offset: %v", response.ClockOffset)
		log.Printf("  Stratum: %d", response.Stratum)
		log.Printf("  Precisão: %v", response.Precision)
		log.Printf("  RTT: %v", response.RTT)

		return nil
	}

	return fmt.Errorf("falha ao sincronizar com todos os servidores: %w", lastErr)
}

// StartAutoSync inicia a sincronização automática
func (n *NTPConfig) StartAutoSync() {
	n.mu.Lock()
	if n.isRunning {
		n.mu.Unlock()
		return
	}
	n.isRunning = true
	n.mu.Unlock()

	log.Printf("Iniciando sincronização automática (intervalo: %v)", n.syncInterval)

	// Sincroniza imediatamente
	if err := n.SyncTime(); err != nil {
		log.Printf("Erro na sincronização inicial: %v", err)
	}

	// Inicia goroutine para sincronização periódica
	go func() {
		ticker := time.NewTicker(n.syncInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := n.SyncTime(); err != nil {
					log.Printf("Erro na sincronização automática: %v", err)
				}
			case <-n.stopChan:
				log.Printf("Sincronização automática parada")
				return
			}
		}
	}()
}

// StopAutoSync para a sincronização automática
func (n *NTPConfig) StopAutoSync() {
	n.mu.Lock()
	defer n.mu.Unlock()

	if !n.isRunning {
		return
	}

	close(n.stopChan)
	n.isRunning = false
	n.stopChan = make(chan struct{})
}

// GetCorrectedTime retorna o horário corrigido com o offset NTP
func (n *NTPConfig) GetCorrectedTime() time.Time {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return time.Now().Add(n.offset)
}

// GetStatus retorna o status atual da configuração NTP
func (n *NTPConfig) GetStatus() map[string]interface{} {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return map[string]interface{}{
		"current_server":    n.currentServer,
		"last_sync":         n.lastSync.Format(time.RFC3339),
		"time_since_sync":   time.Since(n.lastSync).String(),
		"offset":            n.offset.String(),
		"sync_interval":     n.syncInterval.String(),
		"is_running":        n.isRunning,
		"available_servers": n.servers,
		"corrected_time":    n.GetCorrectedTime().Format("15:04:05 02/01/2006"),
		"system_time":       time.Now().Format("15:04:05 02/01/2006"),
	}
}

// AddServer adiciona um novo servidor NTP à lista
func (n *NTPConfig) AddServer(server string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.servers = append(n.servers, server)
	log.Printf("Servidor NTP adicionado: %s", server)
}

// SetSyncInterval define o intervalo de sincronização
func (n *NTPConfig) SetSyncInterval(interval time.Duration) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.syncInterval = interval
	log.Printf("Intervalo de sincronização atualizado: %v", interval)
}
