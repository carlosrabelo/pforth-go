# Arquitetura

Visão geral do interpretador interno, modelo de memória e mecanismo de execução do pforth.

## Pilhas

Duas pilhas LIFO separadas, implementadas como slices Go de `Cell` (int):

| Pilha | Variável | Propósito |
|---|---|---|
| Pilha de Dados (DS) | `f.DS` | Armazenamento de operandos para palavras aritméticas, de memória e lógicas |
| Pilha de Retorno (RS) | `f.RS` | Índices de laços, endereços de retorno de definições de colchetes, armazenamento temporário `>R`/`R>` |

Operações de pilha geram pânico com mensagem descritiva em caso de underflow.

## Modelo de Memória

Array plano de bytes de 65536 elementos (`[]byte`):

- **Células** são valores de 16 bits armazenados em formato little-endian (dois bytes por célula).
- **Ponteiro do Dicionário (DP)** começa em 1024, reservando o primeiro 1KB para uso do sistema.
- **`@` (fetch)** lê 2 bytes e retorna um `Cell`.
- **`!` (store)** escreve um `Cell` como 2 bytes.
- **`C@` / `C!`** lê/escreve bytes individuais.
- **`HERE`** retorna o valor atual do DP.
- **`ALLOT`** avança o DP em n bytes.
- **`,` (comma)** armazena uma célula no DP e avança em 2.
- **`C,` (c-comma)** armazena um byte no DP e avança em 1.

Acesso a bytes/palavras é verificado contra limites e gera pânico em endereços fora da faixa.

## Dicionário

O dicionário é um slice de ponteiros `*Word`, mais um `map[string]int` para busca nome-para-índice:

```go
type Word struct {
    Name      string
    Immediate bool
    Type      WordType    // WordPrimitive, WordColon, WordConstant, WordVariable, WordCreate
    Code      XTCode      // Função Go nativa (apenas primitivas)
    Body      []Cell      // Tokens de execução para definições de colchetes
    PFA       int         // Endereço do campo de parâmetro (VARIABLE/CREATE)
}
```

### Tipos de Palavra

| Tipo | Comportamento |
|---|---|
| `WordPrimitive` | Chama `w.Code(f)` — uma função Go nativa |
| `WordColon` | Salva `f.Body`/`f.IP`, executa `w.Body[]` célula por célula, restaura estado |
| `WordConstant` | Empurra `w.Body[0]` na DS |
| `WordVariable` | Empurra `w.PFA` na DS |
| `WordCreate` | Empurra `w.PFA` na DS |

### Resolução de Nomes

- Todos os nomes são normalizados para maiúsculas via `strings.ToUpper` durante definição e busca.
- `FIND` retorna o XT (índice no slice `f.Words`) e uma flag de sucesso.
- `'` (tick) retorna o XT da próxima palavra analisada.
- `>BODY` converte um XT para seu endereço de campo de parâmetro (XT + 2 offsets de célula).

## Mecanismo de Execução

### Definições de Colchetes

Quando `ExecuteWord` encontra um `WordColon`:

1. Salva o `f.Body` e `f.IP` atuais em uma pilha de chamadas Go (variáveis locais).
2. Define `f.Body = w.Body` e `f.IP = 0`.
3. Laço: lê `f.Body[f.IP]`, incrementa IP, chama `ExecuteWord` no token.
4. Restaura o `f.Body`/`f.IP` anterior.

Esta abordagem é recursiva em Go — cada definição de colchete aninhada adiciona um quadro. O mesmo mecanismo suporta `EXIT` (define IP após o fim → o laço termina, o quadro anterior é restaurado).

### Ramificação

Offsets de ramificação são armazenados inline no body após o XT de ramificação:

- **`BRANCH`**: `f.IP += f.Body[f.IP]` (salto incondicional)
- **`0BRANCH`**: se TOS = 0, `f.IP += f.Body[f.IP]`; senão `f.IP++` (pular offset)

Palavras em tempo de compilação (`IF`, `THEN`, `ELSE`, `BEGIN`, `UNTIL`, `AGAIN`, `WHILE`, `REPEAT`) ajustam offsets forward/backward durante a compilação usando a pilha de dados como rastreador temporário de endereços.

### DO/LOOP

Estado em tempo de execução armazenado na pilha de retorno como pares: `( limit index -- )`.

- **`(DO)`**: desempilha `limit start` da DS, empilha `limit start` na RS.
- **`I`**: copia o índice do topo da RS para a DS.
- **`J`**: lê o índice do próximo par na RS (laços aninhados).
- **`(LOOP)`**: incrementa o índice, se `index < limit` salta de volta pelo offset inline.
- **`(+LOOP)`**: adiciona um valor de passo ao índice, testa terminação do laço.
- **`UNLOOP`**: remove um par da RS (usado por `LEAVE`).

## Estado de Compilação

Quando `STATE` é verdadeiro (`f.State`), tokens são adicionados a `f.compileList` como XTs (tokens de execução) em vez de serem executados imediatamente. Números são compilados como pares `LIT <valor>`.

- **`:`** — cria um novo placeholder `WordColon`, define `STATE` como verdadeiro, registra o nome.
- **`;`** — adiciona `EXIT`, finaliza o body, define `STATE` como falso, reseta a lista de compilação.
- **`IMMEDIATE`** — marca a palavra mais recentemente definida como imediata (executada mesmo em estado de compilação).

## Laço do Interpretador

```
QUIT → (
  se interpretando: imprime "OK"
  lê uma linha para o TIB
  define IN = 0
  interpretLoop()
    analisa um token
    busca no dicionário
    se encontrado:
      se compilando E não imediato: compila XT
      senão: executa palavra
    senão: tenta analisar como número
    se falhou: imprime " ?", abort
)
```

O buffer de entrada de texto (`TIB`) e o ponteiro de entrada (`IN`) suportam re-análise com `WORD` e `EXPECT`.

## REPL

O REPL é controlado por `QUIT`, que faz um laço até `f.Running` se tornar falso (definido por `BYE` ou EOF na stdin). Cada iteração lê uma linha via `bufio.Reader.ReadString('\n')`, portanto não há limite fixo de tamanho de linha.
