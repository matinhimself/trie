package models

import (
	"fmt"
	"github.com/matinhimself/trie/pkg/hashtable"
	"hash/maphash"
	"strconv"
	"unsafe"
)

type StudentID string

type Student struct {
	FullName   string
	StudentID  StudentID
	GPA        float64
	Discipline string
}

func (s Student) String() string {
	return fmt.Sprintf("%-15s %s\n%-15s %s\n%-15s %.2f\n%-15s %s", "Full Name:", s.FullName,
		"Student ID:", s.StudentID, "GPA:", s.GPA, "Discipline:", s.Discipline)
}

func NewStudent(fullName string, studentID StudentID, GPA float64, discipline string) *Student {
	return &Student{FullName: fullName, StudentID: studentID, GPA: GPA, Discipline: discipline}
}
func (s *Student) UpdateStudent(fullName string, studentID StudentID, GPA float64, discipline string) {
	s.FullName = fullName
	s.StudentID = studentID
	s.GPA = GPA
	s.Discipline = discipline
}



//using Golang's collision-resistant hash algorithm. Twice!.
func (s *Student) ToHash() uint64 {
	var h maphash.Hash
	_, _ = h.WriteString(string(s.StudentID))
	_, _ = h.WriteString(strconv.Itoa(int(h.Sum64())))
	var sum = h.Sum64()
	var res uint64
	for _, ch := range (*[4]byte)(unsafe.Pointer(&sum))[:] {
		res = uint64(ch) + (res << 5) + (res >> 7) - res
	}
	return res
}

// using Golang's collision-resistant hash algorithm. Twice!.
//func (s *Student) ToHash() uint64 {
//	var h maphash.Hash
//	_, _ = h.WriteString(string(s.StudentID))
//	_, _ = h.WriteString(strconv.Itoa(int(h.Sum64())))
//	return h.Sum64()
//}


//func (s *Student) ToHash() uint64 {
//	var h uint64
//	for _, ch := range s.GetKey() {
//		h = uint64(ch) + (h << 5) + (h >> 7) - h
//	}
//	return h
//}


// Implements the sha256 hash function
//func (s *Student) ToHash() uint32 {
//	var h uint32
//	sha := sha256.New()
//	sha.Write([]byte(s.StudentID))
//	bh := sha.Sum(nil)
//	return h
//}


//func (s *Student) ToHash() uint32 {
//	var h uint32
//	stId := string(s.StudentID)
//	for i := 0; i < len(stId); i++ {
//		h += uint32(stId[i])
//		h += h << 7
//		h ^= h >> 5
//	}
//	return h
//}

// Implements the Jenkins hash function
//func (s *Student) ToHash() uint32 {
//	var h uint32
//	for _, c := range s.GetKey(){
//		h += uint32(c)
//		h += h << 3
//		h ^= h >> 5
//	}
//	h += h << 3
//	h ^= h >> 11
//	h += h << 15
//
//	return h
//}


// Implements mine algorithm
//func (s *Student) ToHash() uint32 {
//	var h uint32
//	for _, ch := range s.GetKey() {
//		h = uint32(ch) + h << 3 + h << 5 - h
//	}
//	return h
//}

// Safe multiplication for indexing
// in situation of overflow it will always
// return positive number and a false and
// false boolean.
//func Multi64(u1, u2 uint64) (res uint64, ok bool) {
//	if u2 >= (math.MaxUint64 / u1) {
//		u3 := u1 * u2
//		if u3 < 0 {
//			u3 = -u3
//		}
//		return u3, false
//	} else {
//		return u1 * u2, true
//	}
//}

//func reverse(s string) (result string) {
//	for _,v := range s {
//		result = string(v) + result
//	}
//	return
//}

func (s *Student) Equals(other *hashtable.HashAble) bool {
	otherSt, ok := (*other).(*Student)
	return ok && otherSt.StudentID == s.StudentID
}

func (s *Student) GetKey() string {
	return string(s.StudentID)
}
