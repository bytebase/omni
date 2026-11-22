package driver

import (
	"github.com/bytebase/omni/driver"
	hiveanalysis "github.com/bytebase/omni/hive/analysis"
)

// splitCommands cuts script into the statements a request runs, in order, with
// the text and byte range a task log reports.
func splitCommands(script string) ([]driver.Command, error) {
	statements, err := hiveanalysis.Split(script)
	if err != nil {
		return nil, err
	}
	commands := make([]driver.Command, 0, len(statements))
	for _, stmt := range statements {
		if stmt.Empty {
			continue
		}
		commands = append(commands, driver.Command{Text: stmt.Text, ByteStart: stmt.Start, ByteEnd: stmt.End})
	}
	return commands, nil
}
