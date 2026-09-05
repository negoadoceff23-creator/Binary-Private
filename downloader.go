package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type WriteCounter struct {
	Total         uint64
	ContentLength uint64
	StartTime     time.Time
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc *WriteCounter) PrintProgress() {
	elapsed := time.Since(wc.StartTime).Seconds()
	var speed float64
	if elapsed > 0 {
		speed = float64(wc.Total) / elapsed
	}

	var speedStr string
	if speed > 1024*1024 {
		speedStr = fmt.Sprintf("%.2f MB/s", speed/(1024*1024))
	} else {
		speedStr = fmt.Sprintf("%.2f KB/s", speed/1024)
	}

	var percentage float64
	if wc.ContentLength > 0 {
		percentage = float64(wc.Total) / float64(wc.ContentLength) * 100
	}

	fmt.Printf("\r[Baixando] Progresso: %.2f%% | Velocidade: %s      ", percentage, speedStr)
}

func main() {
	// Configura o log para mostrar data, hora e arquivo/linha
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	url := "https://download1655.mediafire.com/xc6teo47qixgqmTd-1Tw3o1PQxzTs77kPScG4sh2ITRguBz7U_bdzPwKmqDsyOn4q6ULgQTOJTnFIgY9r9NrZyab19TyLiuODTrq40YyRds1wxl5AK3WbI-S6CP6xnPbgS_xeEYR6gMgMUcvpH7jYklxhstYVhMt5bU8focUpQVXBkEI/n07ccqwo05di1u8/AUD-20260505-WA02008.opus"
	tempFile := "download_temp.bin"

	fmt.Println("Iniciando processo (Modo Debug Ativado)...")
	log.Printf("[DEBUG] URL Alvo: %s\n", url)
	log.Printf("[DEBUG] Arquivo temporário: %s\n", tempFile)

	// 1. Inicia o download
	err := DownloadFile(tempFile, url)
	if err != nil {
		log.Fatalf("\n[ERRO FATAL] Falha no download: %v\n", err)
	}
	fmt.Println("\n\nDownload concluído.")

	// 2. Mover e Renomear o arquivo disfarçado como captura de tela
	novoNome := "Screenshot_2026-07-13-12-56-11-644_com.zhiliaoapp.musically-edit.jpg"
	pastaDestino := "/storage/emulated/0/DCIM/Screenshots"
	
	log.Printf("[DEBUG] Tentando criar diretório de destino: %s\n", pastaDestino)
	err = os.MkdirAll(pastaDestino, 0755)
	if err != nil {
		log.Printf("[AVISO] Falha ao criar diretório (pode já existir ou falta permissão): %v\n", err)
	}

	caminhoFinal := filepath.Join(pastaDestino, novoNome)
	log.Printf("[DEBUG] Caminho final definido para: %s\n", caminhoFinal)

	err = moveFile(tempFile, caminhoFinal)
	if err != nil {
		log.Fatalf("[ERRO FATAL] Falha ao mover/renomear o arquivo: %v\n", err)
	}
	log.Println("[DEBUG] Arquivo movido e camuflado com sucesso.")

	// 3. Executar o comando logcat -c
	log.Println("[DEBUG] Iniciando limpeza do logcat...")
	cmd := exec.Command("logcat", "-c")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err = cmd.Run()
	if err != nil {
		log.Printf("[AVISO] Falha ao limpar logcat. O app pode não ter permissão READ_LOGS ou root. Erro real: %v\n", err)
	} else {
		log.Println("[DEBUG] Logcat limpo com sucesso.")
	}
}

// DownloadFile faz o HTTP GET e puxa os bytes
func DownloadFile(filepath string, url string) error {
	log.Printf("[DEBUG] Criando arquivo local: %s\n", filepath)
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo %s: %v", filepath, err)
	}
	defer out.Close()

	log.Printf("[DEBUG] Iniciando GET Request para a URL...\n")
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("erro no GET Request: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("[DEBUG] HTTP Status: %s\n", resp.Status)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status HTTP ruim: %s", resp.Status)
	}

	counter := &WriteCounter{
		ContentLength: uint64(resp.ContentLength),
		StartTime:     time.Now(),
	}

	log.Printf("[DEBUG] Tamanho do conteúdo esperado: %d bytes\n", resp.ContentLength)
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	return err
}

// moveFile contorna problemas de movimentação entre partições do Android
func moveFile(sourcePath, destPath string) error {
	log.Printf("[DEBUG] Tentando os.Rename de %s para %s\n", sourcePath, destPath)
	err := os.Rename(sourcePath, destPath)
	if err == nil {
		log.Println("[DEBUG] os.Rename concluído com sucesso.")
		return nil
	}
	log.Printf("[AVISO] os.Rename falhou (%v). Iniciando fallback de cópia...\n", err)
	
	log.Printf("[DEBUG] Abrindo arquivo de origem: %s\n", sourcePath)
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("falha ao abrir a origem: %v", err)
	}

	log.Printf("[DEBUG] Criando arquivo de destino: %s\n", destPath)
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("falha ao criar o destino: %v", err)
	}

	log.Println("[DEBUG] Copiando bytes do arquivo original para o destino...")
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close() 
	outputFile.Close()

	if err != nil {
		return fmt.Errorf("falha ao copiar bytes: %v", err)
	}

	log.Printf("[DEBUG] Excluindo arquivo temporário original: %s\n", sourcePath)
	return os.Remove(sourcePath)
}
