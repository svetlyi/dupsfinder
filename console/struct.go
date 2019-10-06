package console

const basePrompt = ">>>"

const (
	findDupsCmd     = "find-dups"
	showStatsCmd    = "show-stats"
	showDupsCmd     = "show-dups"
	showLastLogsCmd = "show-last-logs"
	runWebServerCmd = "run-web-server"
	exitCmd         = "exit"
)

func getCommandList() []string {
	return []string{
		findDupsCmd,
		showDupsCmd,
		showLastLogsCmd,
		runWebServerCmd,
		showStatsCmd,
		exitCmd,
	}
}
