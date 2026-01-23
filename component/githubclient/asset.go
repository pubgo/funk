package githubclient

import (
	"strings"
	"time"

	"github.com/docker/go-units"
	"github.com/google/go-github/v71/github"
)

func GetAssetList(repositoryReleases []*github.RepositoryRelease) Assets {
	var assetList = make(Assets, 0, len(repositoryReleases))
	for _, a := range repositoryReleases {
		assetList = append(assetList, GetAssets(a)...)
	}
	return assetList
}

func GetAssets(repositoryRelease *github.RepositoryRelease) Assets {
	var assetList = make(Assets, 0, len(repositoryRelease.Assets))
	for _, a := range repositoryRelease.Assets {
		assetList = append(assetList, Asset{
			Name:      repositoryRelease.GetTagName(),
			Filename:  a.GetName(),
			URL:       a.GetBrowserDownloadURL(),
			Type:      a.GetContentType(),
			Size:      a.GetSize(),
			CreatedAt: a.GetCreatedAt().Time,
			OS:        getOS(a.GetName()),
			Arch:      getArch(a.GetName()),

			// maximum file size 64KB
			ChecksumFile: checksumRe.MatchString(strings.ToLower(a.GetName())) && a.GetSize() < 64*1024,
		})
	}
	return assetList
}

type Asset struct {
	Name, Filename, OS, Arch, URL, Type string
	Size                                int
	CreatedAt                           time.Time
	ChecksumFile                        bool
}

func (a Asset) IsChecksumFile() bool {
	return a.ChecksumFile
}

func (a Asset) Key() string {
	return a.OS + "/" + a.Arch
}

func (a Asset) Is32Bit() bool {
	return a.Arch == "386"
}

func (a Asset) IsMac() bool {
	return a.OS == "darwin"
}

func (a Asset) IsWindows() bool {
	return a.OS == "windows"
}

func (a Asset) IsLinux() bool {
	return a.OS == "linux"
}

func (a Asset) IsMacM1() bool {
	return a.IsMac() && a.Arch == "arm64"
}

// IsArchive 根据文件扩展名判断是否为归档文件
func (a Asset) IsArchive() bool {
	filename := strings.ToLower(a.Filename)
	return strings.HasSuffix(filename, ".zip") ||
		strings.HasSuffix(filename, ".tar.gz") ||
		strings.HasSuffix(filename, ".tgz") ||
		strings.HasSuffix(filename, ".tar.bz2") ||
		strings.HasSuffix(filename, ".bz2") ||
		strings.HasSuffix(filename, ".gz") ||
		strings.HasSuffix(filename, ".tar")
}

type Assets []Asset

func (as Assets) HasM1() bool {
	//detect if we have a native m1 asset
	for _, a := range as {
		if a.IsMacM1() {
			return true
		}
	}
	return false
}

func GetSizeFormat(size int) string {
	return units.HumanSize(float64(size))
}
