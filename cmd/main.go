package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/matinhimself/trie/models"
	"github.com/matinhimself/trie/pkg/hashtable"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
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
	CyanBackground    = "\033[42m\033[30m%s\033[0m\n"
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;36m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}



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


func loadMassiveData(hm *hashtable.HashTable) {
	sl := make([]int, 1000)
	for i := 0; i < 20; i++ {
		cd := rand.Intn(9999)
		middle := fmt.Sprintf("%04d", cd)
		for j := 0; j < 200; j++ {
			stId := fmt.Sprintf("%03d", j)
			gpa := math.Mod(rand.Float64(), 10.0) + 10.0
			student := models.NewStudent(
				"student number "+strconv.Itoa((i+1)*(j+1)),
				models.StudentID("980"+middle+"0"+stId),
				gpa,
				"CE",
			)
			sl[hm.Set(student)] += 1
		}
		fmt.Println(sl)
	}
}

func main() {
	fmt.Printf("%s", ClearScreen)
	hm, _ := hashtable.NewHashTable(1000)
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
				} else {
					fmt.Printf("%s", ClearScreen)
					fmt.Println()
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
									hm.Delete(res.Value.GetKey())
									typed = typed[:len(typed)-1]
									break
								} else if _counter == 1 {
									editStudent(res.Value.(*models.Student), hm)
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
		case keyboard.KeyF1:
			{
				fmt.Print(ClearScreen)
				st := addStudent()
				_, found := hm.Get(string(st.StudentID))
				if !found {
					hm.Set(st)
				} else {
					WaitForKey(ErrC("Student ID: "+st.StudentID+" is taken."))
				}
				fmt.Printf("%s", ClearScreen)
				fmt.Println(typed)
			}
		case keyboard.KeyF2:
			{
				fmt.Print(ClearScreen)
				hm.PrintAll()
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
					var sGpa, name , dic, studentID string
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
		case keyboard.KeyF4: {
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

		for i, re := range searchRes[:int(math.Min(float64(len(searchRes)), InlineSearchCount))] {
			if i+1 == selection {
				fmt.Printf(CyanBackground, re)
			} else {
				fmt.Println(Yellow(re))
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
		student := value.(*models.Student)
		err := writer.Write([]string{string(student.StudentID), student.FullName,
			student.Discipline, fmt.Sprintf("%.2f", student.GPA)})
		if err != nil {
			WaitForKey(ErrC(fmt.Sprintf("Something went wrong in writing to csv. %v", err)))
		}
	}

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

func WaitForKey(message string){
	fmt.Println()
	fmt.Println(message)
	fmt.Println(Text("Press any key to continue..."))
	_, _, _ = keyboard.GetKey()
}

func addStudent() *models.Student {
	var name, sStId, stId, dec string
	var gpa float64
	reader := bufio.NewReader(os.Stdin)
	err := getInput("Full Name: ", &name, reader)
	name = strings.TrimSpace(name)
	if err != nil {
		fmt.Println(err)
	}
	err = getInput("Student ID: ", &sStId, reader)
	sStId = strings.TrimSpace(sStId)
	if err != nil {
		fmt.Println(err)
	}
	for _, c := range sStId {
		if unicode.IsDigit(c) {
			stId += string(c)
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
	dec = strings.TrimSpace(dec)
	st := models.NewStudent(name, models.StudentID(stId), gpa, dec)
	return st
}