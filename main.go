package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/mpetavy/common"
	"io/ioutil"
	"strings"
)

var (
	filename = flag.String("f", "", "filename")
)

func init() {
	common.Init(false, "1.0.0", "", "", "2018", "test", "mpetavy", fmt.Sprintf("https://github.com/mpetavy/%s", common.Title()), common.APACHE, nil, nil, nil, run, 0)
}

func run() error {
	var ba []byte
	var err error

	if len(*filename) > 0 && common.FileExists_(*filename) {
		ba, err = ioutil.ReadFile(*filename)
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

	scanner := bufio.NewScanner(bytes.NewReader(ba))
	for scanner.Scan() {
		line := scanner.Text()

		line = strings.TrimSpace(line)

		if len(line) > 0 && !strings.Contains(line, "---") {
			splits := strings.Split(line, "|")

			if len(splits) == 0 {
				continue
			}

			st.AddRow()
			for _, split := range splits {
				st.AddCol(split)
			}
		}
	}

	if st == nil {
		st = common.NewStringTable()
		st.Markdown = true

		st.AddRow()
		for i := 0; i < 3; i++ {
			st.AddCol(fmt.Sprintf("header-%d", i))
		}
		st.AddRow()
		for i := 0; i < 3; i++ {
			st.AddCol(fmt.Sprintf("column-%d", i))
		}
	}

	fmt.Printf(st.String())

	if len(*filename) == 0 {
		err := clipboard.WriteAll(st.String())

		if common.Error(err) {
			return err
		}
	}

	return nil
}

func main() {
	defer common.Done()

	common.Run(nil)
}
