package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Instância global do NTPConfig
var ntpConfig *NTPConfig

// Middleware de logging para registrar todas as requisições
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
		log.Printf("Requisição processada em %v", time.Since(start))
	}
}

// Handler para a rota principal
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bem-vindo ao Servidor NTP!\n")
	fmt.Fprintf(w, "========================\n\n")
	fmt.Fprintf(w, "Hora do Sistema:  %s\n", time.Now().Format("15:04:05 02/01/2006"))
	fmt.Fprintf(w, "Hora NTP (corrigida): %s\n\n", ntpConfig.GetCorrectedTime().Format("15:04:05 02/01/2006"))
	fmt.Fprintf(w, "Offset: %v\n", ntpConfig.offset)
}

// Handler para rota de status
func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok", "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
}

// Handler para rota de informações
func infoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Informações do Servidor:\n")
	fmt.Fprintf(w, "======================\n")
	fmt.Fprintf(w, "Host: %s\n", r.Host)
	fmt.Fprintf(w, "User-Agent: %s\n", r.UserAgent())
	fmt.Fprintf(w, "Método: %s\n", r.Method)
	fmt.Fprintf(w, "URL: %s\n", r.URL.Path)
}

// Handler para status do NTP (JSON)
func ntpStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := ntpConfig.GetStatus()

	jsonData, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		http.Error(w, "Erro ao gerar JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

// Handler para forçar sincronização manual
func ntpSyncHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido. Use POST", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Sincronização manual solicitada")
	err := ntpConfig.SyncTime()

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":         "success",
		"message":        "Sincronização realizada com sucesso",
		"corrected_time": ntpConfig.GetCorrectedTime().Format(time.RFC3339),
		"offset":         ntpConfig.offset.String(),
	})
}

// Handler para obter o horário atual corrigido
func ntpTimeHandler(w http.ResponseWriter, r *http.Request) {
	correctedTime := ntpConfig.GetCorrectedTime()
	systemTime := time.Now()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ntp_time":       correctedTime.Format(time.RFC3339),
		"system_time":    systemTime.Format(time.RFC3339),
		"offset":         ntpConfig.offset.String(),
		"unix_timestamp": correctedTime.Unix(),
	})
}

func main() {
	// Inicializar configuração NTP
	log.Printf("=================================")
	log.Printf("Inicializando configuração NTP...")
	log.Printf("=================================")

	ntpConfig = NewNTPConfig()
	ntpConfig.StartAutoSync()

	// Configurar rotas com middleware de logging
	http.HandleFunc("/", loggingMiddleware(homeHandler))
	http.HandleFunc("/status", loggingMiddleware(statusHandler))
	http.HandleFunc("/info", loggingMiddleware(infoHandler))

	// Rotas NTP
	http.HandleFunc("/ntp/status", loggingMiddleware(ntpStatusHandler))
	http.HandleFunc("/ntp/sync", loggingMiddleware(ntpSyncHandler))
	http.HandleFunc("/ntp/time", loggingMiddleware(ntpTimeHandler))

	// Configurações do servidor
	port := ":8080"

	log.Printf("\n=================================")
	log.Printf("Servidor iniciando na porta %s", port)
	log.Printf("=================================")
	log.Printf("Rotas disponíveis:")
	log.Printf("\n  Rotas Gerais:")
	log.Printf("    - http://localhost%s/", port)
	log.Printf("    - http://localhost%s/status", port)
	log.Printf("    - http://localhost%s/info", port)
	log.Printf("\n  Rotas NTP:")
	log.Printf("    - http://localhost%s/ntp/status (GET)", port)
	log.Printf("    - http://localhost%s/ntp/time (GET)", port)
	log.Printf("    - http://localhost%s/ntp/sync (POST)", port)
	log.Printf("=================================\n")

	// Iniciar servidor
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
