package spdx22json

import (
	"fmt"
	"io"

	"github.com/spdx/tools-golang/jsonloader"

	"github.com/anchore/syft/internal/formats/common/spdxhelpers"
	"github.com/anchore/syft/syft/sbom"
)

func decoder(reader io.Reader) (*sbom.SBOM, error) {
	doc, err := jsonloader.Load2_2(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to decode spdx-json: %w", err)
	}

	return spdxhelpers.ToSyftModel(doc)
}