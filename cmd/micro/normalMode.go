package main

import (
	"strconv"

	"github.com/zyedidia/tcell"
)

const (
	OpNone   int = 0
	OpGlobal     = iota
)

type NormalMode struct {
	count string
	op    int
}

func Repeat(fn func(), count int) {
	for i := 0; i < count; i++ {
		fn()
	}
}

func (m *NormalMode) OnKey(ev *tcell.EventKey) {
	v := CurView()
	r := ev.Rune()
	if (r > '0' && r < '9') || (r == '0' && len(m.count) > 0) {
		m.count += string(r)
		return
	}

	// TODO: Handle error
	nrepeat, _ := strconv.Atoi(m.count)
	// Execute command at least once
	if nrepeat == 0 {
		nrepeat = 1
	}

	if ev.Modifiers() == tcell.ModCtrl {
		switch ev.Key() {
		case tcell.KeyCtrlD:
			v.Cursor.DownN(v.Height / 2)
		case tcell.KeyCtrlU:
			v.Cursor.UpN(v.Height / 2)
		case tcell.KeyCtrlQ:
			v.Quit(true)
		case tcell.KeyCtrlR:
			Repeat(func() {
				v.Redo(true)
			}, nrepeat)
		}
	}

	if ev.Key() == tcell.KeyEsc {
		v.Escape(true)
	}

	if ev.Key() == tcell.KeyRune {
		switch r {
		case 'A':
			Repeat(func() {
				v.EndOfLine(true)
			}, nrepeat)
			mode = &InsertMode{}
		case 'h':
			Repeat(func() {
				v.CursorLeft(true)
			}, nrepeat)
		case 'k':
			Repeat(func() {
				v.CursorUp(true)
			}, nrepeat)
		case 'l':
			Repeat(func() {
				v.CursorRight(true)
			}, nrepeat)
		case 'j':
			Repeat(func() {
				v.CursorDown(true)
			}, nrepeat)
		case 'i':
			mode = &InsertMode{}
		case 'w':
			Repeat(func() {
				v.VimWordRight(true)
				//v.CursorRight(false)
			}, nrepeat)
		case 'b':
			Repeat(func() {
				v.VimWordLeft(true)
			}, nrepeat)
		case '0':
			v.StartOfLine(true)
		case '$':
			v.EndOfLine(true)
		case ':':
			v.CommandMode(true)
		case '/':
			v.Find(true)
		case 'n':
			Repeat(func() {
				v.FindNext(true)
			}, nrepeat)
		case 'N':
			Repeat(func() {
				v.FindPrevious(true)
			}, nrepeat)
		case 'g':
			if m.op == OpNone {
				m.op = OpGlobal
			} else if m.op == OpGlobal {
				v.CursorStart(true)
				m.op = OpNone
			}
		case 'G':
			v.CursorEnd(true)
		case 'o':
			v.EndOfLine(false)
			v.InsertNewline(true)
		case 'O':
			v.CursorUp(false)
			v.EndOfLine(false)
			v.InsertNewline(true)
		case 'u':
			Repeat(func() {
				v.Undo(true)
			}, nrepeat)
		}

		m.count = ""
	}
}
