package models

import (
	"crypto/sha256"
	"fmt"
	"github.com/matinhimself/trie/pkg/hashtable"
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

func (s *Student) ToHash() uint32 {
	var h uint32
	sha := sha256.New()
	sha.Write([]byte(s.StudentID))
	bh := sha.Sum(nil)
	for i := 0; i < len(bh); i++ {
		h += uint32(bh[i])
		h += h << 3
		h ^= h >> 5
	}
	return h
}

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

func (s *Student) Equals(other *hashtable.HashAble) bool {
	otherSt, ok := (*other).(*Student)
	return ok && otherSt.StudentID == s.StudentID
}

//Implements the Jenkins hash function
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


func (s *Student) GetKey() string {
	return string(s.StudentID)
}