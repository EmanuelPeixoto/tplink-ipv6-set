# tplink-ipv6-set

Configura regras de firewall IPv6 no roteador **TP-Link EX141 (C6BF)** via interface web.

Automatiza login → Avançado → Segurança → Firewall IPv6 → editar regras → logout.

## Requisitos

- [Nix](https://nixos.org/) (ou Chromium instalado no PATH)
- Go 1.26+

## Instalação

```bash
git clone https://github.com/seu-user/tplink-ipv6-set
cd tplink-ipv6-set
echo "sua-senha" > senha-wifi.txt
go build -o tplink-ipv6-set .
```

No NixOS, rode dentro do `nix-shell`:

```bash
nix-shell -p chromium
```

## Uso

```bash
# Modo automático (passa o IP como argumento)
./tplink-ipv6-set "2001:db8::1"

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

## Configuração

### Senha

Coloque a senha do roteador no arquivo `senha-wifi.txt` (uma linha).

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
