package console

const basePrompt = ">>>"

const (
	updateIndexDBCmd = "update-index"
	showStatsCmd     = "show-stats"
	showDupsCmd      = "show-dups"
	runWebServerCmd  = "run-web-server"
	exitCmd          = "exit"
)

func getCommandList() []string {
	return []string{
		updateIndexDBCmd,
		showDupsCmd,
		runWebServerCmd,
		showStatsCmd,
		exitCmd,
	}
}
