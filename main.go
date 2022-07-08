package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mpetavy/common"
	"io/ioutil"
)

var (
	filename = flag.String("f", "", "filename")
)

func init() {
	common.Init(false, "1.0.0", "", "", "2018", "test", "mpetavy", fmt.Sprintf("https://github.com/mpetavy/%s", common.Title()), common.APACHE, nil, nil, nil, run, 0)
}

func run() error {
	if !common.FileExists_(*filename) {
		var tb [2][3]string
		for i := 0; i < len(tb[0]); i++ {
			tb[0][i] = fmt.Sprintf("header%d", i)
			tb[1][i] = fmt.Sprintf("sample%d", i)
		}

		ba, err := json.MarshalIndent(tb, "", "    ")
		if common.Error(err) {
			return err
		}

		err = ioutil.WriteFile(*filename, ba, common.DefaultFileMode)
		if common.Error(err) {
			return err
		}

		common.Info(string(ba))
		common.Info("file %s created", *filename)

		return nil
	}

	ba, err := ioutil.ReadFile(*filename)
	if common.Error(err) {
		return err
	}

	var tb [][]string

	err = json.Unmarshal(ba, &tb)
	if common.Error(err) {
		return err
	}

	st := common.NewStringTable()
	st.NoCross = true

	for row := 0; row < len(tb); row++ {
		st.AddRow()
		for col := 0; col < len(tb[row]); col++ {
			st.AddCol(tb[row][col])
		}
	}

	common.Info(st.String())

	return nil
}

func main() {
	defer common.Done()

	common.Run([]string{"f"})
}
