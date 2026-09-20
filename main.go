package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func readPasswordFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text()), nil
	}
	return "", scanner.Err()
}

func prompt(msg, def string) string {
	if def != "" {
		fmt.Printf("%s [%s]: ", msg, def)
	} else {
		fmt.Printf("%s: ", msg)
	}
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return def
	}
	text := strings.TrimSpace(scanner.Text())
	if text == "" {
		return def
	}
	return text
}

const numRegras = 3

func editarRegras(page *rod.Page, ip string, n int) {
	for i := 0; i < n; i++ {
		page.MustElement(fmt.Sprintf(`span.table-grid-icon.edit-modify-icon#edit_%d`, i)).MustClick()
		page.MustWaitStable()

		f := page.MustElement(`input#ipAddr`)
		f.MustSelectAllText()
		f.MustInput(ip)

		page.MustElement(`button#ok`).MustClick()
		page.MustWaitStable()
	}
}

func main() {
	debug := flag.Bool("debug", false, "Abre navegador visível")
	slow := flag.Bool("slow", false, "Executa devagar (1s entre ações)")
	interactive := flag.Bool("i", false, "Modo interativo: pergunta cada valor")
	passwordFile := flag.String("password-file", "", "Caminho do arquivo com a senha (padrão: $TPLINK_PASSWORD_FILE ou senha-wifi.txt)")
	flag.Parse()

	args := flag.Args()
	var ip string
	if *interactive {
		fmt.Println("Modo interativo — vai perguntar o IP durante a execução.")
	} else if len(args) < 1 {
		fmt.Println("Uso: ./tplink-ipv6-set [-i] [-debug] [-slow] [-password-file CAMINHO] <ip>")
		fmt.Println("  -i               modo interativo (pergunta o IP)")
		fmt.Println("  -debug           navegador visível")
		fmt.Println("  -password-file   arquivo com a senha do roteador")
		os.Exit(1)
	} else {
		ip = args[0]
	}

	// Resolve o arquivo de senha: flag > env > padrão local.
	file := *passwordFile
	if file == "" {
		file = os.Getenv("TPLINK_PASSWORD_FILE")
	}
	if file == "" {
		file = "senha-wifi.txt"
	}

	password, err := readPasswordFromFile(file)
	if err != nil {
		fmt.Println("Erro ao ler a senha:", err)
		os.Exit(1)
	}

	chromium, err := exec.LookPath("chromium")
	if err != nil {
		fmt.Println("chromium não encontrado no PATH. Rode: nix-shell -p chromium")
		os.Exit(1)
	}
	u := launcher.New().
		Bin(chromium).
		Headless(!*debug && !*interactive).
		MustLaunch()

	browser := rod.New().
		ControlURL(u).
		MustConnect()
	defer browser.MustClose()

	if *slow {
		browser.SlowMotion(1 * time.Second)
	}
	page := browser.MustPage("http://192.168.0.1")
	page.MustWaitStable()

	// --- LOGIN ---
	page.MustElement(`input#pc-login-password`).MustInput(password)
	page.MustElement(`span.text.button-text`).MustClick()

	// Poll síncrono por 3s: modal "outra sessão ativa"
	for i := 0; i < 15; i++ {
		time.Sleep(200 * time.Millisecond)
		has, err := page.Eval(`() => !!document.querySelector('#confirm-yes')`)
		if err == nil && has.Value.Bool() {
			page.MustElement(`#confirm-yes`).MustClick()
			break
		}
	}
	page.MustWaitStable()

	// --- AVANÇADO ---
	page.MustElement(`span.T_adv.text`).MustClick()
	page.MustWaitStable()

	// --- SEGURANÇA ---
	page.MustElement(`a[url="portFiltering.htm"]`).MustClick()
	page.MustWaitStable()

	// --- FIREWALL IPV6 ---
	page.MustElement(`a[url="ipv6Firewall.htm"]`).MustClick()
	page.MustWaitStable()

	// --- EDITAR REGRAS ---
	if *interactive {
		ip = prompt("IP para as regras", "")
	}
	editarRegras(page, ip, numRegras)

	// --- LOGOUT ---
	page.MustElement(`a#topLogout`).MustClick()
	page.MustWaitStable()
	page.MustElement(`button.button-button.green.pure-button.btn-msg.btn-msg-ok.btn-confirm`).MustClick()
	page.MustWaitStable()

	if *interactive {
		fmt.Println("Concluído. Pressione Enter para fechar.")
		prompt("", "")
	}
}
