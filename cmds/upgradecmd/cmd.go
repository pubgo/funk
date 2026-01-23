package upgradecmd

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/go-version"
	"github.com/olekukonko/tablewriter"
	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/log"
	"github.com/pubgo/funk/v2/pretty"
	"github.com/pubgo/funk/v2/result"
	"github.com/pubgo/redant"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/yarlson/tap"

	"github.com/pubgo/funk/v2/cmds/upgradecmd/githubclient"
)

func New(owner, repo string) *redant.Command {
	return &redant.Command{
		Use:   "upgrade",
		Short: "self upgrade management",
		Children: []*redant.Command{
			{
				Use: "list",
				Handler: func(ctx context.Context, i *redant.Invocation) error {
					client := githubclient.NewPublicRelease(owner, repo)
					releases := assert.Must1(client.List(ctx))

					tt := tablewriter.NewWriter(os.Stdout)
					tt.SetHeader([]string{"Name", "Size", "Url"})

					for _, r := range releases {
						for _, a := range githubclient.GetAssets(r) {
							if a.IsChecksumFile() {
								continue
							}

							if a.OS != runtime.GOOS {
								continue
							}

							if a.Arch != runtime.GOARCH {
								continue
							}

							tt.Append([]string{
								a.Name,
								githubclient.GetSizeFormat(a.Size),
								a.URL,
							})
						}
					}
					tt.Render()
					return nil
				},
			},
		},
		Handler: func(ctx context.Context, i *redant.Invocation) (gErr error) {
			defer result.RecoveryErr(&gErr, func(err error) error {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				pretty.Println(err)
				return err
			})

			client := githubclient.NewPublicRelease(owner, repo)
			r := assert.Must1(client.List(ctx))

			assets := githubclient.GetAssetList(r)
			assets = lo.Filter(assets, func(item githubclient.Asset, index int) bool {
				return !item.IsChecksumFile() && item.OS == runtime.GOOS && item.Arch == runtime.GOARCH
			})
			sort.Slice(assets, func(i, j int) bool {
				v1, err1 := version.NewSemver(assets[i].Name)
				v2, err2 := version.NewSemver(assets[j].Name)
				if err1 != nil || err2 != nil {
					// fallback to string comparison if version parsing fails
					return assets[i].Name > assets[j].Name
				}
				return v1.GreaterThan(v2)
			})

			if len(assets) > 20 {
				assets = assets[:20]
			}

			versionName := tap.Select[string](ctx, tap.SelectOptions[string]{
				Message: "Which version do you prefer?",
				Options: lo.Map(assets, func(item githubclient.Asset, index int) tap.SelectOption[string] {
					return tap.SelectOption[string]{
						Value: item.Name,
						Label: item.Name,
					}
				}),
			})

			if versionName == "" {
				return nil
			}

			log.Info(ctx).Msgf("You chose: %s", versionName)

			asset, ok := lo.Find(assets, func(item githubclient.Asset) bool { return item.Name == versionName })
			assert.If(!ok, "%s not found", versionName)
			var downloadURL = asset.URL

			downloadDir := filepath.Join(os.TempDir(), repo)
			pwd := assert.Must1(os.Getwd())

			// 获取当前可执行文件的完整路径
			execFile, err := os.Executable()
			if err != nil {
				// 如果 os.Executable() 失败，回退到使用 os.Args[0]
				execFile = os.Args[0]
				if !filepath.IsAbs(execFile) {
					// 如果不是绝对路径，尝试在 PATH 中查找
					if found, err := exec.LookPath(execFile); err == nil {
						execFile = found
					}
				}
			}
			// 解析符号链接，获取真实路径
			execFile = assert.Must1(filepath.EvalSymlinks(execFile))

			log.Info().Func(func(e *zerolog.Event) {
				e.Str("download_dir", downloadDir)
				e.Str("pwd", pwd)
				e.Str("exec_file", execFile)
				e.Msgf("start download %s", downloadURL)
			})

			c := &getter.Client{
				Ctx:              ctx,
				Src:              downloadURL,
				Dst:              downloadDir,
				Pwd:              pwd,
				Mode:             getter.ClientModeDir,
				ProgressListener: defaultProgressBar,
			}
			assert.Must(c.Get())
			assert.Must(os.Rename(filepath.Join(downloadDir, asset.Filename), execFile))

			return nil
		},
	}
}
