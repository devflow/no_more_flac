package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var dir string
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func main() {
	rand.Seed(time.Now().UnixNano())

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	dir = filepath.Dir(ex)

	if _, err := os.Stat(dir + "\\ffmpeg.exe"); os.IsNotExist(err) {
		fmt.Println("missing ffmpeg.exe in same location.")
	}

	err = filepath.Walk(dir, startConvertFile)

	fmt.Println("All Done.")
}

func startConvertFile(path string, f os.FileInfo, err error) error {
	if f.IsDir() || strings.ToLower(filepath.Ext(path)) != ".flac" {
		return nil
	}

	absPath, err := filepath.Rel(dir, path)

	if err != nil {
		return err
	}

	targetName := filepath.Join(filepath.Dir(path), strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))+".mp3")

	if _, err := os.Stat(targetName); os.IsNotExist(err) {
		fmt.Printf("Converting... %s", absPath)
	} else {
		fmt.Printf("already exists file as same name (mp3) : %s\n", absPath)
		return nil
	}

	randName := filepath.Join(filepath.Dir(path), RandStringRunes(32)+".mp3")

	if err = convertFlacToMp3(path, randName); err != nil {
		fmt.Printf(".... Failed to convert audio. \n")
	}

	if err = os.Rename(randName, targetName); err != nil {
		fmt.Printf(".... Failed to move new file \n")
	}

	if err = os.Remove(path); err != nil {
		fmt.Printf(".... Failed to remove original flac file... ")

		if err = os.Remove(targetName); err != nil {
			fmt.Printf("failed to remove original flac file \n")
		} else {
			fmt.Printf("removed converted file. \n")
		}
	}

	fmt.Printf(".... DONE \n")

	return nil
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func convertFlacToMp3(in string, out string) error {
	return exec.Command(dir+"\\ffmpeg.exe", "-i", in, "-codec:a", "libmp3lame", "-b:a", "320k", "-map_metadata", "0", "-id3v2_version", "3", "-y", out).Run()
}
