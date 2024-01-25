package runner

import "github.com/projectdiscovery/gologger"

const banner = `
          __    __                  __            
   ____  / /_  / /__________ ______/ /_____  _____
  / __ \/ __ \/ __/ ___/ __ \/ ___/ //_/ _ \/ ___/
 / /_/ / / / / /_/ /  / /_/ / /__/ ,< /  __/ /
 \__, /_/ /_/\__/_/   \__,_/\___/_/|_|\___/_/
/____/
`

// Name
const ToolName = `ghtracker`

// showBanner is used to show the banner to the user
func showBanner() {
	gologger.Print().Msgf("%s\n", banner)
	gologger.Print().Msgf("\t\tVersion v1.0.0\n\n")
}
