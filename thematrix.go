package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"strings"
	"time"

	//"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/nullboundary/glfont"
)

const (
	/*
	windowWidth  = 1920
	windowHeight = 1080
	numRows = 29
	numCols = 80
	*/
	windowWidth  = 2560
	windowHeight = 1440
	numRows = 39
	numCols = 106

	startCharCode = 33
	endCharCode = 78
	fontSize = 46
	colFallingTicks = 20
	tickTime = float64(0.04)
	//tickTime = float64(0.2)
	startDelay = 5
	totalAnimTime = 30
	randomChangeChance = 0.2
	minTicksToRandomChange = 3
	maxTicksToRandomChange = 20
	numEmptySingleCols = 7
	numEmptyDoubleCols = 4
	numHighlightCols = 12
)

var (
	MatrixCharset string
	Ticks int
	MinColsSeparation int
	MaxColsSeparation int
	EndingTicks int
)

func getRandomCharacter() byte {
	charCode := rand.Intn(endCharCode - startCharCode + 1) + startCharCode

	return byte(charCode)

	/*
	if charCode != 37 {
		return string(rune(charCode))
	} else {
		return "%%"
	}
	*/
}

func init() {
	rand.Seed(time.Now().UnixNano())

	runtime.LockOSThread()
}

func main() {

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, _ := glfw.CreateWindow(int(windowWidth), int(windowHeight), "glfontExample", nil, nil)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)
	
	if err := gl.Init(); err != nil { 
		panic(err)
	}

	//load font (fontfile, font scale, window width, window height
	font, err := glfont.LoadFont("MatrixCodeNFI.ttf", int32(fontSize), windowWidth, windowHeight)
	if err != nil {
		log.Panicf("LoadFont: %v", err)
	}

	MatrixCharset := ""

	for i := startCharCode; i <= endCharCode; i++ {
		if i != 37 {
			currCharStr := string(rune(i))
			MatrixCharset += currCharStr
		} else {
			//currCharStr := string(rune(79))
			MatrixCharset += "%%"
		}
	}

	fmt.Println(MatrixCharset)

	EndingTicks = int((float64(1) / tickTime) * totalAnimTime)

	colIndices := make([]int, numCols)
	for i := range colIndices {
	    colIndices[i] = i
	}
	rand.Shuffle(len(colIndices), func(i, j int) { colIndices[i], colIndices[j] = colIndices[j], colIndices[i] })

	MaxColsSeparation = (numCols / (numEmptySingleCols + numEmptyDoubleCols + 2))
	fmt.Println("MaxColsSeparation: ", MaxColsSeparation)
	MinColsSeparation = int(math.Max(0, float64(MaxColsSeparation) - float64(1)))

	emptyColsCounts := make([]int, numEmptySingleCols + numEmptyDoubleCols)

	for i := 0; i < numEmptySingleCols; i++ {
		emptyColsCounts[i] = 1
	}
	for i := numEmptySingleCols; i < numEmptySingleCols + numEmptyDoubleCols; i++ {
		emptyColsCounts[i] = 2
	}

	rand.Shuffle(len(emptyColsCounts), func(i, j int) { emptyColsCounts[i], emptyColsCounts[j] = emptyColsCounts[j], emptyColsCounts[i] })

	emptyColsBitmap := make([]int, numCols)

	currEmptyCol := rand.Intn(MaxColsSeparation - MinColsSeparation + 1) + MinColsSeparation
	for i := range emptyColsCounts {
		currEmptyCol += rand.Intn(MaxColsSeparation - MinColsSeparation + 1) + MinColsSeparation

		if (emptyColsCounts[i] == 1) {
			emptyColsBitmap[currEmptyCol] = 1
		} else if (emptyColsCounts[i] == 2) {
			emptyColsBitmap[currEmptyCol] = 1
			emptyColsBitmap[currEmptyCol + 1] = 1
		}
	}

	fmt.Println(emptyColsCounts)

	colIndices = make([]int, numCols)
	for i := range colIndices {
	    colIndices[i] = i
	}
	rand.Shuffle(len(colIndices), func(i, j int) { colIndices[i], colIndices[j] = colIndices[j], colIndices[i] })

	var highlightColsToSwitch []int
    var availiableColsForHighlight []int

	highlightColsBitmap := make([]int, numCols)

	for i := 0; i < numHighlightCols; {
		if emptyColsBitmap[i] == 0 {
			highlightColsBitmap[colIndices[i]] = 1
			i++
		}
	}

	var emptyRow = make([]byte, numCols)
	for i := range emptyRow {
	    emptyRow[i] = ' '
	}

	var highlightRow = make([]byte, numCols)

	var codeRows [numRows][numCols]byte
	var realCodeRows [numRows][numCols]byte

	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			if (emptyColsBitmap[j] == 1) {
				codeRows[i][j] = ' '
			} else {
				codeRows[i][j] = getRandomCharacter()
			}
		}
	}

	var columnEnds [numCols]int
	var columnLens [numCols]int
	var columnStartTicks [numCols]int

	for i := 0; i < numCols; i++ {
		columnStartTicks[i] = rand.Intn(100) + int((float64(1) / tickTime) * startDelay)
		columnEnds[i] = 0
		columnLens[i] = numRows + rand.Intn(int(math.Round(numRows * 1.5)))
	}

	var secondColumnEnds [numCols]int
	var secondColumnLens [numCols]int
	var secondColumnStartTicks [numCols]int

	for i := 0; i < numCols; i++ {
		secondColumnStartTicks[i] = columnStartTicks[i] + columnLens[i] + rand.Intn(int(math.Round(numRows * 1.5)))
		secondColumnEnds[i] = 0
		secondColumnLens[i] = numRows + rand.Intn(int(math.Round(numRows * 1.5)))
	}

	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			realCodeRows[i][j] = ' '
		}
	}

	
	//fmt.Println("First row: ", string(codeRows[0][:]))

	lastRandomChangeTicks := 0
	ticksUntilNextRandomChange := rand.Intn(maxTicksToRandomChange - minTicksToRandomChange + 1) + minTicksToRandomChange


	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		Ticks = int(glfw.GetTime() / tickTime)
		//fmt.Printf("\rTicks: %d", Ticks)

		if Ticks - lastRandomChangeTicks > ticksUntilNextRandomChange {
			for i := 0; i < numRows; i++ {
				for j := 0; j < numCols; j++ {
					if (emptyColsBitmap[j] == 1) {
						codeRows[i][j] = ' '
					} else if rand.Float64() < randomChangeChance {
						codeRows[i][j] = getRandomCharacter()
					}
				}
			}	

			lastRandomChangeTicks = Ticks
			ticksUntilNextRandomChange = rand.Intn(maxTicksToRandomChange - minTicksToRandomChange + 1) + minTicksToRandomChange
		}

		for colIndex := 0; colIndex < numCols; colIndex++ {
			//columnLens[colIndex] = int(float32(Ticks - columnStartTicks[colIndex]) / float32(colFallingTicks))
			columnEnds[colIndex] = Ticks - columnStartTicks[colIndex]

			secondColumnEnds[colIndex] = Ticks - secondColumnStartTicks[colIndex]

			//fmt.Printf("Column lengths: %d, %d, %d, %d, %d\n", columnLens[0], columnLens[1], columnLens[2], columnLens[3], columnLens[4])
		}

		// Clear out both of the highlight switch buffers
		highlightColsToSwitch = nil
		availiableColsForHighlight = nil


		if Ticks < EndingTicks {
		// Update strands that have fallen past the screen
		// Also change highlight columns if we can
			for col := 0; col < numCols; col++ {
				// Primary strands
				if columnEnds[col] - columnLens[col] > numRows - 1 {
					//fmt.Println("Primary strand at index", col, "fell off the screen!");
					columnStartTicks[col] = secondColumnStartTicks[col] + secondColumnLens[col] + rand.Intn(int(math.Round(numRows * 1.5)))
					columnEnds[col] = 0
					columnLens[col] = numRows + rand.Intn(int(math.Round(numRows * 1.5)))
				}

				// Secondary strands
				if secondColumnEnds[col] - secondColumnLens[col] > numRows - 1 {
					//fmt.Println("Secondary strand at index", col, "fell off the screen!");
					secondColumnStartTicks[col] = columnStartTicks[col] + columnLens[col] + rand.Intn(int(math.Round(numRows * 1.5)))
					secondColumnEnds[col] = 0
					secondColumnLens[col] = numRows + rand.Intn(int(math.Round(numRows * 1.5)))
				}

				if columnEnds[col] > numRows {
					if highlightColsBitmap[col] == 1 {
						highlightColsToSwitch = append(highlightColsToSwitch, col)
					} else {
						availiableColsForHighlight = append(availiableColsForHighlight, col)
					}
				}
			}
		}

		if availiableColsForHighlight != nil {
			rand.Shuffle(len(availiableColsForHighlight), func(i, j int) { availiableColsForHighlight[i], availiableColsForHighlight[j] = availiableColsForHighlight[j], availiableColsForHighlight[i] })
		}

		highlightColsSwitchLen := int(math.Min(float64(len(availiableColsForHighlight)), float64(len(highlightColsToSwitch))))

		for i := 0; i < highlightColsSwitchLen; i++ {
			highlightColsBitmap[highlightColsToSwitch[i]] = 0
			highlightColsBitmap[availiableColsForHighlight[i]] = 1
		}


		for i := 0; i < numRows; i++ {
			for j := 0; j < numCols; j++ {
				// First part of each expression makes sure the start of the strand isn't already past the row
				// Second part of each expression makes sure the end of the strand isn't before the row
				if (columnEnds[j] - columnLens[j] > i || columnEnds[j] < i) &&
				   (secondColumnEnds[j] - secondColumnLens[j] > i || secondColumnEnds[j] < i) {
					realCodeRows[i][j] = ' '
				} else {
					realCodeRows[i][j] = codeRows[i][j]
				}
			}	
		}


		for i := 0; i < numCols; i++ {
			columnEndRow := columnEnds[i]
			if highlightColsBitmap[i] == 1 && columnEndRow >= 0 && columnEndRow < numRows {
				//fmt.Println("Highlighting column:", i)
				font.SetColor(0.8, 1.0, 0.9, 1.0)

				copy(highlightRow, emptyRow)
				highlightRow[i] = codeRows[columnEndRow][i]
				//fmt.Println("Highlight Row:", string(highlightRow[:]))
				//fmt.Printf("Setting highlightRow column to index: (%d, %d)\n", columnEndRow, i)
				//fmt.Println("The codeRows value is:", codeRows[columnEndRow][i])

				highlightRowStr := string(highlightRow[:])

				highlightRowStr = strings.Replace(highlightRowStr, "%", "%%", -1)
				font.Printf(0, float32((float64(fontSize) * 0.8) + float64(fontSize) * (0.8 * float64(columnEndRow))), 1.0, highlightRowStr)
			}
		}

     	//set color and draw text
		font.SetColor(0.2, 0.851, 0.376, 1.0) //r,g,b,a font color
		//font.SetColor(0.631, 1.0, 0.796, 1.0) // brighter character

		for i := 0; i < numRows; i++ {
			currRow := string(realCodeRows[i][:])
			currRow = strings.Replace(currRow, "%", "%%", -1)
			font.Printf(0, float32((float64(fontSize) * 0.8) + float64(fontSize) * (0.8 * float64(i))), 1.0, currRow)
		}

		//font.Printf(0, 48, 1.0, MatrixCharset) //x,y,scale,string,printf args

		//fmt.Println(glfw.GetTime())

		

		window.SwapBuffers()
		glfw.PollEvents()

	}

	//fmt.Println("First row (actual): ", string(realCodeRows[0][:]))
}