package unit

import "fmt"

var logo = `
'   ______     __    __      __      _    _____      __      _      _____   
'  (   __ \    ) )  ( (     /  \    / )  (_   _)    /  \    / )    / ___ \  
'   ) (__) )  ( (    ) )   / /\ \  / /     | |     / /\ \  / /    / /   \_) 
'  (    __/    ) )  ( (    ) ) ) ) ) )     | |     ) ) ) ) ) )   ( (  ____  
'   ) \ \  _  ( (    ) )  ( ( ( ( ( (      | |    ( ( ( ( ( (    ( ( (__  ) 
'  ( ( \ \_))  ) \__/ (   / /  \ \/ /     _| |__  / /  \ \/ /     \ \__/ /  
'   )_) \__/   \______/  (_/    \__/     /_____( (_/    \__/       \____/   
'                                                                           
`
var line_bottom = `-------------------------------RUNING-----------------------------------`
var line_top = `------------------------------------------------------------------------`

func PrintLogo() {
	fmt.Println(line_top)
	fmt.Println(logo)
	fmt.Println(line_bottom)
}

var GlabalObject glabalObject

type glabalObject struct {
	Name string
	IP   string
	Port int
}

func init() {
	GlabalObject = glabalObject{
		Name: "SimpleServer",
		IP:   "127.0.0.1",
		Port: 8888,
	}
}
