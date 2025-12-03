package velty_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/velty"
	"testing"
)

type Result struct {
	// Name is the name of the tool or step invoked
	Name string `json:"name"`
	// Args holds the original arguments passed to the tool
	Args map[string]interface{} `json:"args"`
	// Result is the string output from the tool invocation
	Result string `json:"result"`
	//Error tool call error
	Error string `json:"error"`
}

var results = []Result{
	{
		Name: "ImageResize",
		Args: map[string]interface{}{
			"width":  800,
			"height": 600,
			"format": "png",
		},
		Result: "Image resized to 800x600 in png format.",
		Error:  "",
	},
	{
		Name: "TextSummarize",
		Args: map[string]interface{}{
			"text":     "Go is an open-source programming language designed at Google.",
			"language": "en",
		},
		Result: "Go is a Google-designed open-source language.",
		Error:  "",
	},
	{
		Name: "DataExport",
		Args: map[string]interface{}{
			"type":     "csv",
			"filename": "report.csv",
		},
		Result: "",
		Error:  "Failed to write to report.csv: permission denied",
	},
}

func TestIt(t *testing.T) {
	options := []velty.Option{velty.BufferSize(8192)}

	planner := velty.New(options...)

	planner.DefineVariable("Results", []Result{})
	var template = `
#foreach($res in $Results)
- ToolName: $res.Name#if($res.Result)
- Result: $res.Result#end#if($res.Error)
- Error: $res.Error#end
- Args:  $res.Args
#end
`

	exec, newState, err := planner.Compile([]byte(template))

	if !assert.Nil(t, err) {
		return
	}
	state := newState()
	state.SetValue("Results", results)
	err = exec.Exec(state)
	assert.Nil(t, err)
	output := state.Buffer.String()

	fmt.Println(output)
}
