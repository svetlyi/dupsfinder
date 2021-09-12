package config

var (
	// nil value for an option means it is not set
	DBPath                   string = "./dupsfinder_dups.db"
	ProcNum                  int    = 1
	DupsPerPage              int    = 100
	IgnoreFilesLessThanBytes int64  = 1000
	LogFile                  string = "dupsfinder.log"
	WebServerPort            int    = 55786
)
