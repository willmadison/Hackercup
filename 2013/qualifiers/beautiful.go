package main

import (
  "fmt"
  "os"
  "bufio"
  "strconv"
  "flag"
  "strings"
  "regexp"
  "sort"
)

var inputFileName *string = flag.String("file", "", "The input file for this script")

type Reverse struct {
	sort.Interface
}

func (r Reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

func main() {
  flag.Parse()

  if *inputFileName == "" {
    panic("No Input File Provided!")
  }

  inputFile, openError := os.Open(*inputFileName)

  if openError != nil {
    errorDesc := fmt.Sprint("Error opening %s for reading!", inputFileName)
    panic(errorDesc)
  }

  bufferedInput := bufio.NewReader(inputFile)

  var currentLine []byte

  currentLine, isTruncated, readError := bufferedInput.ReadLine()

  if readError != nil {
    fmt.Fprintf(os.Stdout, "Error reading from %s!: ", inputFileName, readError)
  }

  //First line of file is always the number of test cases...
  numCases,_ := strconv.Atoi(string(currentLine))

  //For each line in the file
  for caseNumber := 1; caseNumber <= numCases && readError == nil; caseNumber++ {
    currentLine, isTruncated, readError = bufferedInput.ReadLine()

    if isTruncated {
      fmt.Fprintf(os.Stdout, "Buffer was too small! Line was truncated to %s ", string(currentLine))
    } else {
      findMaxBeauty(caseNumber, strings.ToLower(string(currentLine)))
    }

  }

  defer inputFile.Close()
}

func findMaxBeauty(caseNumber int, currentLine string) {
  maxBeauty := 0

  var charactersToOccurrences = map[string]int{}

  for _, char := range currentLine {
    character := string(char)

    if isAlphabetic(character) {
      charactersToOccurrences[character] += 1
    }
  }

  occurrenceValues := []int{}

  for _, numOccurrences := range charactersToOccurrences {
    occurrenceValues = append(occurrenceValues , numOccurrences)
  }

  //Sort the occurrence values in reverse order
  sort.Sort(Reverse{sort.IntSlice(occurrenceValues)})

  topBeautyScore := 26

  for _, occurrences := range occurrenceValues {
    maxBeauty += topBeautyScore * occurrences
    topBeautyScore -= 1
  }

  fmt.Fprintf(os.Stdout, "Case #%d: %d\n", caseNumber, maxBeauty)
}

func isAlphabetic(character string) bool {
  matches, _ := regexp.MatchString("[[:alpha:]]", character)
  return matches
}

