package constants

type Status struct {
	Progress int
	Name string
}

var TODO = Status{1, "TODO"}
var GOING = Status{2, "GOING"}
var DONE = Status{3, "DONE"}

var Statuses []Status
func init()  {
   Statuses = []Status{TODO, GOING, DONE}
}