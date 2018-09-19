package printer

type Printer interface {
	Println(a ...interface{}) (int, error)
	Printf(format string, a ...interface{}) (int, error)
	Print(fa ...interface{}) (int, error)
	NextLevel() Printer
}
