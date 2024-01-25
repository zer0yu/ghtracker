package options

type GHTopDepOptions struct {
	Version      bool   // Version is a flag that prints the current version.
	URL          string // URL is the URL to fetch repositories from.
	Repositories bool   // Repositories is a flag that indicates if repositories should be fetched.
	Table        bool   // Table is a flag that indicates if the output should be in table format.
	Description  bool   // Description is a flag that indicates if the description of the repository should be fetched.
	Rows         int    // Rows is the number of rows to show.
	MinStar      int    // MinStar is the minimum number of stars to show.
	Search       string // Search is the search code at dependents (repositories/packages).
	Token        string // Token is the GitHub token to use for authenticated requests.
	OutputFile   string // Output is the file to write founds to.
}
