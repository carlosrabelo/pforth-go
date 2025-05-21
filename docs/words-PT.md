# Referência de Palavras

Referência completa de todas as palavras embutidas do pforth, organizadas por categoria. A notação de pilha segue a convenção Forth `( antes -- depois )` onde o topo da pilha está à direita.

## Manipulação de Pilha

| Palavra | Efeito na Pilha | Descrição |
|---|---|---|
| `DUP` | `( n -- n n )` | Duplica o topo da pilha |
| `?DUP` | `( n -- 0 \| n n )` | Duplica apenas se não for zero |
| `DROP` | `( n -- )` | Remove o topo da pilha |
| `SWAP` | `( a b -- b a )` | Troca os dois itens do topo |
| `OVER` | `( a b -- a b a )` | Copia o segundo item para o topo |
| `ROT` | `( a b c -- b c a )` | Rotaciona o terceiro item para o topo |
| `-ROT` | `( a b c -- c a b )` | Rotaciona o item do topo para o terceiro |
| `NIP` | `( a b -- b )` | Remove o segundo item |
| `TUCK` | `( a b -- b a b )` | Duplica o topo sob o segundo |
| `2DUP` | `( a b -- a b a b )` | Duplica o par do topo |
| `2DROP` | `( a b -- )` | Remove o par do topo |
| `2SWAP` | `( a b c d -- c d a b )` | Troca os dois pares do topo |
| `2OVER` | `( a b c d -- a b c d a b )` | Copia o terceiro par para o topo |
| `>R` | `( n -- )` ( R: `-- n` ) | Move da pilha de dados para a de retorno |
| `R>` | `( -- n )` ( R: `n --` ) | Move da pilha de retorno para a de dados |
| `R@` | `( -- n )` ( R: `n -- n` ) | Copia o topo da pilha de retorno para a de dados |
| `DEPTH` | `( -- +n )` | Número de itens na pilha de dados |

## Memória

| Palavra | Efeito na Pilha | Descrição |
|---|---|---|
| `@` | `( addr -- n )` | Busca célula de 16 bits da memória |
| `!` | `( n addr -- )` | Armazena célula de 16 bits na memória |
| `C@` | `( addr -- c )` | Busca byte da memória |
| `C!` | `( c addr -- )` | Armazena byte na memória |
| `+!` | `( n addr -- )` | Adiciona n à célula de memória |
| `?` | `( addr -- )` | Busca e imprime o valor no endereço |
| `FILL` | `( addr len c -- )` | Preenche região de memória com byte |
| `ERASE` | `( addr len -- )` | Preenche região de memória com zeros |
| `CMOVE` | `( dst src len -- )` | Copia bytes para frente |
| `CMOVE>` | `( dst src len -- )` | Copia bytes para trás |

## Aritmética

| Palavra | Efeito na Pilha | Descrição |
|---|---|---|
| `+` | `( a b -- soma )` | Adição |
| `-` | `( a b -- dif )` | Subtração (b - a) |
| `*` | `( a b -- prod )` | Multiplicação |
| `/` | `( a b -- quoc )` | Divisão (b / a) |
| `MOD` | `( a b -- resto )` | Resto (b mod a) |
| `/MOD` | `( a b -- resto quoc )` | Divide, retorna resto e quociente |
| `NEGATE` | `( n -- -n )` | Negação em complemento de dois |
| `ABS` | `( n -- \|n\| )` | Valor absoluto |
| `MIN` | `( a b -- min )` | Mínimo |
| `MAX` | `( a b -- max )` | Máximo |
| `1+` | `( n -- n+1 )` | Incrementa em 1 |
| `1-` | `( n -- n-1 )` | Decrementa em 1 |
| `2+` | `( n -- n+2 )` | Incrementa em 2 |
| `2-` | `( n -- n-2 )` | Decrementa em 2 |
| `2*` | `( n -- n*2 )` | Multiplica por 2 (shift à esquerda) |
| `2/` | `( n -- n/2 )` | Divide por 2 (shift à direita) |

## Comparação e Lógica

| Palavra | Efeito na Pilha | Descrição |
|---|---|---|
| `=` | `( a b -- flag )` | Igual |
| `<>` | `( a b -- flag )` | Diferente |
| `<` | `( a b -- flag )` | Menor que (com sinal) |
| `>` | `( a b -- flag )` | Maior que (com sinal) |
| `<=` | `( a b -- flag )` | Menor ou igual (com sinal) |
| `>=` | `( a b -- flag )` | Maior ou igual (com sinal) |
| `0=` | `( n -- flag )` | Teste de zero |
| `0<>` | `( n -- flag )` | Teste de não-zero |
| `0<` | `( n -- flag )` | Teste de negativo |
| `0>` | `( n -- flag )` | Teste de positivo |
| `U<` | `( a b -- flag )` | Menor que (sem sinal) |
| `U>` | `( a b -- flag )` | Maior que (sem sinal) |
| `U<=` | `( a b -- flag )` | Menor ou igual (sem sinal) |
| `U>=` | `( a b -- flag )` | Maior ou igual (sem sinal) |
| `AND` | `( a b -- c )` | E bit-a-bit |
| `OR` | `( a b -- c )` | OU bit-a-bit |
| `XOR` | `( a b -- c )` | OU exclusivo bit-a-bit |
| `INVERT` | `( n -- ~n )` | NÃO bit-a-bit |

Flags retornam -1 (verdadeiro) ou 0 (falso), seguindo a convenção Forth.

## E/S

| Palavra | Efeito na Pilha | Descrição |
|---|---|---|
| `KEY` | `( -- c )` | Lê um caractere da stdin |
| `EMIT` | `( c -- )` | Escreve um caractere na stdout |
| `CR` | `( -- )` | Escreve carriage return + newline |
| `SPACE` | `( -- )` | Escreve um espaço |
| `SPACES` | `( n -- )` | Escreve n espaços |
| `TYPE` | `( addr len -- )` | Escreve len bytes a partir de addr |
| `.` | `( n -- )` | Imprime inteiro com sinal na base |
| `U.` | `( u -- )` | Imprime inteiro sem sinal na base |
| `."` | `( -- )` (imediata) | Imprime string literal durante interpretação |
| `PAGE` | `( -- )` | Limpa a tela (emite form feed) |
| `BYE` | `( -- )` | Sai do pforth |

## Interpretador / Diversos

| Palavra | Efeito na Pilha | Descrição |
|---|---|---|
| `WORD` | `( delim -- addr )` | Analisa próxima palavra delimitada por char, retorna endereço de string contada |
| `FIND` | `( addr -- xt flag )` | Busca string contada no dicionário |
| `EXECUTE` | `( xt -- )` | Executa palavra pelo token de execução |
| `INTERPRET` | `( -- )` | Interpreta texto no TIB a partir da posição IN atual |
| `QUIT` | `( -- )` | Reinicia pilhas e entra no laço REPL |
| `ABORT` | `( -- )` | Reinicia pilhas, para compilação, descarta TIB |
| `EXPECT` | `( addr len -- )` | Lê linha da stdin para o buffer, armazena contagem no endereço 1000 |
| `.S` | `( -- )` | Imprime conteúdo da pilha de dados com marcadores de profundidade `< ... >` |
| `STATE` | `( -- flag )` | Estado atual do interpretador (-1 compilando, 0 interpretando) |
| `BASE` | `( -- addr )` | Endereço da variável base |
| `DP` | `( -- addr )` | Endereço da variável ponteiro do dicionário |
| `SPAN` | `( -- addr )` | Endereço da variável span (armazena contagem de char do EXPECT) |
| `HERE` | `( -- addr )` | Ponteiro atual do dicionário |
| `PAD` | `( -- addr )` | Endereço do buffer temporário (HERE + 68) |

## Palavras do Compilador

### Palavras Definidoras

| Palavra | Efeito na Pilha | Descrição |
|---|---|---|
| `:` | `( -- )` (imediata) | Inicia definição de colchete |
| `;` | `( -- )` (imediata) | Finaliza definição de colchete |
| `CONSTANT` | `( n -- )` | Cria constante nomeada com valor n |
| `VARIABLE` | `( -- )` | Cria variável nomeada (2 bytes reservados) |
| `CREATE` | `( -- )` | Cria cabeçalho de dicionário nomeado, sem tamanho alocado |
| `IMMEDIATE` | `( -- )` | Marca a palavra mais recente como imediata |

### Diretivas do Compilador

| Palavra | Efeito na Pilha | Descrição |
|---|---|---|
| `[` | `( -- )` (imediata) | Entra em estado de interpretação |
| `]` | `( -- )` (imediata) | Entra em estado de compilação |
| `LITERAL` | `( n -- )` (imediata) | Compila o topo da pilha como literal |
| `[COMPILE]` | `( -- )` (imediata) | Força compilação da próxima palavra imediata |
| `'` | `( -- xt )` | Tick — obtém XT da próxima palavra |
| `>BODY` | `( xt -- pfa )` | Converte XT para endereço do campo de parâmetro |
| `RECURSE` | `( -- )` (imediata) | Compila recursivamente a definição atual |
| `S"` | `( -- addr len )` (imediata) | String literal — empurra addr len em tempo de execução |

### Alocação de Memória

| Palavra | Efeito na Pilha | Descrição |
|---|---|---|
| `ALLOT` | `( n -- )` | Avança DP em n bytes |
| `,` | `( n -- )` | Compila uma célula em HERE, avança DP em 2 |
| `C,` | `( c -- )` | Compila um byte em HERE, avança DP em 1 |

## Ramificação

### Condicional

| Palavra | Efeito na Pilha | Descrição |
|---|---|---|
| `IF` | `( -- )` (imediata) | Inicia ramificação condicional |
| `THEN` | `( -- )` (imediata) | Finaliza ramificação condicional |
| `ELSE` | `( -- )` (imediata) | Ramificação alternativa |

```
IF <parte-verdadeira> THEN
IF <parte-verdadeira> ELSE <parte-falsa> THEN
```

### Laços Indefinidos

| Palavra | Efeito na Pilha | Descrição |
|---|---|---|
| `BEGIN` | `( -- )` (imediata) | Inicia laço |
| `UNTIL` | `( -- )` (imediata) | Laço até flag verdadeira |
| `AGAIN` | `( -- )` (imediata) | Laço incondicional |
| `WHILE` | `( -- )` (imediata) | Saída condicional do laço BEGIN |
| `REPEAT` | `( -- )` (imediata) | Finaliza laço BEGIN/WHILE/REPEAT |

```
BEGIN ... UNTIL
BEGIN ... AGAIN
BEGIN ... WHILE ... REPEAT
```

### Laços Definidos

| Palavra | Efeito na Pilha | Descrição |
|---|---|---|
| `DO` | `( limit inicio -- )` (imediata) | Inicia laço contado |
| `LOOP` | `( -- )` (imediata) | Finaliza laço contado, passo = 1 |
| `+LOOP` | `( passo -- )` (imediata) | Finaliza laço contado, passo = n |
| `I` | `( -- indice )` | Índice do laço atual |
| `J` | `( -- indice )` | Índice do laço externo |
| `LEAVE` | `( -- )` (imediata) | Sai do laço antecipadamente |
| `UNLOOP` | `( -- )` | Descarta parâmetros de laço da pilha de retorno |

```
10 0 DO I . LOOP         \ imprime 0 1 2 ... 9
10 0 DO I . 3 +LOOP      \ imprime 0 3 6 9
```

## Primitivas de Execução

Estas palavras não são tipicamente chamadas diretamente; são emitidas pelo compilador:

| Palavra | Descrição |
|---|---|
| `LIT` | Empurra próxima célula inline como literal |
| `EXIT` | Retorna da definição de colchete |
| `BRANCH` | Salto relativo incondicional |
| `0BRANCH` | Salto relativo condicional (zero) |
| `(DO)` | Inicializa parâmetros DO/LOOP na pilha de retorno |
| `(LOOP)` | Incrementa índice e faz laço se não terminou |
| `(+LOOP)` | Adiciona passo e faz laço se não terminou |
