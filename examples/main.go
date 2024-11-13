package main

import (
	"log"

	gitexec "github.com/instruqt/git-exec"
)

func main() {
	git, err := gitexec.New()
	if err != nil {
		log.Fatal(err)
	}

	// tmp, err := os.MkdirTemp("", "git-exec")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer os.RemoveAll(tmp)

	git.SetWorkingDirectory("/Users/erik/code/instruqt/git-exec")

	output, err := git.Show("HEAD")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("show", output)

	// output, err := git.Init(tmp)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(output)

}
