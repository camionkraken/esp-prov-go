package main

import (
	"encoding/json"
	"esp-prov-go/core/security"
	"esp-prov-go/softap"
	"fmt"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	provisioner, err := softap.NewSoftapProvisioner("", &security.Security0{})

	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := provisioner.GetProtoVersion()

	if err != nil {
		fmt.Println(err)
		return
	}

	j, err := json.Marshal(resp)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(j))

	err = provisioner.EstablishSession()

	if err == nil {
		fmt.Println("Session established successfully")
	} else {
		fmt.Println(err)
		return
	}
}
