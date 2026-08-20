package main

import (
	//"fmt"
	"strings"
	"syscall/js"
)

var morseMap=map[rune]string{
	'A': ".-", 'B': "-...", 'C': "-.-.", 'D': "-..",
	'E': ".", 'F': "..-.", 'G': "--.", 'H': "....",
	'I': "..", 'J': ".---", 'K': "-.-", 'L': ".-..",
	'M': "--", 'N': "-.", 'O': "---", 'P': ".--.",
	'Q': "--.-", 'R': ".-.", 'S': "...", 'T': "-",
	'U': "..-", 'V': "...-", 'W': ".--", 'X': "-..-",
	'Y': "-.--", 'Z': "--..",

	'0': "-----", '1': ".----", '2': "..---",
	'3': "...--", '4': "....-", '5': ".....",
	'6': "-....", '7': "--...", '8': "---..",
	'9': "----.",

	' ': "/",
}
var reverseMorseMap = map[string]string{}
var mode = "textToMorse"

func init() {
	for ch, code := range morseMap {
		reverseMorseMap[code] = string(ch)
	}
}
func morseToText(morse string) string {
	words := strings.Split(morse, " / ")
	var result []string
	for _, word := range words {
		letters := strings.Split(strings.TrimSpace(word), " ")
		wordStr := ""
		for _, letter := range letters {
			if ch, ok := reverseMorseMap[letter]; ok {
				wordStr += ch
			}
		}
		result = append(result, wordStr)
	}
	return strings.Join(result, " ")
}
func swapMode(this js.Value, args []js.Value) interface{} {
	document := js.Global().Get("document")
	inputEl := document.Call("getElementById", "input")
	outputEl := document.Call("getElementById", "output")
	inputLabel := document.Call("getElementById", "inputLabel")
	outputLabel := document.Call("getElementById", "outputLabel")

	inputVal := inputEl.Get("value").String()
	outputVal := outputEl.Get("value").String()
	inputEl.Set("value", outputVal)
	outputEl.Set("value", inputVal)

	if mode == "textToMorse" {
		mode = "morseToText"
		inputLabel.Set("textContent", "Morse Code")
		outputLabel.Set("textContent", "Text")
	} else {
		mode = "textToMorse"
		inputLabel.Set("textContent", "Text")
		outputLabel.Set("textContent", "Morse Code")
	}
	return nil
}

func textToMorse(text string) string{
	text=strings.ToUpper(text)
	var result []string
	for _, ch:=range text{
		if code, ok:=morseMap[ch]; ok{
			result=append(result, code)
		}
	}
	return strings.Join(result, "")
}

func translate(this js.Value, args []js.Value) interface{}{
	document:=js.Global().Get("document")
	inputEl:=document.Call("getElementById", "input")
	if inputEl.IsNull() || inputEl.IsUndefined(){
		return nil
	}

	input := inputEl.Get("value").String()

	outputEl:=document.Call("getElementById", "output")

	if outputEl.IsNull() || outputEl.IsUndefined(){
		return nil
	}

	if mode == "textToMorse" {
		outputEl.Set("value", textToMorse(input))
	} else {
		outputEl.Set("value", morseToText(input))
	}

	return nil
}
// func playMorse(this js.Value, args []js.Value) interface{} {
// 	document := js.Global().Get("document")
// 	outputEl := document.Call("getElementById", "output")
// 	if outputEl.IsNull() || outputEl.IsUndefined() {
// 		return nil
// 	}
// 	output := outputEl.Get("value").String()
// 	if output == "" {
// 		return nil
// 	}

// 	go func() {
// 		ac := js.Global().Get("AudioContext")
// 		if ac.IsUndefined() {
// 			ac = js.Global().Get("webkitAudioContext")
// 		}
// 		audioCtx := ac.New()

// 		for _, ch := range output {
// 			switch ch {
// 			case '.':
// 				beep(audioCtx, 0.1)
// 				sleep(200)
// 			case '-':
// 				beep(audioCtx, 0.3)
// 				sleep(400)
// 			case ' ':
// 				sleep(300)
// 			case '/':
// 				sleep(700)
// 			}
// 		}
// 	}()

// 	return nil
// }

func playMorse(this js.Value, args []js.Value) interface{} {
	document := js.Global().Get("document")
	outputEl := document.Call("getElementById", "output")
	if outputEl.IsNull() || outputEl.IsUndefined(){
		return nil
	}
	output:=outputEl.Get("value").String()
	if output==""{
		return nil
	}

	visual:=document.Call("getElementById", "visual")
	visual.Set("innerHTML", "")

	for _, ch:=range output{
		span := document.Call("createElement", "span")
		span.Get("style").Set("fontSize", "24px")
		span.Get("style").Set("color", "#555")

		switch ch {
		case '.':
			span.Set("textContent", "● ")
		case '-':
			span.Set("textContent", "━ ")
		case ' ':
			span.Set("textContent", "  ")
		case '/':
			span.Set("textContent", " / ")
		}
		visual.Call("appendChild", span)
	}
	

	go func() {
		ac := js.Global().Get("AudioContext")
		if ac.IsUndefined(){
			ac=js.Global().Get("webkitAudioContext")
		}
		audioCtx:=ac.New()
		children:=visual.Get("childNodes")

		for i, ch := range output {
			elem:=children.Index(i)
			switch ch {
			case '.':
				elem.Get("style").Set("color", "#6366f1")
				beep(audioCtx, 0.1)
				sleep(200)
				elem.Get("style").Set("color", "#888")
			case '-':
				elem.Get("style").Set("color", "#6366f1")
				beep(audioCtx, 0.3)
				sleep(400)
				elem.Get("style").Set("color", "#888")
			case ' ':
				sleep(300)
			case '/':
				sleep(700)
			}
			i++
		}
	}()

	return nil
}

func beep(audioCtx js.Value, duration float64) {
	oscillator := audioCtx.Call("createOscillator")
	gainNode := audioCtx.Call("createGain")

	oscillator.Call("connect", gainNode)
	gainNode.Call("connect", audioCtx.Get("destination"))

	oscillator.Get("frequency").Set("value", 600)
	gainNode.Get("gain").Set("value", 0.3)

	now := audioCtx.Get("currentTime").Float()
	oscillator.Call("start", now)
	oscillator.Call("stop", now+duration)
}

func sleep(ms int) {
	ch := make(chan struct{})
	js.Global().Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		close(ch)
		return nil
	}), ms)
	<-ch
}

func main(){
	js.Global().Set("translate", js.FuncOf(translate))
	js.Global().Set("playMorse", js.FuncOf(playMorse))
	js.Global().Set("swapMode", js.FuncOf(swapMode))
	select{}
}