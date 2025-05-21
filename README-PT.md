# pforth

Interpretador e compilador FORTH implementado em Go, portado das implementações originais em assembly Z80 e MOS 6502.

## Destaques

- Interpretador/compilador Forth com REPL interativo
- Definições com dois-pontos com suporte completo ao compilador (IF/THEN/ELSE, BEGIN/UNTIL/WHILE/REPEAT, DO/LOOP)
- Modelo de memória de 64KB endereçável por byte com pilhas de dados e retorno separadas
- Definições recursivas via RECURSE
- Saída de strings com ."
- Mais de 80 palavras primitivas cobrindo pilha, memória, aritmética, lógica, E/S e operações do compilador
- Dicionário em duas camadas: primitivas em Go e palavras Forth carregadas de `libs/core.fs`
- Programas de demonstração para ordenação, fatorial, Fibonacci, crivo e Torre de Hanói

## Pré-requisitos

- **Go 1.23+** — necessário para compilar a partir do código fonte; [download](https://go.dev/dl/)

## Instalação

### Compilar a partir do código fonte

```bash
git clone https://github.com/carlosrabelo/pforth.git
cd pforth-go
make build
```

Instalar em `~/.local/bin` (padrão), ou em todo o sistema em `/usr/local/bin` (sudo apenas para a cópia):

```bash
make install
make install-system
make uninstall
make uninstall-system
```

## Uso

```bash
make run
```

Ou diretamente:

```bash
./bin/pforth
```

Carregar um programa e entrar no REPL:

```bash
./bin/pforth demos/hello.fs
```

### Exemplo de sessão

```forth
1 2 + .
3 ok

: QUADRADO DUP * ;
ok
5 QUADRADO .
25 ok

: FATORIAL
  DUP 0= IF DROP 1 ELSE DUP 1- RECURSE * THEN ;
ok
10 FATORIAL .
3628800 ok
```

## Documentação

- [Arquitetura](docs/architecture-PT.md) — interpretador interno, modelo de memória, mecanismo de execução
- [Referência de Palavras](docs/words-PT.md) — referência completa de todas as palavras embutidas

## Estrutura do Projeto

```
pforth/              # Código fonte Go
├── cmd/pforth/      # Ponto de entrada
└── internal/forth/  # Motor principal: pilhas, memória, dicionário, primitivas
libs/                # Bibliotecas Forth (core.fs)
demos/               # Programas Forth de exemplo
docs/                # Arquitetura e referência de palavras
bin/                 # Binários compilados (ignorados pelo git)
.make/               # Scripts de build e instalação
```

## Desenvolvimento

```bash
make build             # Compila o binário para bin/pforth
make run               # Compila e executa o REPL
make test              # Executa todos os testes
make quality           # Formata, verifica e faz lint
make install           # Instala o binário em ~/.local/bin
make install-system    # Instala o binário em /usr/local/bin
make uninstall         # Remove de ~/.local/bin
make uninstall-system  # Remove de /usr/local/bin
```

## Licença

Este projeto está licenciado sob a licença MIT — consulte [LICENSE](LICENSE) para detalhes.
