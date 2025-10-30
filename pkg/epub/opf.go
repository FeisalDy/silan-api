package epub

import (
	"encoding/xml"
	"errors"
	"fmt"
)

// FindOpfPath locates the OPF file path from container.xml bytes map
func FindOpfPath(files map[string][]byte) (string, error) {
	data, ok := files["META-INF/container.xml"]
	if !ok {
		return "", errors.New("META-INF/container.xml not found")
	}

	type Rootfile struct {
		FullPath string `xml:"full-path,attr"`
	}
	type Container struct {
		XMLName   xml.Name `xml:"container"`
		Rootfiles struct {
			Rootfile []Rootfile `xml:"rootfile"`
		} `xml:"rootfiles"`
	}

	var container Container
	if err := xml.Unmarshal(data, &container); err != nil {
		return "", fmt.Errorf("invalid container.xml: %w", err)
	}

	if len(container.Rootfiles.Rootfile) == 0 {
		return "", errors.New("no rootfile found in container.xml")
	}

	return container.Rootfiles.Rootfile[0].FullPath, nil
}

// ParseOPF parses OPF XML content into OPFPackage
func ParseOPF(opfContent []byte) (*OPFPackage, error) {
	type Package struct {
		XMLName  xml.Name    `xml:"package"`
		Metadata OPFMetadata `xml:"metadata"`
		Manifest struct {
			Items []OPFManifestItem `xml:"item"`
		} `xml:"manifest"`
		Spine OPFSpine `xml:"spine"`
	}

	var pkg Package
	if err := xml.Unmarshal(opfContent, &pkg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OPF: %w", err)
	}

	return &OPFPackage{
		Metadata: pkg.Metadata,
		Manifest: pkg.Manifest.Items,
		Spine:    pkg.Spine,
	}, nil
}
