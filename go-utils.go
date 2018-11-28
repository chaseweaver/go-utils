package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 || (len(os.Args) > 1 && os.Args[1] == "-h" || os.Args[1] == "-help") {
		fmt.Print("\ngo-utils v0.1-alpha-20181128-0116\n")
		fmt.Print("- Chase Weaver <github.com/chaseweaver/go-utils>\n\n")
		fmt.Print("Common go-utils commands for use in various situations:\n\n")
		for _, value := range Commands {
			fmt.Println(fmt.Sprintf("%-16s %10s", value.Name, value.Description))
		}
		fmt.Println()
		os.Exit(0)
	}

	fmt.Println()
	if cmd, ok := Commands[os.Args[1]]; ok {
		if cmd.FlagSet != nil {
			cmd.FlagSet.Parse(os.Args[2:])
		}
		cmd.Action()
	} else {
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
	fmt.Println()
}
