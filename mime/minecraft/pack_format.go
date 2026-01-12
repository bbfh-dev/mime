package minecraft

import (
	"fmt"
	"strings"
)

const (
	USES_SUPPORTED_FORMATS uint8 = 1 + iota
	USES_MIN_MAX_FORMAT
)

type VersionRange struct {
	Min Version
	Max Version
}

type Version struct {
	Digits [2]int
	Flag   uint8
}

func (version Version) Value() any {
	if version.Digits[1] == 0 {
		return version.Digits[0]
	}

	if version.Digits[1] > 9 {
		panic(
			fmt.Sprintf(
				"(Assertion fail) Version [%d, %d] is not yet supported internally",
				version.Digits[0],
				version.Digits[1],
			),
		)
	}
	return float64(version.Digits[0]) + float64(version.Digits[1])/10
}

var DataPackFormats = map[string]Version{
	"1.13":    {Digits: [2]int{4, 0}},
	"1.13.1":  {Digits: [2]int{4, 0}},
	"1.13.2":  {Digits: [2]int{4, 0}},
	"1.14":    {Digits: [2]int{4, 0}},
	"1.14.1":  {Digits: [2]int{4, 0}},
	"1.14.2":  {Digits: [2]int{4, 0}},
	"1.14.3":  {Digits: [2]int{4, 0}},
	"1.14.4":  {Digits: [2]int{4, 0}},
	"1.15":    {Digits: [2]int{5, 0}},
	"1.15.1":  {Digits: [2]int{5, 0}},
	"1.15.2":  {Digits: [2]int{5, 0}},
	"1.16":    {Digits: [2]int{5, 0}},
	"1.16.1":  {Digits: [2]int{5, 0}},
	"1.16.2":  {Digits: [2]int{6, 0}},
	"1.16.3":  {Digits: [2]int{6, 0}},
	"1.16.4":  {Digits: [2]int{6, 0}},
	"1.16.5":  {Digits: [2]int{6, 0}},
	"1.17":    {Digits: [2]int{7, 0}},
	"1.17.1":  {Digits: [2]int{7, 0}},
	"1.18":    {Digits: [2]int{8, 0}},
	"1.18.1":  {Digits: [2]int{8, 0}},
	"1.18.2":  {Digits: [2]int{9, 0}},
	"1.19":    {Digits: [2]int{10, 0}},
	"1.19.1":  {Digits: [2]int{10, 0}},
	"1.19.2":  {Digits: [2]int{10, 0}},
	"1.19.3":  {Digits: [2]int{10, 0}},
	"1.19.4":  {Digits: [2]int{12, 0}},
	"1.20":    {Digits: [2]int{15, 0}},
	"1.20.1":  {Digits: [2]int{15, 0}},
	"1.20.2":  {Digits: [2]int{18, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.20.3":  {Digits: [2]int{26, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.20.4":  {Digits: [2]int{26, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.20.5":  {Digits: [2]int{41, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.20.6":  {Digits: [2]int{41, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21":    {Digits: [2]int{48, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21.1":  {Digits: [2]int{48, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21.2":  {Digits: [2]int{57, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21.3":  {Digits: [2]int{57, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21.4":  {Digits: [2]int{61, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21.5":  {Digits: [2]int{71, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21.6":  {Digits: [2]int{80, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21.7":  {Digits: [2]int{81, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21.8":  {Digits: [2]int{81, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21.9":  {Digits: [2]int{88, 0}, Flag: USES_MIN_MAX_FORMAT},
	"1.21.10": {Digits: [2]int{88, 0}, Flag: USES_MIN_MAX_FORMAT},
	"1.21.11": {Digits: [2]int{94, 1}, Flag: USES_MIN_MAX_FORMAT},
}

var ResourcePackFormats = map[string]Version{
	"1.13":    {Digits: [2]int{4, 0}},
	"1.13.1":  {Digits: [2]int{4, 0}},
	"1.13.2":  {Digits: [2]int{4, 0}},
	"1.14":    {Digits: [2]int{4, 0}},
	"1.14.1":  {Digits: [2]int{4, 0}},
	"1.14.2":  {Digits: [2]int{4, 0}},
	"1.14.3":  {Digits: [2]int{4, 0}},
	"1.14.4":  {Digits: [2]int{4, 0}},
	"1.15":    {Digits: [2]int{5, 0}},
	"1.15.1":  {Digits: [2]int{5, 0}},
	"1.15.2":  {Digits: [2]int{5, 0}},
	"1.16":    {Digits: [2]int{5, 0}},
	"1.16.1":  {Digits: [2]int{5, 0}},
	"1.16.2":  {Digits: [2]int{6, 0}},
	"1.16.3":  {Digits: [2]int{6, 0}},
	"1.16.4":  {Digits: [2]int{6, 0}},
	"1.16.5":  {Digits: [2]int{6, 0}},
	"1.17":    {Digits: [2]int{7, 0}},
	"1.17.1":  {Digits: [2]int{7, 0}},
	"1.18":    {Digits: [2]int{8, 0}},
	"1.18.1":  {Digits: [2]int{8, 0}},
	"1.18.2":  {Digits: [2]int{8, 0}},
	"1.19":    {Digits: [2]int{9, 0}},
	"1.19.1":  {Digits: [2]int{9, 0}},
	"1.19.2":  {Digits: [2]int{9, 0}},
	"1.19.3":  {Digits: [2]int{12, 0}},
	"1.19.4":  {Digits: [2]int{13, 0}},
	"1.20":    {Digits: [2]int{15, 0}},
	"1.20.1":  {Digits: [2]int{15, 0}},
	"1.20.2":  {Digits: [2]int{18, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.20.3":  {Digits: [2]int{22, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.20.4":  {Digits: [2]int{22, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.20.5":  {Digits: [2]int{32, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.20.6":  {Digits: [2]int{32, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21":    {Digits: [2]int{34, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21.1":  {Digits: [2]int{34, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21.2":  {Digits: [2]int{42, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21.3":  {Digits: [2]int{42, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21.4":  {Digits: [2]int{46, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21.5":  {Digits: [2]int{55, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21.6":  {Digits: [2]int{63, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21.7":  {Digits: [2]int{64, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21.8":  {Digits: [2]int{64, 0}, Flag: USES_SUPPORTED_FORMATS},
	"1.21.9":  {Digits: [2]int{69, 0}, Flag: USES_MIN_MAX_FORMAT},
	"1.21.10": {Digits: [2]int{69, 0}, Flag: USES_MIN_MAX_FORMAT},
	"1.21.11": {Digits: [2]int{75, 0}, Flag: USES_MIN_MAX_FORMAT},
}

func IsVersionSupported(version string) bool {
	for version_fragment := range strings.SplitSeq(version, "-") {
		if _, ok := DataPackFormats[version_fragment]; !ok {
			return false
		}
	}
	return true
}
