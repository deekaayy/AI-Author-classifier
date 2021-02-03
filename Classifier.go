package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)



func CountWords(fileName string,restriction int) map[string] int {
	// If the restriction parameter is -1, it means that we want to map the whole file.
	// -1 is like a flag, if we pass -1 as a restriction it means we want every word to be in the returned map
	myMap := make(map[string]int)
	var tokenizedString []string

	input,err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		tokenizedString = append(tokenizedString[:], strings.Fields(scanner.Text())...)
	}
	for i:= range tokenizedString {
		myMap[tokenizedString[i]]++
	}

	words := make([]string, 0, len(myMap))
	for name := range myMap {
		words = append(words, name)
	}

	sort.Slice(words, func(i, j int) bool {
		return myMap[words[i]] > myMap[words[j]]
	})


	if restriction == -1 {
		return myMap
	}
	ToReturnMap := make(map[string] int,restriction)

	iterations:= 0
	for _, k := range words {
		ToReturnMap[k] = myMap[k]
		if iterations > restriction-2 {
			break
		}
		iterations++
	}
	return ToReturnMap
}

// Given the two counters of the authors, print out which one has more positives.
func FinalVerdict(aleko int, zahari int) {
	sum := aleko + zahari

	if aleko > zahari {
		curAleko:= float64(aleko)/float64(sum) * 100
		fmt.Printf("Aleko Konstantinov matching possibility: %.f%%",curAleko)
	} else {
		curZahari:= float64(zahari)/float64(sum) * 100
		fmt.Printf("Zahari Stoyanov matching possibility: %.f%%",curZahari)
	}
	fmt.Println()
}

func Classify() {
	// filling the maps with the appropriate data
	AlekoStopWords := CountWords(`C:\Users\Ivan\go\src\week2_part1\classifier\aleko.txt`,10000)
	ZahariStopWords := CountWords(`C:\Users\Ivan\go\src\week2_part1\classifier\zahari.txt`,10000)
	InputMapped := CountWords(`C:\Users\Ivan\go\src\week2_part1\classifier\input.txt`,-1)

	AlekoY := 0
	ZahariY := 0

	for i:= range InputMapped {
		// if neither of both authors have an occurrence of this word, just go to the next one
		if AlekoStopWords[i] == 0 && ZahariStopWords[i] == 0 {
			continue
		}
		// if Aleko Konstantinov doesn't use this word at all, increment the Zahari Stoyanov counter twice
		// (this approach seems to be more consistent in terms of guessing the right author)
		if AlekoStopWords[i] == 0 {
			ZahariY+=2
			continue
		}
		// same as the above
		if ZahariStopWords[i] == 0 {
			AlekoY+=2
			continue
		}
		// here we compute the frequency of the current word by both authors
		// for example AlekoFreq = 0.0032; ZahariFreq = 0.0030
		AlekoFreq := float64(InputMapped[i]) / float64(AlekoStopWords[i])
		ZahariFreq := float64(InputMapped[i]) / float64(ZahariStopWords[i])

		// if one of the authors have a higher frequency we increment his counter
		if AlekoFreq < ZahariFreq {
			ZahariY++
		} else if AlekoFreq > ZahariFreq {
			AlekoY++
		} else {
			AlekoY++
			ZahariY++
		}
	}
	fmt.Printf("AlekoY: %d\nZahariY: %d\n",AlekoY,ZahariY)
	FinalVerdict(AlekoY,ZahariY)
}

func main() {
	Classify()
}

// Sources for Aleko Konstantinov --> http://www.slovo.bg/showauthor.php3?ID=169&LangID=1
// Sources for Zahari Stoyanov --> http://www.slovo.bg/showauthor.php3?ID=149&LangID=1