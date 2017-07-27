package main

import (
	"strings"

	"github.com/zyedidia/tcell"
)

type InsertMode struct {
	count string
}

func (m *InsertMode) OnKey(e *tcell.EventKey) {
	v := CurView()
	if e.Key() == tcell.KeyRune {
		// Check viewtype if readonly don't insert a rune (readonly help and log view etc.)
		if v.Type.readonly == false {
			for _, c := range v.Buf.cursors {
				v.SetCursor(c)

				// Insert a character
				if v.Cursor.HasSelection() {
					v.Cursor.DeleteSelection()
					v.Cursor.ResetSelection()
				}
				v.Buf.Insert(v.Cursor.Loc, string(e.Rune()))

				for pl := range loadedPlugins {
					_, err := Call(pl+".onRune", string(e.Rune()), v)
					if err != nil && !strings.HasPrefix(err.Error(), "function does not exist") {
						TermMessage(err)
					}
				}

				if recordingMacro {
					curMacro = append(curMacro, e.Rune())
				}
			}
			v.SetCursor(&v.Buf.Cursor)
		}
	} else if e.Key() == tcell.KeyEsc {
		mode = &NormalMode{}
	}
}
