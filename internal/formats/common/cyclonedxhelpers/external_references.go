package cyclonedxhelpers

import (
	"fmt"
	"strings"

	syftFile "github.com/anchore/syft/syft/file"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/anchore/syft/syft/pkg"
)

//nolint:funlen, gocognit
func encodeExternalReferences(p pkg.Package) *[]cyclonedx.ExternalReference {
	var refs []cyclonedx.ExternalReference
	if hasMetadata(p) {
		switch metadata := p.Metadata.(type) {
		case pkg.ApkMetadata:
			if metadata.URL != "" {
				refs = append(refs, cyclonedx.ExternalReference{
					URL:  metadata.URL,
					Type: cyclonedx.ERTypeDistribution,
				})
			}
		case pkg.CargoPackageMetadata:
			if metadata.Source != "" {
				refs = append(refs, cyclonedx.ExternalReference{
					URL:  metadata.Source,
					Type: cyclonedx.ERTypeDistribution,
				})
			}
		case pkg.NpmPackageJSONMetadata:
			if metadata.URL != "" {
				refs = append(refs, cyclonedx.ExternalReference{
					URL:  metadata.URL,
					Type: cyclonedx.ERTypeDistribution,
				})
			}
			if metadata.Homepage != "" {
				refs = append(refs, cyclonedx.ExternalReference{
					URL:  metadata.Homepage,
					Type: cyclonedx.ERTypeWebsite,
				})
			}
		case pkg.GemMetadata:
			if metadata.Homepage != "" {
				refs = append(refs, cyclonedx.ExternalReference{
					URL:  metadata.Homepage,
					Type: cyclonedx.ERTypeWebsite,
				})
			}
		case pkg.JavaMetadata:
			if len(metadata.ArchiveDigests) > 0 {
				for _, digest := range metadata.ArchiveDigests {
					refs = append(refs, cyclonedx.ExternalReference{
						URL:  "",
						Type: cyclonedx.ERTypeBuildMeta,
						Hashes: &[]cyclonedx.Hash{{
							Algorithm: cyclonedx.HashAlgorithm(digest.Algorithm),
							Value:     digest.Value,
						}},
					})
				}
			}
		case pkg.PythonPackageMetadata:
			if metadata.DirectURLOrigin != nil && metadata.DirectURLOrigin.URL != "" {
				ref := cyclonedx.ExternalReference{
					URL:  metadata.DirectURLOrigin.URL,
					Type: cyclonedx.ERTypeVCS,
				}
				if metadata.DirectURLOrigin.CommitID != "" {
					ref.Comment = fmt.Sprintf("commit: %s", metadata.DirectURLOrigin.CommitID)
				}
				refs = append(refs, ref)
			}
		}
	}
	if len(refs) > 0 {
		return &refs
	}
	return nil
}

func decodeExternalReferences(c *cyclonedx.Component, metadata interface{}) {
	if c.ExternalReferences == nil {
		return
	}
	switch meta := metadata.(type) {
	case *pkg.ApkMetadata:
		meta.URL = refURL(c, cyclonedx.ERTypeDistribution)
	case *pkg.CargoPackageMetadata:
		meta.Source = refURL(c, cyclonedx.ERTypeDistribution)
	case *pkg.NpmPackageJSONMetadata:
		meta.URL = refURL(c, cyclonedx.ERTypeDistribution)
		meta.Homepage = refURL(c, cyclonedx.ERTypeWebsite)
	case *pkg.GemMetadata:
		meta.Homepage = refURL(c, cyclonedx.ERTypeWebsite)
	case *pkg.JavaMetadata:
		var digests []syftFile.Digest
		if ref := findExternalRef(c, cyclonedx.ERTypeBuildMeta); ref != nil {
			if ref.Hashes != nil {
				for _, hash := range *ref.Hashes {
					digests = append(digests, syftFile.Digest{
						Algorithm: string(hash.Algorithm),
						Value:     hash.Value,
					})
				}
			}
		}

		meta.ArchiveDigests = digests
	case *pkg.PythonPackageMetadata:
		if meta.DirectURLOrigin == nil {
			meta.DirectURLOrigin = &pkg.PythonDirectURLOriginInfo{}
		}
		meta.DirectURLOrigin.URL = refURL(c, cyclonedx.ERTypeVCS)
		meta.DirectURLOrigin.CommitID = strings.TrimPrefix(refComment(c, cyclonedx.ERTypeVCS), "commit: ")
	}
}

func findExternalRef(c *cyclonedx.Component, typ cyclonedx.ExternalReferenceType) *cyclonedx.ExternalReference {
	if c.ExternalReferences != nil {
		for _, r := range *c.ExternalReferences {
			if r.Type == typ {
				return &r
			}
		}
	}
	return nil
}

func refURL(c *cyclonedx.Component, typ cyclonedx.ExternalReferenceType) string {
	if r := findExternalRef(c, typ); r != nil {
		return r.URL
	}
	return ""
}

func refComment(c *cyclonedx.Component, typ cyclonedx.ExternalReferenceType) string {
	if r := findExternalRef(c, typ); r != nil {
		return r.Comment
	}
	return ""
}
