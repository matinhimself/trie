package models

import (
	"fmt"
	"strconv"
)

type StudentID string

type Student struct {
	FullName   string
	StudentID  StudentID
	GPA        float64
	Discipline string
}



func (s Student) String() string {
	return fmt.Sprintf("%s,%s,%.2f,%s", s.StudentID, s.FullName, s.GPA, s.Discipline)
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

func getBinaryString(s string) string {
	return ""
}



func (s *Student) ToHash() uint32 {
	var h uint32
	stid := s.StudentID
	h1, _ := strconv.Atoi(string(stid))
	h = uint32(h1)
	h *= 2654435761
	return h
}

func (s *Student) GetKey() string {
	return string(s.StudentID)
}

//
//func (s *Student) ToHash() uint64 {
//	var h float64
//	for i := 0; i < len(s.StudentID); i++ {
//		h += math.Pow(97,float64(i)) * float64(s.StudentID[i])
//	}
//	h = math.Mod(h, float64(4999))
//
//	return uint64(h)
//}
