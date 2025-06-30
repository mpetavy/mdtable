package main

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/csv"
	"flag"
	"github.com/atotto/clipboard"
	"github.com/mpetavy/common"
	"os"
	"strings"
)

var (
	inputFilename  = flag.String("i", "", "input filename")
	outputFilename = flag.String("o", "", "output filename")
	format         = flag.String("format", "markdown", "format of output 'markdown', 'table', 'html'")
)

//go:embed go.mod
var resources embed.FS

func init() {
	common.Init("", "", "", "", "test", "", "", "", &resources, nil, nil, run, 0)
}

func run() error {
	var ba []byte
	var err error

	if *inputFilename != "" {
		ba, err = os.ReadFile(*inputFilename)
		if common.Error(err) {
			return err
		}
	} else {
		t, err := clipboard.ReadAll()
		if err != nil {
			return err
		}

		if len(t) > 0 {
			ba = []byte(t)
		}
	}

	if ba == nil {
		common.Info("Please provide content via clipboard or file")

		return nil
	}

	st := common.NewStringTable()

	output := strings.Builder{}

	comma := ""

	scanner := bufio.NewScanner(bytes.NewReader(ba))
	if scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "\t") {
			comma = "\t"
		}

		if strings.Contains(line, ",") {
			comma = ","
		}
	}

	if comma != "" {
		common.Info("Reading as CSV")
		common.Info("")

		c := csv.NewReader(bytes.NewReader(ba))
		c.Comma = rune(comma[0])

		recs, err := c.ReadAll()
		if common.Error(err) {
			return err
		}

		for row := range len(recs) {
			st.AddRow()
			for col := range len(recs[row]) {
				st.AddCol(strings.TrimSpace(recs[row][col]))
			}
		}
	} else {
		common.Info("Reading as strings")
		common.Info("")

		crlf, err := common.NewSeparatorSplitFunc(nil, []byte("\n"), false)
		if common.Error(err) {
			return err
		}

		cell := strings.Builder{}

		inMyTable := false

		scanner := bufio.NewScanner(bytes.NewReader(ba))
		scanner.Split(crlf)
		for scanner.Scan() {
			line := scanner.Text()

			switch {
			case strings.HasPrefix(line, "++TABLE"):
				if inMyTable && cell.Len() > 0 {
					st.AddCol(cell.String())
				}

				cell.Reset()

				inMyTable = !inMyTable

				if inMyTable {
					st.NoHeader = true
				} else {
					if st.Rows() > 0 {
						output.WriteString(st.Markdown())

						st.Clear()
						st.NoHeader = false
					}
				}

				continue
			case strings.HasPrefix(line, "++ROW"):
				if cell.Len() > 0 {
					st.AddCol(cell.String())
				}

				cell.Reset()

				st.AddRow()

				continue
			case strings.HasPrefix(line, "++COL"):
				if cell.Len() == 0 {
					cell.WriteString("<na>")
				}

				st.AddCol(cell.String())

				cell.Reset()

				continue
			}

			if inMyTable {
				if cell.Len() > 0 {
					cell.WriteString("<br>")
				}

				cell.WriteString(strings.TrimSpace(line))

				continue
			}

			if !strings.Contains(line, "|") {
				if st.Rows() > 0 {
					output.WriteString(st.Markdown())

					st.Clear()
					st.NoHeader = false
				}

				output.WriteString(line)

				continue
			}

			line = strings.TrimSpace(line)

			if strings.HasPrefix(line, "|") {
				line = line[1:]
			}

			if strings.HasSuffix(line, "|") {
				line = line[:len(line)-1]
			}

			if len(line) > 0 && !strings.Contains(line, "---") {
				splits := common.Split(line, "|")

				if len(splits) == 0 {
					continue
				}

				st.AddRow()
				for _, split := range splits {
					split = strings.TrimSpace(strings.ReplaceAll(split, "|", ""))
					st.AddCol(strings.TrimSpace(split))
				}
			}
		}
	}

	if st.Rows() > 0 {
		switch *format {
		case "table":
			output.WriteString(st.Table())
		case "csv":
			output.WriteString(st.CSV())
		case "html":
			output.WriteString(st.HTML())
		case "markdown":
			output.WriteString(st.Markdown())
		default:
			output.WriteString(st.JSON())
		}
	}

	err = clipboard.WriteAll(output.String())
	if common.Error(err) {
		return err
	}

	if *outputFilename != "" {
		err := common.FileBackup(*outputFilename)
		if common.Error(err) {
			return err
		}

		err = os.WriteFile(*outputFilename, []byte(output.String()), common.DefaultFileMode)
		if common.Error(err) {
			return err
		}
	}

	common.Info(output.String())

	return nil
}

func main() {
	common.Run(nil)
}
