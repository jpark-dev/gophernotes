package main

import (

	"log"
    "os"
    "sort"

	interp "github.com/cosmos72/gomacro/fast"

    "github.com/gopherdata/gophernotes/complete"
)

type Completion struct {
	class,
	name,
	typ string
}

type CompletionResponse struct {
	partial     int
	completions []Completion
}

/************************************************************
* entry function
************************************************************/
func handleCompleteRequest(ir *interp.Interp, receipt msgReceipt) error {
	// Extract the data from the request.
	reqcontent := receipt.Msg.Content.(map[string]interface{})
    code := reqcontent["code"].(string)
    cursorPos := reqcontent["cursor_pos"].(float64)

	// Tell the front-end that the kernel is working and when finished notify the
	// front-end that the kernel is idle again.
	if err := receipt.PublishKernelStatus(kernelBusy); err != nil {
		log.Printf("Error publishing kernel status 'busy': %v\n", err)
	}

	// Redirect the standard out from the REPL.
	oldStdout := os.Stdout
	_, wOut, err := os.Pipe()
	if err != nil {
		return err
	}
	os.Stdout = wOut

	// Redirect the standard error from the REPL.
	oldStderr := os.Stderr
	_, wErr, err := os.Pipe()
	if err != nil {
		return err
	}
	os.Stderr = wErr

    // evaluate the code to the cursor position
    matches, prefix, err := complete.FindMatch(ir, code, int(cursorPos))
    sort.Strings(matches)

	// Close and restore the streams.
	wOut.Close()
	os.Stdout = oldStdout

	wErr.Close()
	os.Stderr = oldStderr

    // prepare the reply
    content := make(map[string]interface{})

    if err == nil {
        content["cursor_start"] = cursorPos - float64(len(prefix))
        content["cursor_end"] = cursorPos
        content["matches"] = matches
        content["status"] = "ok"
    } else {
        content["ename"] = "ERROR"
        content["evalue"] = err.Error()
        content["traceback"] = nil
        content["status"] = "error"
    }

    return receipt.Reply("complete_reply", content)
}

