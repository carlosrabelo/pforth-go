package forth

import (
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

type execFrame struct {
	Body []Cell
	IP   int
}

func (f *Forth) ExecuteWord(xt int) {
	if xt < 0 || xt >= len(f.Words) {
		forthError("INVALID XT: %d", xt)
	}
	w := f.Words[xt]
	switch w.Type {
	case WordPrimitive:
		w.Code(f)
	case WordColon:
		saved := execFrame{Body: f.Body, IP: f.IP}
		f.Body = w.Body
		f.IP = 0
		for f.IP < len(f.Body) {
			if atomic.LoadInt32(&f.Interrupted) != 0 {
				atomic.StoreInt32(&f.Interrupted, 0)
				f.Body = saved.Body
				f.IP = saved.IP
				forthError("INTERRUPTED")
			}
			nxt := int(f.Body[f.IP])
			f.IP++
			f.ExecuteWord(nxt)
		}
		f.Body = saved.Body
		f.IP = saved.IP
	case WordConstant:
		f.DSPush(w.Body[0])
	case WordVariable, WordCreate:
		f.DSPush(Cell(w.PFA))
	}
}

func (f *Forth) QUIT() {
	f.State = false
	f.DS = f.DS[:0]

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-sigChan:
				atomic.StoreInt32(&f.Interrupted, 1)
			case <-quit:
				return
			}
		}
	}()
	defer func() {
		signal.Stop(sigChan)
		close(quit)
	}()

	f.EmitStr("pForth Go ready\r\n")
	for f.Running {
		f.replStep()
	}
}

func (f *Forth) replStep() {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprint(r)
			if _, ok := r.(*ForthError); ok {
				f.EmitStr("ABORT: ")
			} else {
				f.EmitStr("PANIC: ")
			}
			f.EmitStr(msg)
			f.EmitCR()
			f.Abort()
		}
	}()
	if !f.State {
		f.EmitStr("OK\r\n")
	}
	for f.IN >= len(f.TIB) {
		line, err := f.ReadLine()
		if err != nil {
			f.Running = false
			return
		}
		f.TIB = []byte(line + "\r")
		f.IN = 0
	}
	f.interpretLoop()
}

func (f *Forth) Abort() {
	atomic.StoreInt32(&f.Interrupted, 0)
	f.State = false
	f.DS = f.DS[:0]
	f.Body = nil
	f.IP = 0
	f.IN = len(f.TIB)
}

func (f *Forth) InterpretLine(line string) {
	f.TIB = []byte(line + "\r")
	f.IN = 0
	f.interpretLoop()
}
