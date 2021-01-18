package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/matinhimself/trie/models"
	"github.com/matinhimself/trie/pkg/hashtable"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const (
	RedColor          = "\033[1;31m%s\033[0m\n"
	GreenColor        = "\033[32m%s\033[0m\n"
	CyanBackground    = "\033[47m\033[30m%s\033[0m\n"
	Purple            = "\033[35m%s\033[0m\n"
	ClearScreen       = "\033[H\033[2J"
	InlineSearchCount = 3
)

func getInput(format string, destination *string, reader *bufio.Reader, params ...interface{}) error {
	fmt.Printf(format, params...)
	text, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	// convert CRLF to LF
	*destination = strings.Replace(text, "\n", "", -1)
	return nil
}

func loadData(hm *hashtable.HashTable) {
	sl := make([]int, 1000)
	for i := 0; i < 5; i++ {
		stid := fmt.Sprintf("%03d", i)
		student := models.NewStudent(
			"Matin Habibi",
			models.StudentID("980122680"+stid),
			19.5,
			"CE",
		)
		sl[hm.Set(student)] += 1
	}
	for i := 0; i < 5; i++ {
		stid := fmt.Sprintf("%03d", i)
		student := models.NewStudent(
			"Matin Habibi",
			models.StudentID("970122680"+stid),
			19.5,
			"CE",
		)
		sl[hm.Set(student)] += 1
	}
}
func loadMassiveData(hm *hashtable.HashTable) {
	sl := make([]int, 1000)
	for i := 0; i < 20; i++ {
		cd := rand.Intn(9999)
		middle := fmt.Sprintf("%04d", cd)
		for j := 0; j < 200; j++ {
			stid := fmt.Sprintf("%03d", j)
			gpa := math.Mod(rand.Float64(), 10.0) + 10.0
			student := models.NewStudent(
				"student number "+strconv.Itoa((i+1)*(j+1)),
				models.StudentID("980"+middle+"0"+stid),
				gpa,
				"CE",
			)
			sl[hm.Set(student)] += 1
		}
	}
}

func main() {
	fmt.Printf("%s", ClearScreen)
	hm, _ := hashtable.NewHashTable(1000)
	loadData(hm)
	//loadMassiveData(hm)
	menu(hm)
}

func menu(hm *hashtable.HashTable) {

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press ESC to quit")
	var typed string
	var selection int
	var searchRes []string
	LOOP:
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		switch key {
		case keyboard.KeyBackspace, keyboard.KeyBackspace2:
			{
				if len(typed) >= 1 {
					typed = typed[:len(typed)-1]
					fmt.Printf("%s", ClearScreen)
					fmt.Println(typed)
				}
			}
		case keyboard.KeyEnter:
			{
				if selection > 0 {
					typed = searchRes[selection-1]
					selection = 0
				} else {
					res, found := hm.Get(typed)
					if found {
						_counter := 0
						for {
							fmt.Printf(ClearScreen)
							fmt.Println(*res)
							fmt.Println()
							menu := []string{"() Delete", "() Edit", "() Back"}
							for i, s := range menu {
								if i == _counter {
									fmt.Printf(CyanBackground, s)
								} else {
									if i == 2 {
										fmt.Printf(RedColor, s)
									} else {
										fmt.Println(s)
									}
								}
							}
							_, secKey, err := keyboard.GetKey()
							if err != nil {
								panic(err)
							}
							if secKey == keyboard.KeyArrowUp {
								if _counter > 0 {
									_counter--
								}
							}
							if secKey == keyboard.KeyArrowDown {
								if _counter < 2 {
									_counter++
								}
							}
							if secKey == keyboard.KeyEnter {
								if _counter == 0 {
									hm.Delete(res.Value.GetKey())
									typed = typed[:len(typed)-1]
									break
								} else if _counter == 1 {
									editStudent(res.Value.(*models.Student))
									break
								} else if _counter == 2 {
									break
								}
							}
							if secKey == keyboard.KeyEsc {
								break
							}

						}
					} else {
						continue
					}
				}

				fmt.Printf("%s", ClearScreen)
				fmt.Println(typed)
			}
		case keyboard.KeyArrowDown:
			{
				fmt.Print(ClearScreen)
				if selection < int(math.Min(InlineSearchCount, float64(len(searchRes)))) {
					selection++
				}
				fmt.Println(typed)
			}
		case keyboard.KeyArrowUp:
			{
				fmt.Print(ClearScreen)
				if selection > 0 {
					selection--
				}
				fmt.Println(typed)
			}
		case keyboard.KeyF1:{
				fmt.Print(ClearScreen)
				st := addStudent()
				_, found := hm.Get(string(st.StudentID))
				if !found {
					hm.Set(st)
				} else {
					fmt.Printf(RedColor, "Student ID: " + st.StudentID + " is taken.")
					fmt.Printf(GreenColor, "Press enter to continue.")
					_, _, _ = keyboard.GetKey()
				}
				fmt.Printf("%s", ClearScreen)
				fmt.Println(typed)
			}
		case keyboard.KeyF2:
			{
				fmt.Print(ClearScreen)
				export(hm)
				continue
			}
		case keyboard.KeyEsc:
			break LOOP
		default:
			if unicode.IsDigit(char) {
				fmt.Printf("%s", ClearScreen)

				typed += string(char)
				fmt.Println(typed)
			} else {
				continue
			}
		}

		if len(typed) == 0 {
			searchRes = hm.GetAllKeys()
		} else {
			searchRes = hm.GetKeysWithPrefix(typed)
		}


		if len(searchRes) > 0 && searchRes[0] == typed {
			fmt.Print(ClearScreen)
			fmt.Printf(Purple, typed)
			searchRes = searchRes[1:]
		}


		for i, re := range searchRes[:int(math.Min(float64(len(searchRes)), InlineSearchCount))] {
			if i+1 == selection {
				fmt.Printf(CyanBackground, re)
			} else {
				fmt.Printf(GreenColor, re)
			}

		}

	}
}

func export(hm *hashtable.HashTable) {
	file, err := os.Create("export.csv")
	if err != nil {
		println(err)
	}
	defer func() {
		_ = file.Close()
	}()

	data := hm.GetAllPairs()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, value := range data {
		if value == nil {
			continue
		}
		stundent := value.(*models.Student)
		err := writer.Write([]string{string(stundent.StudentID), stundent.FullName,
			stundent.Discipline, fmt.Sprintf("%.2f", stundent.GPA)})
		if err != nil {
			_ = fmt.Errorf("something went wrong in writing to csv: %v", err)
		}
	}
	hm.PrintAll()

}

func getKeyboardInput(Fixed, placeHolder string) string {
	typed := placeHolder
	for {
		fmt.Print(ClearScreen)
		fmt.Printf(Fixed)
		fmt.Printf(typed)
		char, key, _ := keyboard.GetKey()
		if unicode.IsDigit(char) || unicode.IsLetter(char) || unicode.IsPunct(char) {
			typed += string(char)
		} else
		if key == keyboard.KeySpace {
			typed += " "
		} else
		if key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2 {
			if len(typed) > 0 {
				typed = typed[:len(typed)-1]
			}
		} else
		if key == keyboard.KeyEnter {
			break
		}
	}

	return typed
}

func editStudent(st *models.Student) {

	fmt.Print(ClearScreen)
	var curString string
	var name, dec, sgpa string
	var gpa float64

	curString += fmt.Sprintf(strings.TrimSpace(GreenColor), fmt.Sprintf("%-12s", "Name:"))
	name = getKeyboardInput(curString, st.FullName)
	curString += name + "\n"

	curString += fmt.Sprintf(strings.TrimSpace(GreenColor), fmt.Sprintf("%-12s", "Discipline:"))
	dec = getKeyboardInput(curString, st.Discipline)
	curString += dec + "\n"

	curString += fmt.Sprintf(strings.TrimSpace(GreenColor), fmt.Sprintf("%-12s", "GPA:"))
	sgpa = getKeyboardInput(curString, fmt.Sprintf("%.2f", st.GPA))
	gpa, err := strconv.ParseFloat(sgpa, 64)
	if err != nil {
		// TODO: Handle invalid float Input
		fmt.Printf(RedColor, "Wrong float format")
	}
	curString += sgpa + "\n"

	st.UpdateStudent(name, st.StudentID, gpa, dec)

}

func addStudent() *models.Student {
	var name, sstid, stid, dec string
	var gpa float64
	reader := bufio.NewReader(os.Stdin)

	err := getInput("Full Name: ", &name, reader)
	if err != nil {
		fmt.Println(err)
	}
	err = getInput("Student ID: ", &sstid, reader)
	if err != nil {
		fmt.Println(err)
	}
	for _, c := range sstid {
		if unicode.IsDigit(c) {
			stid += string(c)
		}
	}

	fmt.Print("GPA: ")
	_, err = fmt.Fscan(reader, &gpa)
	_, err = fmt.Fscanln(reader)
	if err != nil {
		print(err)
	}
	if err != nil {
		fmt.Print(ClearScreen)
		fmt.Println("Not acceptable float.")
		return addStudent()
	}

	err = getInput("Discipline: ", &dec, reader)
	if err != nil {
		fmt.Println(err)
	}
	st := models.NewStudent(name, models.StudentID(stid), gpa, dec)
	return st
}
