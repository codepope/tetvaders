package main

import "github.com/codepope/tetvaders/pkg/tetvaders"

func main() {
	t := tetvaders.Tetvaders{}
	t.Init()
	t.Run()
}
