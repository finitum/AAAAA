package makepkg

import "strings"

type SrcInfo struct {
	PackageName string
	PackageDescription string
	PackageVersion string
	PackageRelease string
	MakeDependencies []string
	Dependencies []string

	OtherFields map[string][]string
}

func ParseSrcInfo(srcinfo string) *SrcInfo {
	lines := strings.Split(srcinfo, "\n")

	info := SrcInfo{
		OtherFields: make(map[string][]string),
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.SplitN(line, "=", 2)

		if len(parts) < 2 {
			continue
		}

		switch parts[0] {
		case "pkgname":
			info.PackageName = parts[1]
		case "pkgdesc":
			info.PackageDescription = parts[1]
		case "pkgver":
			info.PackageVersion = parts[1]
		case "pkgrel":
			info.PackageRelease = parts[1]
		case "makedepends":
			info.MakeDependencies = append(info.MakeDependencies,parts[1])
		case "depends":
			info.Dependencies = append(info.Dependencies,parts[1])
		default:
			info.OtherFields[parts[0]] = append(info.OtherFields[parts[0]],parts[1])
		}
	}

	return &info
}
