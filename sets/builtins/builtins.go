package builtins

//go:generate genny -in ../generic/generic.go -out builtins_generated.go -pkg builtins gen "Item=BUILTINS"
//go:generate addgeneratedheader builtins_generated.go "by $GOPACKAGE/$GOFILE"
