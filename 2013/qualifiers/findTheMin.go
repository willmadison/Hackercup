package main

import (
  "fmt"
  "os"
  "bufio"
  "strconv"
  "flag"
  "strings"
  "sort"
  "time"
)

var inputFileName *string = flag.String("file", "", "The input file for this script")

type Case struct {
  a,b,c,r,numNonNegatives,numKnownValues int
}

type Set struct {
  set map[int]int
}

func (s *Set) Add(x int) {
  s.set[x] += 1
}

func (s *Set) Remove(x int) {
  s.set[x] -= 1

  if s.set[x] <= 0 {
    delete(s.set, x)
  }
}

func (s *Set) Contains(x int) bool {
  _, isPresent := s.set[x]

  return isPresent
}

func (s *Set) Values() []int {
  values := []int{}

  for entry,_ := range s.set {
    values = append(values, entry)
  }

  sort.Ints(values)

  return values
}

func main() {
  start := time.Now()

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
  linesPerCase := 2

  //For each line in the file
  for caseNumber := 1; caseNumber <= numCases && readError == nil; caseNumber++ {
    problemCase := Case{}

    for line := 0; line < linesPerCase; line++ {
      currentLine, isTruncated, readError = bufferedInput.ReadLine()

      if isTruncated {
        fmt.Fprintf(os.Stdout, "Buffer was too small! Line was truncated to %s ", string(currentLine))
      } else {
        fields := strings.Fields(string(currentLine))

        switch (line) {
        case 0: //First Line of the Pair
          problemCase.numNonNegatives,_ = strconv.Atoi(fields[0])
          problemCase.numKnownValues,_ = strconv.Atoi(fields[1])
        case 1: //Second Line of the Pair
          problemCase.a,_ = strconv.Atoi(fields[0])
          problemCase.b,_ = strconv.Atoi(fields[1])
          problemCase.c,_ = strconv.Atoi(fields[2])
          problemCase.r,_ = strconv.Atoi(fields[3])
        }
      }
    }

    findTheMin(caseNumber, problemCase)
  }

  defer inputFile.Close()

  fmt.Fprintf(os.Stdout, "Total Runtime: %v\n", time.Now().Sub(start))
}

func findTheMin(caseNumber int, problemCase Case) {
  //fmt.Fprintf(os.Stdout, "Case #%d: Case Details: %+v\n", caseNumber, problemCase)

  m := []int{problemCase.a} //From problem statement m[0] = a

  //Generate the known values...
  for i := 1; i < problemCase.numKnownValues; i++ {
    nextValue := (problemCase.b * m[i - 1] + problemCase.c) % problemCase.r
    m = append(m, nextValue)
  }

  candidates := &Set{set : map[int]int{}}

  for _, entry := range m {
    candidates.Add(entry)
  }

  pool := []int{}
  generatedValues := []int{}

  totalInPool := 0
  poolEntry := 0

  //Fill our pool up with exactly k + 1 entries.
  for ; poolEntry < problemCase.r; poolEntry++ {
    if !candidates.Contains(poolEntry) {
      pool = append(pool, poolEntry)
      totalInPool++
    }

    if totalInPool == problemCase.numKnownValues + 1 {
      break
    }
  }

  //If our pool is empty simply fill it with much larger values which we'll replace the 0th element of each iteration...
  if totalInPool == 0 {
    poolEntry = problemCase.r

    for totalInPool < problemCase.numKnownValues + 1 {
      pool = append(pool, poolEntry)
      totalInPool++
      poolEntry++
    }
  }

  cycleComplete := false

  var nthValue int

  //Find the nth item...
  for j := problemCase.numKnownValues; j < problemCase.numNonNegatives; {
    if !cycleComplete {
      currentSize := len(m)

      nthValue = pool[0]

      if len(pool) > 1 {
        pool = pool[1:]
      }

      m = append(m, nthValue)
      generatedValues = append(generatedValues, nthValue)
      cycleComplete = len(generatedValues) == problemCase.numKnownValues + 1

      //Update our set of candidates...
      windowLowerBound := currentSize - problemCase.numKnownValues
      window := m[windowLowerBound:]

      newlyEligibleCandidate := window[0]

      candidates.Remove(newlyEligibleCandidate)
      candidates.Add(nthValue)

      if !candidates.Contains(newlyEligibleCandidate) && newlyEligibleCandidate >= 0 {
        if newlyEligibleCandidate < pool[0] {
          pool[0] = newlyEligibleCandidate
        } else {
          pool = append(pool, newlyEligibleCandidate)

          if needsSort := newlyEligibleCandidate < pool[len(pool)-1]; needsSort {
            sort.Ints(pool)
          }
        }
      }

      j++
    } else {
      //We can simply cycle through the values we've generated thus far as minumums will begin to repeat.
      for _, next := range generatedValues {
        nthValue = next
        j++

        if j == problemCase.numNonNegatives {
          break
        }
      }
    }
  }

  fmt.Fprintf(os.Stdout, "Case #%d: %d\n", caseNumber, nthValue)
}
