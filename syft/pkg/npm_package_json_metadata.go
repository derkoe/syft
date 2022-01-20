package pkg

import "github.com/anchore/packageurl-go"

// NpmPackageJSONMetadata holds extra information that is used in pkg.Package
type NpmPackageJSONMetadata struct {
	Name        string   `mapstructure:"name" json:"name"`
	Version     string   `mapstructure:"version" json:"version"`
	Files       []string `mapstructure:"files" json:"files,omitempty"`
	Author      string   `mapstructure:"author" json:"author"`
	Licenses    []string `mapstructure:"licenses" json:"licenses"`
	Homepage    string   `mapstructure:"homepage" json:"homepage"`
	Description string   `mapstructure:"description" json:"description"`
	URL         string   `mapstructure:"url" json:"url"`
}

// PackageURL returns the PURL for the specific NPM package (see https://github.com/package-url/purl-spec)
func (p NpmPackageJSONMetadata) PackageURL() string {
	pURL := packageurl.NewPackageURL(
		packageurl.TypeNPM,
		"",
		p.Name,
		p.Version,
		nil,
		"")

	return pURL.ToString()

}