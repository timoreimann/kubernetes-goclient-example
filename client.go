/*
MIT License

Copyright (c) 2016 Timo Reimann

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/rest"
)

const (
	defaultServer string = "http://127.0.0.1:8001"
	serverEnvVar  string = "SERVER"
	tokenEnvVar   string = "TOKEN"
	caFileEnvVar  string = "CA_FILE"
)

var logger *log.Logger

func init() {
	// Ignore timestamps.
	logger = log.New(os.Stdout, "", 0)
}

func usage(msg string) {
	logger.Printf("%s\n\n", msg)
	fmt.Fprintf(os.Stderr, "usage: %s version | deploy", path.Base(os.Args[0]))
	os.Exit(1)
}

func main() {
	// Skip 0-th argument containing the binary's name.
	args := os.Args[1:]
	if len(args) < 1 {
		usage("insufficient number of parameters")
	}

	opName := args[0]
	var op operation

	switch opName {
	case "version":
		op = &versionOperation{}
	case "deploy":
		op = &deployOperation{
			image: "nginx:latest",
			name:  "nginx",
			port:  8080,
		}
	default:
		usage(fmt.Sprintf("unknown operation: %s", opName))
	}

	server, token, caData, err := parseConnectionParams()
	if err != nil {
		logger.Fatalf("failed to parse configuration parameters: %s", err)
	}

	config := &rest.Config{
		Host:            server,
		BearerToken:     token,
		TLSClientConfig: rest.TLSClientConfig{CAData: caData},
	}
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Fatalf("could not connect to Kubernetes API: %s", err)
	}

	op.Do(c)
}

func parseConnectionParams() (server, token string, caData []byte, err error) {
	server = os.Getenv(serverEnvVar)
	if len(server) == 0 {
		server = defaultServer
	}

	token = os.Getenv(tokenEnvVar)

	caFile := os.Getenv(caFileEnvVar)
	if len(caFile) > 0 {
		caData, err = ioutil.ReadFile(caFile)
	}

	return server, token, caData, err
}
