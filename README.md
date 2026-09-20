# tplink-ipv6-set

Configura regras de firewall IPv6 no roteador **TP-Link EX141 (C6BF)** via interface web.

Automatiza login → Avançado → Segurança → Firewall IPv6 → editar regras → logout.

## Requisitos

- [Nix](https://nixos.org/) (ou Chromium instalado no PATH)
- Go 1.26+

## Instalação

```bash
git clone https://github.com/EmanuelPeixoto/tplink-ipv6-set
cd tplink-ipv6-set
go build -o tplink-ipv6-set .
```

> A senha do roteador **não** deve ficar dentro do repositório. Coloque-a em um
> arquivo fora do repo (ex.: `~/.config/senha-wifi.txt` com `chmod 600`) e
> informe o caminho via `-password-file` ou `TPLINK_PASSWORD_FILE`.

No NixOS, rode dentro do `nix-shell`:

```bash
nix-shell -p chromium
```

## Uso

```bash
# Modo automático (passa o IP como argumento)
./tplink-ipv6-set "2001:db8::1"

# Modo automático com caminho da senha
./tplink-ipv6-set -password-file ~/.config/senha-wifi.txt "2001:db8::1"

# Ou via variável de ambiente
TPLINK_PASSWORD_FILE=~/.config/senha-wifi.txt ./tplink-ipv6-set "2001:db8::1"

# Modo interativo (pergunta o IP)
./tplink-ipv6-set -i

# Devagar e visível (debug)
./tplink-ipv6-set -debug -slow "2001:db8::1"
```

### Flags

| Flag | Descrição |
|------|-----------|
| `-i` | Modo interativo: pergunta o IP durante a execução |
| `-debug` | Abre o navegador visível |
| `-slow` | Pausa de 1s entre cada ação |
| `-password-file` | Caminho do arquivo com a senha (padrão: `$TPLINK_PASSWORD_FILE` ou `senha-wifi.txt`) |

## Configuração

### Senha

Coloque a senha do roteador em um arquivo local fora do repositório (uma linha)
com permissões restritas:

```bash
install -m 600 /dev/null ~/.config/senha-wifi.txt
echo "sua-senha" > ~/.config/senha-wifi.txt
```

Depois passe o caminho com `-password-file` ou defina `TPLINK_PASSWORD_FILE`.

### Número de regras

Altere a constante no `main.go`:

```go
const numRegras = 3
```

O seletor de cada regra segue o padrão `#edit_0`, `#edit_1`, etc.

## Fluxo

1. Acessa `http://192.168.0.1`
2. Login com senha
3. Clica em **Avançado**
4. Expande **Segurança**
5. Clica em **Firewall IPv6**
6. Para cada regra: clica no ícone de editar, preenche `input#ipAddr`, clica OK
7. Logout

## Modelo testado

- **TP-Link EX141 (C6BF)** — firmware stock
