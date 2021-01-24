package main

import (
	"encoding/csv"
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/matinhimself/trie/models"
	"github.com/matinhimself/trie/pkg/hashtable"
	"io"
	"log"
	"os"
	"strconv"
	"unicode"
)

const (
	ClearScreen       = "\033[H\033[2J"
	InlineSearchCount = 5
)

var (
	ErrC    = Red
	Text    = Teal
	Succeed = Green
)

var (
	CyanBackground = "\033[42m\033[30m%s\033[0m\n"
	Red            = color.Red.Sprint
	Green          = color.Green.Sprint
	Yellow         = color.Yellow.Sprint
	Magenta        = color.Magenta.Sprint
	Teal           = color.Cyan.Sprint
)

func main() {
	fmt.Print(ClearScreen)
	hm, _ := hashtable.NewHashTable(1000)
	menu(hm)
}

func help() {
	fmt.Println("Type student id to search")
	fmt.Println(Yellow("Press"), Teal("\n  ESC"), " to quit.")
	fmt.Println(Teal("  F1 "), " to add new student.")
	fmt.Println(Teal("  F2 "), " to show complete list of students.")
	fmt.Println(Teal("  F3 "), " to load from exported csv.")
	fmt.Println(Teal("  F4 "), " to export students into a csv file.")
	fmt.Println(Teal("  F5 "), " to show manual.")
}

func menu(hm *hashtable.HashTable) {

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		err := keyboard.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	help()

	var typed string
	var selection int
	var searchRes = hm.GetAllKeys()
	var startIndex int

LOOP:
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		switch key {
		case keyboard.KeyF5:
			{
				fmt.Print(ClearScreen)
				help()
				continue
			}
		case keyboard.KeyBackspace, keyboard.KeyBackspace2:
			{
				if len(typed) > 0 {
					typed = typed[:len(typed)-1]
					fmt.Print(ClearScreen)
					fmt.Println(typed)
				} else {
					fmt.Print(ClearScreen)
					fmt.Println()
				}
			}
		case keyboard.KeyEnter:
			{
				if selection > 0 {
					typed = searchRes[selection-1+startIndex]
					selection, startIndex = 0, 0
				} else {
					res, found := hm.Get(typed)
					if found {
						stu := res.Value.(*models.Student)
						studentProfile(stu, &typed, hm)
					} else {
						continue
					}
				}

				fmt.Print(ClearScreen)
				fmt.Println(typed)
			}
		case keyboard.KeyArrowDown:
			{
				fmt.Print(ClearScreen)
				if selection < min(InlineSearchCount, len(searchRes)) {
					selection++
				} else if len(searchRes) > InlineSearchCount+startIndex {
					startIndex++
				}
				fmt.Println(typed)
			}
		case keyboard.KeyArrowUp:
			{
				fmt.Print(ClearScreen)
				if startIndex > 0 && selection == 1 {
					startIndex--
				} else if selection > 0 {
					selection--
				}
				fmt.Println(typed)
			}
		case keyboard.KeyF1:
			{
				fmt.Print(ClearScreen)
				st := addStudent()
				_, found := hm.Get(string(st.StudentID))
				if !found {
					hm.Set(st)
				} else {
					WaitForKey(ErrC("Student ID: " + st.StudentID + " is taken."))
				}
				fmt.Printf("%s", ClearScreen)
				fmt.Println(typed)
			}
		case keyboard.KeyF2:
			{
				fmt.Print(ClearScreen)
				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"Student ID", "Name", "Field", "GPA"})
				t.SetAutoIndex(true)
				t.SetStyle(table.StyleLight)
				if len(typed) > 0 {
					for _, pair := range hm.GetPairsWithPrefix(typed) {
						st := pair.Value.(*models.Student)
						t.AppendRow(table.Row{
							pair.Key,
							st.FullName,
							st.Discipline,
							st.GPA,
						})
					}
				} else {
					for _, pair := range hm.GetAllPairs() {
						st := pair.Value.(*models.Student)
						t.AppendRow(table.Row{
							pair.Key,
							st.FullName,
							st.Discipline,
							st.GPA,
						})
					}
				}
				t.Render()
				continue
			}
		case keyboard.KeyF3:
			{
				f, err := os.Open("export.csv")
				if err != nil {
					WaitForKey(ErrC("No file named export.csv found in directory"))
					fmt.Printf("%s", ClearScreen)
					fmt.Println(typed)
					f.Close()
					continue LOOP
				}
				r := csv.NewReader(f)
				for {
					record, err := r.Read()
					if err == io.EOF {
						f.Close()
						break
					}
					if len(record) < 4 {
						f.Close()
						continue
					}
					var sGpa, name, dic, studentID string
					var gpa float64
					studentID = record[0]
					name = record[1]
					dic = record[2]
					sGpa = record[3]
					gpa, err = strconv.ParseFloat(sGpa, 64)
					if err != nil {
						f.Close()
						continue
					}
					st := models.NewStudent(name, models.StudentID(studentID), gpa, dic)
					hm.Set(st)
				}
				WaitForKey(ErrC("Students imported successfully."))
				fmt.Printf("%s", ClearScreen)
				fmt.Println(typed)
				f.Close()
				continue LOOP
			}
		case keyboard.KeyF4:
			{
				export(hm)
				WaitForKey(Succeed("Students exported successfully."))
				fmt.Print(ClearScreen)
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
			fmt.Println(Magenta(typed))
			searchRes = searchRes[1:]
		}

		for i, re := range searchRes[min(startIndex, len(searchRes)):min(len(searchRes), startIndex+InlineSearchCount)] {
			if i+1 == selection {
				fmt.Printf(CyanBackground, re)
			} else {
				fmt.Println(Yellow(re))
			}

		}

	}
}

func studentProfile(student *models.Student, typed *string, hm *hashtable.HashTable) {

	_counter := 0
	for {
		fmt.Printf(ClearScreen)
		fmt.Println(*student)
		fmt.Println()
		menu := []string{"() Delete", "() Edit", "() Back"}
		for i, s := range menu {
			if i == _counter {
				fmt.Printf(CyanBackground, s)
			} else {
				if i == 2 {
					fmt.Printf(Red(s))
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
				hm.Delete(student.GetKey())
				*typed = (*typed)[:len(*typed)-1]
				break
			} else if _counter == 1 {
				editStudent(student, hm)
				break
			} else if _counter == 2 {
				break
			}
		}
		if secKey == keyboard.KeyEsc {
			break
		}

	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func export(hm *hashtable.HashTable) {
	file, err := os.Create("export.csv")
	if err != nil {
		println(err)
	}
	defer func() {
		err = file.Close()
		log.Println(Red(err))
	}()

	data := hm.GetAllPairs()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, pair := range data {
		if pair.Value == nil {
			continue
		}
		student := pair.Value.(*models.Student)
		err := writer.Write([]string{string(student.StudentID), student.FullName,
			student.Discipline, fmt.Sprintf("%.2f", student.GPA)})
		if err != nil {
			WaitForKey(ErrC(fmt.Sprintf("Something went wrong in writing to csv. %v", err)))
		}
	}

}

func getKeyboardInput(Fixed, placeHolder string) string {
	typed := placeHolder
	fmt.Print(ClearScreen)
	fmt.Print(Fixed)
	fmt.Print(typed)

	for {
		char, key, _ := keyboard.GetKey()
		if unicode.IsDigit(char) || unicode.IsLetter(char) || unicode.IsPunct(char) {
			typed += string(char)
			fmt.Print(string(char))
		} else
		if key == keyboard.KeySpace {
			typed += " "
			fmt.Print(" ")
		} else
		if key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2 {
			if len(typed) > 0 {
				typed = typed[:len(typed)-1]
				fmt.Print(ClearScreen)
				fmt.Print(Fixed)
				fmt.Print(typed)
			}
		} else
		if key == keyboard.KeyEnter {
			break
		}
	}

	return typed
}

func editStudent(st *models.Student, hm *hashtable.HashTable) {

	fmt.Print(ClearScreen)
	var curString string
	var name, dec, sGpa, stId string
	var gpa float64

	curString += Text(fmt.Sprintf("%-12s", "Student ID:"))
	stId = getKeyboardInput(curString, string(st.StudentID))
	curString += stId + "\n"

	curString += Text(fmt.Sprintf("%-12s", "Name:"))
	name = getKeyboardInput(curString, st.FullName)
	curString += name + "\n"

	curString += Text(fmt.Sprintf("%-12s", "Discipline:"))
	dec = getKeyboardInput(curString, st.Discipline)
	curString += dec + "\n"

	curString += Text(fmt.Sprintf("%-12s", "GPA:"))
	sGpa = getKeyboardInput(curString, fmt.Sprintf("%.2f", st.GPA))
	gpa, err := strconv.ParseFloat(sGpa, 64)
	if err != nil {
		// TODO: Handle invalid float Input
		fmt.Print(Red("Wrong float format"))
	}
	curString += sGpa + "\n"

	if string(st.StudentID) == stId {
		st.UpdateStudent(name, st.StudentID, gpa, dec)
	} else {

		_, found := hm.Get(stId)
		if found {
			WaitForKey(ErrC("Student ID " + stId + " is taken."))

			editStudent(st, hm)
		} else {
			tempSt := models.NewStudent(name, models.StudentID(stId), gpa, dec)
			hm.Delete(string(st.StudentID))
			hm.Set(tempSt)
		}
	}

}

func WaitForKey(message string) {
	fmt.Println()
	fmt.Println(message)
	fmt.Println(Text("Press any key to continue..."))
	_, _, _ = keyboard.GetKey()
}

func addStudent() *models.Student {
	fmt.Print(ClearScreen)
	var curString string
	var name, dec, sGpa, stId string
	var gpa float64

	curString += Text(fmt.Sprintf("%-12s", "Student ID:"))
	stId = getKeyboardInput(curString, "")
	curString += stId + "\n"

	curString += Text(fmt.Sprintf("%-12s", "Name:"))
	name = getKeyboardInput(curString, "")
	curString += name + "\n"

	curString += Text(fmt.Sprintf("%-12s", "Discipline:"))
	dec = getKeyboardInput(curString, "")
	curString += dec + "\n"

	curString += Text(fmt.Sprintf("%-12s", "GPA:"))
	sGpa = getKeyboardInput(curString, "")
	gpa, err := strconv.ParseFloat(sGpa, 64)
	if err != nil {
		// TODO: Handle invalid float Input
		fmt.Print(Red("Wrong float format"))
	}
	curString += sGpa + "\n"

	st := models.NewStudent(name, models.StudentID(stId), gpa, dec)
	return st
}
