package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func main() {

	// Main flag
	cmd := flag.NewFlagSet("filter", flag.ExitOnError)

	// More flags
	redValue := cmd.Float64("red", 1, "Red scale")
	greenValue := cmd.Float64("green", 1, "Green scale")
	blueValue := cmd.Float64("blue", 1, "Blue scale")
	filePath := cmd.String("file", "", "File to process") // Mandatory flag

	if len(os.Args) < 2 {
		log.Fatalf("expected at least one command")
	}

	switch os.Args[1] {
	case "filter":
		handleCmd(cmd, redValue, greenValue, blueValue, filePath)
	default:
		log.Fatalf("expected [filter] subcommand, got [%s]", os.Args[1])
	}
	sep := strings.Split(*filePath, ".ppm")
	if len(sep) != 2 {
		log.Fatal("fail finding .ppm extension")
	}

	read, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("fail reading file: %v", err)
	}

	offset, err := GetOffset(read)
	if err != nil {
		log.Fatalf("fail decoding file: %v", err)
	}
	wg.Add(3)
	read2 := make([]byte, cap(read))
	copy(read2, read)
	read3 := make([]byte, cap(read))
	copy(read3, read)
	go BuildWorker(offset, "green", sep[0], *greenValue, read)
	go BuildWorker(offset, "blue", sep[0], *blueValue, read2)
	go BuildWorker(offset, "red", sep[0], *redValue, read3)
	wg.Wait()
	fmt.Println("Successfully filter")
}

func handleCmd(cmd *flag.FlagSet, red, green, blue *float64, file *string) {
	cmd.Parse(os.Args[2:])
	if *file == "" {
		log.Fatalf("file argument is required")
	}
}

// GetOffset reads fron slice of bytes and return offset value
func GetOffset(bytes []byte) (int, error) {
	var head int
	// Capture version
	if string(bytes[0:2]) != "P6" {
		return 0, fmt.Errorf("file version not suported")
	}
	// Capture Comment
	if bytes[3] == 0x23 {
		for i, b := range bytes[3:] {
			if b == 0x0A {
				head = i + 4
				break
			}
		}
	} else {
		head = 3
	}
	// Capture width & Length
	for i, b := range bytes[head:] {
		if b == 0x0A {
			head += i + 1
			break
		}
	}
	// Capture depth
	for i, b := range bytes[head:] {
		if b == 0x0A {
			head += i + 1
			break
		}
	}
	return head, nil
}

// BuildWorker this is a go routine to write filtred images
func BuildWorker(offset int, color, path string, intense float64, bytes []byte) {
	cut := map[string]int{
		"red":   0,
		"green": 1,
		"blue":  2,
	}
	name := fmt.Sprintf("%s_%s.ppm", path, color)
	for i, b := range bytes[offset:] {
		var newB byte
		if i%3 == cut[color] {
			value := int(float64(b) * intense)
			if value >= 255 {
				value = 255
			}
			newB = byte(value)
		}
		bytes[offset+i] = newB
	}
	err := os.WriteFile(name, bytes, 0677)
	if err != nil {
		log.Fatalf("error creating file: %v", err)
	}
	wg.Done()
}
