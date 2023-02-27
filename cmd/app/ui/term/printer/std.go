package printer

var globalPrint = NewPrint()

func Print(p []byte) {
	globalPrint.WriteToStd(p)
}

func PrintSf(p string, a ...interface{}) {
	globalPrint.WriteToStdf(p, a...)
}

func PrintError(err error) {
	globalPrint.WriteToError([]byte(err.Error()))
}

func PrintErrorS(err string) {
	globalPrint.WriteToError([]byte(err))
}

func PrintS(s string) {
	globalPrint.WriteToStd([]byte(s))
}
