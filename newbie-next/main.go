package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"slices"
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
	for {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		_, err := os.Open("Notes.json")
		fmt.Println(strings.Repeat("-", 12))
		if os.IsNotExist(err) {
			_, err = os.Create("Notes.json")
			if err != nil {
				printerror(err.Error())
			}
		}
		notes := loadNotes("Notes.json")
		for i := 0; i < len(notes); i++ {
			fmt.Printf("[%d]%s\n", i, notes[i].Header)
			for _, line := range notes[i].Text {
				fmt.Println(line)
			}
			fmt.Println(strings.Repeat("-", 12))
		}
		fmt.Println(len(notes))

		reader := bufio.NewReader(os.Stdin)
		fmt.Println(strings.Repeat("-", 12) + "\nNotes\n" + strings.Repeat("-", 12))
		fmt.Println("[1]Create note\n[2]Delete note\n[0]Exit")

		var operation int
		_, err = fmt.Scan(&operation)
		if err != nil {
			printerror(err.Error())
			continue
		}
		switch operation {
		case 1:
			fmt.Println("Please, enter header of the note:")
			texth, err := reader.ReadString('\n')
			if err != nil {
				printerror(err.Error())
				continue
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
		case 2:
			if len(notes) == 0 {
				fmt.Println("Notes are empty!")
				continue
			}
			fmt.Println("Please enter the index of the note: ")
			var deletenote int
			_, err = fmt.Scan(&deletenote)
			if err != nil {
				printerror(err.Error())
				continue
			}
			if deletenote >= len(notes) || deletenote < 0 {
				printerror("Incorrect index")
				continue
			}
			notes = slices.Delete(notes, deletenote, deletenote+1)
			data, err := json.MarshalIndent(notes, "", " ")
			if err != nil {
				printerror(err.Error())
				continue
			}
			err = os.WriteFile("Notes.json", data, 0644)
			if err != nil {
				printerror(err.Error())
				continue
			}
			fmt.Println("Deleted!")
		case 0:
			fmt.Println("Quiting program...")
			return
		}
	}
}
