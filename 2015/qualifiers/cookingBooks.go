package main

import (
  "fmt"
  "os"
  "bufio"
  "strconv"
  "flag"
  "strings"
)

var inputFileName *string = flag.String("file", "", "The input file for this script")

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

  defer inputFile.Close()

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
      cookTheBooks(caseNumber, string(currentLine))
    }

  }
}

func cookTheBooks(caseNumber int, currentLine string) {
    minValue, maxValue := 0, 0

    numDigits := len(currentLine)

    if (numDigits > 1) {
        digits := strings.Split(currentLine, "")

        _, maxLocation := findRightmostMaxValue(digits)

        swapDigits(digits, maxLocation, 0)

        maxValue, _ = strconv.Atoi(strings.Join(digits, ""))

        digits = strings.Split(currentLine, "")

        minFound, minLocation := findRightmostMinValue(digits)

        firstLargerIndex := firstIndexLargerThan(digits, minFound)

        if minLocation > 0 && firstLargerIndex != -1 {

            if minFound == 0 && firstLargerIndex == 0 {
                firstLargerIndex += 1
            }

            swapDigits(digits, minLocation, firstLargerIndex)
        }

        minValue, _ = strconv.Atoi(strings.Join(digits, ""))
    } else {
        givenValue, _ := strconv.Atoi(currentLine)
        minValue, maxValue = givenValue, givenValue
    }

    fmt.Fprintf(os.Stdout, "Case #%d: %d %d\n", caseNumber, minValue, maxValue)
}

func findRightmostMaxValue(digits []string) (int, int) {
    maxValue, location := 0, -1

    intValue := 0

    for index, value := range digits {
        intValue, _ = strconv.Atoi(value)

        if intValue >= maxValue {
            maxValue, location = intValue, index
        }
    }

    return maxValue, location
}

func findRightmostMinValue(digits []string) (int, int) {
    minValue, location := 9, -1

    intValue := 0

    for index, value := range digits {
        intValue, _ = strconv.Atoi(value)

        if intValue <= minValue {
            minValue, location = intValue, index
        }
    }

    return minValue, location
}

func firstIndexLargerThan(digits []string, threshold int) int {
    firstLargerIndex := -1

    intValue := 0

    for index, value := range digits {
        intValue, _ = strconv.Atoi(value)

        if intValue > threshold {
            firstLargerIndex = index
            break
        }
    }

    return firstLargerIndex
}


func swapDigits(digits []string, first, second int) {
    digits[first], digits[second] = digits[second], digits[first]
}
