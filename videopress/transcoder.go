package videopress

import (
	"bytes"
	"fmt"
	"log"

	"os/exec"
)

func ConvertToWebm(filename string) (string, error) {
	log.Print("Encoding ", filename, " to webm")
	outFname := "out.webm"
	cmd := exec.Command("ffmpeg",
		"-i", filename,
		"-c:v", "libvpx",
		"-crf", "10",
		"-b:v", "1M",
		"-c:a", "libvorbis",
		outFname)

	cmd.Dir = "uploads"

	var outerr bytes.Buffer
	cmd.Stderr = &outerr

	//err := cmd.Run()
	out, err := cmd.Output()
	fmt.Printf("%s\n", outerr.String())

	if err != nil {
		//return "", err
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out)
	return "hello", nil
}
