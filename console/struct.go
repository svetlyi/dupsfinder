package console

const basePrompt = ">>>"

const (
	findDupsCmd     = "find-dups"
	showStatsCmd    = "show-stats"
	showDupsCmd     = "show-dups"
	runWebServerCmd = "run-web-server"
	exitCmd         = "exit"
)

func getCommandList() []string {
	return []string{
		findDupsCmd,
		showDupsCmd,
		runWebServerCmd,
		showStatsCmd,
		exitCmd,
	}
}
