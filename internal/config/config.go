package config

import (
 "errors"
 "fmt"
 "os"
 "path/filepath"
 "strings"

 "github.com/BurntSushi/toml"
)

type Config struct {
 Project Project `toml:"project"`
 Git Git `toml:"git"`
}
type Project struct {
 Name string `toml:"name"`
 Directory string `toml:"directory"`
}
type Git struct {
 Remote string `toml:"remote"`
 Branch string `toml:"branch"`
}

const DefaultTOML = `# Dock by Webigram - local configuration
# Relative paths resolve relative to config.toml
[project]
name = "My Project"
directory = "./workspace"

[git]
remote = "origin"
branch = "main"
`

func Ensure(path string) error {
 _, err := os.Stat(path)
 if err == nil { return nil }
 if !errors.Is(err, os.ErrNotExist) { return err }
 f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
 if errors.Is(err, os.ErrExist) { return nil }
 if err != nil { return err }
 defer f.Close()
 _, err = f.WriteString(DefaultTOML)
 return err
}

func Load(path string) (Config, error) {
 var cfg Config
 if _, err := toml.DecodeFile(path, &cfg); err != nil { return cfg, fmt.Errorf("read config: %w", err) }
 if strings.TrimSpace(cfg.Project.Directory) == "" { return cfg, errors.New("project.directory is required") }
 if strings.TrimSpace(cfg.Git.Remote) == "" { return cfg, errors.New("git.remote is required") }
 if strings.TrimSpace(cfg.Git.Branch) == "" { return cfg, errors.New("git.branch is required") }
 if !filepath.IsAbs(cfg.Project.Directory) { cfg.Project.Directory = filepath.Join(filepath.Dir(path), cfg.Project.Directory) }
 cfg.Project.Directory = filepath.Clean(cfg.Project.Directory)
 return cfg, nil
}
