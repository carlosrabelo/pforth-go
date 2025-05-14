: FIB ( n -- fib[n] )
  0 1 ROT 0 DO
    TUCK +
  LOOP DROP
;

: TEST-FIB
  20 FIB . CR
;