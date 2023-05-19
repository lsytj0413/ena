// Copyright (c) 2023 The Songlin Yang Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package main provide an example for use conma to retrieve configuration items
package main

import (
	"fmt"
	"reflect"

	"github.com/lsytj0413/ena/conma"
)

type Config struct {
	v1 string `conma:"k1"`
	v2 string `conma:"k2"`
	v3 string `conma:"k3"`
	v4 string `conma:"k4:=v4"`

	c Config1 `conma:"p"`
}

type Config1 struct {
	v1 string `conma:"k1"`
	v2 string `conma:"k2"`
	v3 string `conma:"k3"`
	v4 string `conma:"k4:=v4"`
}

// K1=ev1 K2=ev2 P_K1=pev1 P_K2=pev2 go run main.go -- --k1=ov1 --p.k1=pov1 -p.k4=pov4
func main() {
	// configuration priority is: option > env > file
	conma.DefaultConfigMgr().AddConfigReader(conma.NewFileConfigReader("./config.yaml"))
	conma.DefaultConfigMgr().AddConfigReader(conma.NewEnvConfigReader())
	conma.DefaultConfigMgr().AddConfigReader(conma.NewOptionConfigReader())
	err := conma.DefaultConfigMgr().ReadConfig()
	if err != nil {
		panic(err)
	}

	var c Config
	err = conma.DefaultConfigMgr().Unmarshal(&c)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", c)
	if !reflect.DeepEqual(c, Config{
		v1: "ov1",
		v2: "ev2",
		v3: "fv3",
		v4: "v4",
		c: Config1{
			v1: "pov1",
			v2: "pev2",
			v3: "pfv3",
			v4: "pov4",
		},
	}) {
		panic("unexpect config value")
	}
}
