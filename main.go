package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/mpetavy/common"
	"os"
	"strings"
)

var (
	filename = flag.String("f", "", "filename")
)

func init() {
	common.Init("mdtable", "", "", "", "2018", "test", "mpetavy", fmt.Sprintf("https://github.com/mpetavy/%s", common.Title()), common.APACHE, nil, nil, nil, run, 0)
}

func run() error {
	var ba []byte
	var err error

	if len(*filename) > 0 && common.FileExists_(*filename) {
		ba, err = os.ReadFile(*filename)
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
	st.Markdown = true

	crlf, err := common.NewSeparatorSplitFunc(nil, []byte("\n"), false)
	if common.Error(err) {
		return err
	}

	output := strings.Builder{}

	scanner := bufio.NewScanner(bytes.NewReader(ba))
	scanner.Split(crlf)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.Contains(line, "|") {
			if st.Rows() > 0 {
				output.WriteString(st.String())

				st.Clear()
			}

			output.WriteString(line)

			continue
		}

		line = strings.TrimSpace(line)

		if strings.HasSuffix(line, "|") {
			line = line[:len(line)-1]
		}

		if len(line) > 0 && !strings.Contains(line, "---") {
			splits := strings.Split(line, "|")

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

	if st.Rows() > 0 {
		output.WriteString(st.String())
	}

	if len(*filename) > 0 {
		err := common.FileBackup(*filename)
		if common.Error(err) {
			return err
		}

		err = os.WriteFile(*filename, []byte(output.String()), common.DefaultFileMode)
		if common.Error(err) {
			return err
		}
	} else {
		err := clipboard.WriteAll(output.String())
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
