package test

import (
	"fmt"
	"github.com/matinhimself/trie/models"
	"github.com/matinhimself/trie/pkg/hashtable"
	"hash/maphash"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"testing"
)

func TestHashTableDelete(t *testing.T) {
	stId := "980122680000"
	hm, _ := hashtable.NewHashTable(200)
	st := models.NewStudent("Test test", models.StudentID(stId), 16.5, "TE")
	hm.Set(st)
	hm.Delete(stId)
	if _, found := hm.Get(stId); found {
		t.Error("Student didn't delete.'")
	}
}

func loadMassiveData(studentCount int, middleCount int, hm *hashtable.HashTable, ls []int) {
	for i := 0; i < middleCount; i++ {
		middle := fmt.Sprintf("%04d", rand.Intn(9999))
		for j := 0; j < studentCount; j++ {
			stId := fmt.Sprintf("%03d", j)
			gpa := math.Mod(rand.Float64(), 10.0) + 10.0
			student := models.NewStudent(
				"student number "+strconv.Itoa((i+1)*(j+1)),
				models.StudentID("91"+middle+"0"+stId),
				gpa,
				"CE",
			)
			gpa = math.Mod(rand.Float64(), 10.0) + 10.0
			student2 := models.NewStudent(
				"student number "+strconv.Itoa((i+1)*(j+1))+"v2",
				models.StudentID("91"+middle+"1"+stId),
				gpa,
				"CE",
			)
			ls[hm.Set(student)] += 1
			ls[hm.Set(student2)] += 1
		}
	}
}

func TestHashTableCollision(t *testing.T) {
	ls := make([]int, 1000)
	hm, _ := hashtable.NewHashTable(1000)
	loadMassiveData(75, 200, hm, ls)
	min, max := MinMax(ls)
	t.Log(Magenta("\nCollision Test Result:"), Teal("\nElements Count: "),
		200*75*4, Teal("\nHashmap Size: "), hm.Size(), Teal("\nRange: "), min, "-", max)
	sort.Sort(sort.Reverse(sort.IntSlice(ls)))
}

func i64tob(val uint64) []byte {
	r := make([]byte, 8)
	for i := uint64(0); i < 8; i++ {
		r[i] = byte((val >> (i * 8)) & 0xff)
	}
	return r
}

type testHashAble struct {
	val uint64
}
func (t *testHashAble) GetKey() string {return strconv.FormatUint(t.val, 10) }
func (t *testHashAble) ToHash() uint64 {
	var h maphash.Hash
	_, _ = h.Write(i64tob(t.val))
	return h.Sum64()
}

func (t *testHashAble) Equals(other *hashtable.HashAble) bool {
	otherSt, ok := (*other).(*testHashAble)
	return ok && otherSt.val == t.val
}

func BenchmarkHashAdd(b *testing.B){
	hm, _ := hashtable.NewHashTable(1000)
	for i := 0; i < b.N; i++ {
		hm.Set(&testHashAble{val: uint64(b.N)})
	}
}

func TestHashTableAdd(t *testing.T) {
	hm, _ := hashtable.NewHashTable(200)
	loadMassiveData(75, 50, hm, make([]int, 1000))
}

var (
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

func MinMax(array []int) (int, int) {
	var max = array[0]
	var min = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}
