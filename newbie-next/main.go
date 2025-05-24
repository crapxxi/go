package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func printerror(msg string) {
	fmt.Println("Error: " + msg)
	time.Sleep(time.Second * 2)
}

type Notes struct {
	Header string   `json:"header"`
	Text   []string `json:"list"`
}

func saveNotes(newNote Notes, filename string) {
	notes := loadNotes(filename)
	notes = append(notes, newNote)

	data, err := json.MarshalIndent(notes, "", "  ")
	if err != nil {
		printerror("Error marshaling JSON: " + err.Error())
		return
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		printerror("Error writing file: " + err.Error())
		return
	}
	fmt.Println("Notes saved!")
}

func loadNotes(filename string) []Notes {
	var notes []Notes

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []Notes{}
		}
		printerror("Error reading file: " + err.Error())
		return []Notes{}
	}

	if len(data) == 0 {
		return []Notes{}
	}

	err = json.Unmarshal(data, &notes)
	if err != nil {
		printerror("Error unmarshaling JSON: " + err.Error())
		return []Notes{}
	}

	return notes
}

func main() {
	file, err := os.Open("Notes.json")
	fmt.Println(strings.Repeat("-", 12))
	if os.IsNotExist(err) {
		file, err = os.Create("Notes.json")
		if err != nil {
			printerror(err.Error())
		}
	}
	defer file.Close()
	notes := loadNotes("Notes.json")
	for _, note := range notes {
		fmt.Println(note.Header)
		for _, line := range note.Text {
			fmt.Println(line)
		}
		fmt.Println(strings.Repeat("-", 12))
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println(strings.Repeat("-", 12) + "\nNotes\n" + strings.Repeat("-", 12))
	fmt.Println("[1]Create note\n[2]Edit note\n[3]Delete note\n[0]Exit")

	var operation int
	_, err = fmt.Scan(&operation)
	if err != nil {
		printerror(err.Error())
		return
	}
	switch operation {
	case 1:
		fmt.Println("Please, enter header of the note:")
		texth, err := reader.ReadString('\n')
		if err != nil {
			printerror(err.Error())
			return
		} else {
			texth = strings.TrimSpace(texth)
			fmt.Println("Please, enter the main text(Press CTRL+C when you're done):")
			lines := []string{}
			input := make(chan string)

			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

			go func() {
				for {
					textm, errm := reader.ReadString('\n')
					if errm != nil {
						break
					}
					textm = strings.TrimSpace(textm)
					input <- textm
				}
			}()
		loop:
			for {
				select {
				case line := <-input:
					lines = append(lines, line)
				case <-sigs:
					break loop
				}
			}

			thisnote := Notes{
				Header: texth,
				Text:   lines,
			}
			saveNotes(thisnote, "Notes.json")
		}
	}
}
