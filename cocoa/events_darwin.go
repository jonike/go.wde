/*
   Copyright 2012 the go.wde authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package cocoa

// #include "gomacdraw/gmd.h"
import "C"

import (
	"fmt"
	"github.com/skelterjohn/go.wde"
	// "strings"
)

func getButton(b int) (which wde.Button) {
	switch b {
	case 0:
		which = wde.LeftButton
	}
	return
}

func containsGlyph(haystack []string, needle string) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}
	return false
}

func (w *Window) EventChan() (events <-chan interface{}) {
	downKeys := make(map[string]bool)
	ec := make(chan interface{})
	go func(ec chan<- interface{}) {

		var noX int = 1<<31 - 1
		noX++
		var lastX, lastY int = noX, 0

	eventloop:
		for {
			e := C.getNextEvent(w.cw)
			switch e.kind {
			case C.GMDNoop:
				continue
			case C.GMDMouseDown:
				var mde wde.MouseDownEvent
				mde.Where.X = int(e.data[0])
				mde.Where.Y = int(e.data[1])
				mde.Which = getButton(int(e.data[2]))
				lastX = mde.Where.X
				lastY = mde.Where.Y
				ec <- mde
			case C.GMDMouseUp:
				var mue wde.MouseUpEvent
				mue.Where.X = int(e.data[0])
				mue.Where.Y = int(e.data[1])
				mue.Which = getButton(int(e.data[2]))
				lastX = mue.Where.X
				lastY = mue.Where.Y
				ec <- mue
			case C.GMDMouseDragged:
				var mde wde.MouseDraggedEvent
				mde.Where.X = int(e.data[0])
				mde.Where.Y = int(e.data[1])
				mde.Which = getButton(int(e.data[2]))
				if lastX != noX {
					mde.From.X = int(lastX)
					mde.From.Y = int(lastY)
				} else {
					mde.From.X = mde.Where.X
					mde.From.Y = mde.Where.Y
				}
				lastX = mde.Where.X
				lastY = mde.Where.Y
				ec <- mde
			case C.GMDMouseMoved:
				var mme wde.MouseMovedEvent
				mme.Where.X = int(e.data[0])
				mme.Where.Y = int(e.data[1])
				if lastX != noX {
					mme.From.X = int(lastX)
					mme.From.Y = int(lastY)
				} else {
					mme.From.X = mme.Where.X
					mme.From.Y = mme.Where.Y
				}
				lastX = mme.Where.X
				lastY = mme.Where.Y
				ec <- mme
			case C.GMDMouseEntered:
				var me wde.MouseEnteredEvent
				me.Where.X = int(e.data[0])
				me.Where.Y = int(e.data[1])
				if lastX != noX {
					me.From.X = int(lastX)
					me.From.Y = int(lastY)
				} else {
					me.From.X = me.Where.X
					me.From.Y = me.Where.Y
				}
				lastX = me.Where.X
				lastY = me.Where.Y
				ec <- me
			case C.GMDMouseExited:
				var me wde.MouseExitedEvent
				me.Where.X = int(e.data[0])
				me.Where.Y = int(e.data[1])
				if lastX != noX {
					me.From.X = int(lastX)
					me.From.Y = int(lastY)
				} else {
					me.From.X = me.Where.X
					me.From.Y = me.Where.Y
				}
				lastX = me.Where.X
				lastY = me.Where.Y
				ec <- me
			case C.GMDKeyDown:
				var letter string
				var ke wde.KeyEvent
				keycode := int(e.data[1])

				blankLetter := containsInt(blankLetterCodes, keycode)
				if !blankLetter {
					letter = fmt.Sprintf("%c", e.data[0])
				}

				ke.Key = keyMapping[keycode]

				if !downKeys[ke.Key] {
					ec <- wde.KeyDownEvent(ke)
				}

				downKeys[ke.Key] = true

				ec <- wde.KeyTypedEvent{
					KeyEvent: ke,
					Chord:    wde.ConstructChord(downKeys),
					Glyph:    letter,
				}

			case C.GMDKeyUp:
				var ke wde.KeyUpEvent
				ke.Key = keyMapping[int(e.data[1])]
				delete(downKeys, ke.Key)
				ec <- ke
			case C.GMDResize:
				var re wde.ResizeEvent
				re.Width = int(e.data[0])
				re.Height = int(e.data[1])
				ec <- re
			case C.GMDClose:
				ec <- wde.CloseEvent{}
				break eventloop
				return
			}
		}
		close(ec)
	}(ec)
	events = ec
	return
}
