package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Helper struct {
	Red, Green, Blue   float64
	Size               int
	FilePath, FileName string
}

type ImageDescriptor struct {
	Version, Comment     string
	Width, Length, Depth int
	Head                 int
	Filter               Helper
}

func main() {

	// Main flag
	cmd := flag.NewFlagSet("filter", flag.ExitOnError)

	// More flags
	redValue := cmd.Float64("red", 1, "Red scale")
	greenValue := cmd.Float64("green", 1, "Green scale")
	blueValue := cmd.Float64("blue", 1, "Blue scale")
	sizeValue := cmd.Int("size", 256, "Reading size")
	filePath := cmd.String("file", "", "File to process") // Mandatory flag

	if len(os.Args) < 2 {
		log.Fatalf("expected at least one command")
	}

	switch os.Args[1] {
	case "filter":
		handleCmd(cmd, redValue, greenValue, blueValue, sizeValue, filePath)
	default:
		log.Fatalf("expected [filter] subcommand, got [%s]", os.Args[1])
	}
	sep := strings.Split(*filePath, ".ppm")
	if len(sep) != 2 {
		log.Fatal("fail finding .ppm extension")
	}

	helper := Helper{*redValue, *greenValue, *blueValue, *sizeValue, *filePath, sep[0]}

	read, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("fail reading file: %v", err)
	}

	var imageDesciptor ImageDescriptor
	err = DecodeFile(&read, &imageDesciptor)
	if err != nil {
		log.Fatalf("fail decoding file: %v", err)
	}
	imageDesciptor.Filter = helper
	read2 := make([]byte, cap(read))
	read3 := make([]byte, cap(read))
	copy(read2, read)
	copy(read3, read)
	BuildWorker("green", 1, imageDesciptor, read)
	BuildWorker("blue", 1, imageDesciptor, read2)
	BuildWorker("red", 1, imageDesciptor, read3)

	// time.Sleep(time.Second * 10)
}

func handleCmd(cmd *flag.FlagSet, red, green, blue *float64, size *int, file *string) {
	cmd.Parse(os.Args[2:])
	if *file == "" {
		log.Fatalf("file argument is required")
	}
}

func DecodeFile(bytes *[]byte, imageDesciptor *ImageDescriptor) error {
	var head int
	// Capture version
	imageDesciptor.Version = string((*bytes)[0:2])
	// Capture Comment
	if (*bytes)[3] == 0x23 {
		for i, b := range (*bytes)[3:] {
			if b == 0x0A {
				head = i + 4
				break
			}
			imageDesciptor.Comment += string(b)
		}
	}
	// Capture width & Length
	var widthLength string
	for i, b := range (*bytes)[head:] {
		if b == 0x0A {
			head += i + 1
			break
		}
		widthLength += string(b)
	}

	wL := strings.Split(widthLength, " ")
	w, err := strconv.Atoi(wL[0])
	if err != nil {
		return fmt.Errorf("failed getting width value: %s", err.Error())
	}
	l, err := strconv.Atoi(wL[1])
	if err != nil {
		return fmt.Errorf("failed getting length value: %s", err.Error())
	}
	imageDesciptor.Width = w
	imageDesciptor.Length = l
	// Capture depth
	var depth string
	for i, b := range (*bytes)[head:] {
		if b == 0x0A {
			head += i + 1
			break
		}
		depth += string(b)
	}
	d, err := strconv.Atoi(depth)
	if err != nil {
		return fmt.Errorf("failed getting depth value: %s", err.Error())
	}
	imageDesciptor.Depth = d
	imageDesciptor.Head = head
	return nil
}

func BuildWorker(color string, intense float64, imageD ImageDescriptor, bytes []byte) error {
	cut := map[string]int{
		"red":   0,
		"green": 1,
		"blue":  2,
	}
	head := bytes[:imageD.Head]
	body := bytes[imageD.Head:]
	name := fmt.Sprintf("%s_%s.ppm", imageD.Filter.FileName, color)
	for i, b := range body {
		var newB byte
		if i%3 == cut[color] {
			value := int(float64(b) * intense)
			if value >= 255 {
				value = 255
			}
			newB = byte(value)
		}
		body[i] = newB
	}
	head = append(head, body...)
	return os.WriteFile(name, head, 0677)
}
