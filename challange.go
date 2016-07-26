package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/arbovm/levenshtein"
)

// Reading vocabulary file
func main() {
	vocabulary := loadVocabulary("./vocabulary.txt")
	words := loadWords("./187")

	count := len(words)
	c := make(chan int, 30)
	for i, mul := range words {
		go func(word string, mul int) {
			length := len(word)
			wordDistance := length

			// Trying to find distance in main vocab
			if slice, ok := vocabulary[length]; ok {
				wordDistance = searchInSlice(slice, word)
			}

			// Now we will try to find distance
			shoulder := 1
			if wordDistance >= shoulder {
				for shoulder < wordDistance {
					leftDistance := wordDistance
					rightDistance := wordDistance

					if leftSlice, ok := vocabulary[length-shoulder]; ok {
						leftDistance = searchInSlice(leftSlice, word)
					}

					if leftDistance >= wordDistance {
						if rightSlice, ok := vocabulary[length+shoulder]; ok {
							rightDistance = searchInSlice(rightSlice, word)
						}

						if rightDistance < wordDistance {
							wordDistance = rightDistance
						}
					} else {
						wordDistance = leftDistance
					}

					shoulder++
				}
			}

			c <- wordDistance * mul
		}(i, mul)
	}

	totalDistane := 0
	for i := 1; i <= count; i++ {
		totalDistane = totalDistane + <-c
	}

	fmt.Println(totalDistane)
}

// Load vocabulary file
func loadVocabulary(filename string) map[int][]string {
	vocabulary := make(map[int][]string)

	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		word := strings.ToUpper(scanner.Text())
		wordLength := len(word)

		vocabulary[wordLength] = append(vocabulary[wordLength], word)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return vocabulary
}

// Load words file
func loadWords(filename string) map[string]int {
	words := make(map[string]int)
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		for i := range parts {
			word := strings.ToUpper(parts[i])
			if _, ok := words[word]; ok {
				words[word]++
			} else {
				words[word] = 1
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return words
}

// Searching in slice
func searchInSlice(slice []string, word string) int {
	wordDistance := len(word)

	for j := range slice {
		tmp := levenshtein.Distance(word, slice[j])

		if tmp < wordDistance {
			wordDistance = tmp
		}

		if wordDistance == 0 {
			break
		}
	}

	return wordDistance
}
