// Minecraft-related hardcoded values and utility functions
package minecraft

import (
	"fmt"
	"strings"

	liberrors "github.com/bbfh-dev/lib-errors"
)

func OverlayVersions(folder string) (string, string, error) {
	parts := strings.SplitN(folder, "__", 2)
	if len(parts) != 2 {
		return "", "", &liberrors.DetailedError{
			Label:   liberrors.ERR_FORMAT,
			Context: nil,
			Details: fmt.Sprintf(
				"overlay folder must be `[mc_version]__[mc_version]`, but got: %q",
				folder,
			),
		}
	}

	if !IsVersionSupported(parts[0]) {
		return "", "", &liberrors.DetailedError{
			Label:   liberrors.ERR_VALIDATE,
			Context: nil,
			Details: fmt.Sprintf(
				"unknown version %q, could be an invalid format or outdated vintage version",
				parts[0],
			),
		}
	}

	if !IsVersionSupported(parts[1]) {
		return "", "", &liberrors.DetailedError{
			Label:   liberrors.ERR_VALIDATE,
			Context: nil,
			Details: fmt.Sprintf(
				"unknown version %q, could be an invalid format or outdated vintage version",
				parts[1],
			),
		}
	}

	return parts[0], parts[1], nil
}
