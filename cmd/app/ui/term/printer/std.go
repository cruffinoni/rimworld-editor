package printer

var globalPrint = NewPrint()

func Print(p []byte) {
	globalPrint.WriteToStd(p)
}

func PrintSf(p string, a ...any) {
	globalPrint.WriteToStdf(p, a...)
}

func PrintError(err error) {
	if err == nil {
		globalPrint.WriteToError([]byte("<nil>"))
	} else {
		globalPrint.WriteToError([]byte(err.Error()))
	}
}

func PrintErrorS(err string) {
	globalPrint.WriteToError([]byte(err))
}

func PrintErrorSf(err string, args ...any) {
	globalPrint.WriteToErrf(err, args...)
}

func PrintS(s string) {
	globalPrint.WriteToStd([]byte(s))
}
