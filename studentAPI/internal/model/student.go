package model

type Student struct {
	Id    int
	Name  string
	Age   int
	Score int
}

type StudentArr []Student

func (s StudentArr) Len() int {
	return len(s)
}

func (s StudentArr) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func (s StudentArr) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
