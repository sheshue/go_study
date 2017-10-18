package main

import (
	"fmt"
)

func main() {
	shiyan := make(map[string]string)
	shiyan["golang"] = "docker"
	shiyan["python"] = "flask web framework"
	shiyan["linux"] = "sys administrator"
	fmt.Print("traverse all keys:")
	for key := range shiyan {
		fmt.Printf(" %s", key)
	}
	fmt.Println()
	delete(shiyan, "linux")
	shiyan["golang"] = "beego web framework"
	v, found := shiyan["linux"]
	fmt.Println("found key \" linux \" yes or false: %t, value of key \"linux \": \" %s \"", found, v)
	fmt.Println("traverse all keys/value after changed:")
	for k, v := range shiyan {
		fmt.Printf("\" %s \": \" %s \" \n", k, v)
	}
}
