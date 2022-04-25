package main

import (
	"flag"
	"fmt"
	"github.com/alt-golang/config-server/web/bindings/gin/context"
	"os"
)

func main() {

	versionPtr := flag.Bool("v", false, "output the version number")
	configPtr := flag.String("c", "config"+fmt.Sprint(os.PathSeparator)+"internal", "the internal config directory for the server itself")
	helpPtr := flag.Bool("h", false, "output usage")
	flag.Parse()

	if *helpPtr == false {
		if *versionPtr {
			fmt.Printf("github.com/alt-golang/config-server %s\n", Version)
		} else if *configPtr == "" {
			fmt.Println("Error: the internal config directory (flag -c) is required")
			flag.Usage()
		} else {
			context.Start()
		}
	} else {
		flag.Usage()
	}

}
