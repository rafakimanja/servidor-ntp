package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

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
	fmt.Fprintf(w, "Bem-vindo ao Servidor Web Simples!\n")
	fmt.Fprintf(w, "==================================\n\n")
	fmt.Fprintf(w, "Hora atual: %s\n", time.Now().Format("15:04:05 02/01/2006"))
	fmt.Fprintf(w, "Timezone: %s\n", time.Now().Location().String())
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

func main() {
	// Configurar fuso horário para América/São Paulo
	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Printf("Aviso: Não foi possível carregar timezone 'America/Sao_Paulo': %v", err)
		log.Printf("Usando timezone do sistema")
	} else {
		time.Local = location
		log.Printf("Timezone configurado: %s", time.Local.String())
	}

	// Configurar rotas com middleware de logging
	http.HandleFunc("/", loggingMiddleware(homeHandler))
	http.HandleFunc("/status", loggingMiddleware(statusHandler))
	http.HandleFunc("/info", loggingMiddleware(infoHandler))

	// Configurações do servidor
	port := ":8080"

	log.Printf("\n=================================")
	log.Printf("Servidor iniciando na porta %s", port)
	log.Printf("=================================")
	log.Printf("Rotas disponíveis:")
	log.Printf("- http://localhost%s/", port)
	log.Printf("- http://localhost%s/status", port)
	log.Printf("- http://localhost%s/info", port)

	// Iniciar servidor
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
