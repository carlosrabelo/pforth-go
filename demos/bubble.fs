CREATE DATA 9 C, 8 C, 7 C, 6 C, 5 C, 4 C, 3 C, 2 C, 1 C, 0 C,

: PRINT
  10 0 DO DATA I + C@ . LOOP CR
;

: BUBBLE ( addr len -- )
  DUP 1 DO
    DUP I - 0 DO
      OVER J + C@ OVER J 1+ + C@
      OVER OVER > IF
        SWAP OVER J 1+ + C! J + C!
      ELSE
        2DROP
      THEN
    LOOP
  LOOP
  DROP
;

: TEST-BUBBLE
  ." Before: " PRINT
  DATA 10 BUBBLE
  ." After:  " PRINT
;