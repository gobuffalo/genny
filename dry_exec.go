package genny

// DryExec will just log out the command and it's arguments
// it will NOT run the command
func DryExec(g Generator) Generator {
	g = WithExec(g, nil)
	return g
}
