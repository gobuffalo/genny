package genny

// WetRun will wrap the Generator with `WetExec`
// as well as wrap it in `WetFilesHandler` to write files to disk
func WetRun(g Generator) Generator {
	g = WetExec(g)
	return WetFilesHandler(g)
}
