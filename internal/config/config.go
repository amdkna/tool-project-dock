package config

import (
 "errors"
 "fmt"
 "os"
 "path/filepath"
 "strings"

 "github.com/BurntSushi/toml"
)

type Config struct { Repositories []Repository `toml:"repositories"` }
type Repository struct {
 Name string `toml:"name"`
 URL string `toml:"url"`
 Directory string `toml:"directory"`
 Branch string `toml:"branch"`
 Remote string `toml:"remote"`
}

const DefaultTOML = `# Dock by Webigram
# Use single quotes for Windows paths. Add another [[repositories]] block per app.

[[repositories]]
name = "Sample Project"
url = "https://github.com/example/project.git"
directory = './workspace/sample-project'
branch = "main"
remote = "origin"
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
 meta, err := toml.DecodeFile(path, &cfg)
 if err != nil { return cfg, fmt.Errorf("read config: %w", err) }
 if u := meta.Undecoded(); len(u)>0 { return cfg, fmt.Errorf("unknown configuration field: %s",u[0].String()) }
 if len(cfg.Repositories)==0 { return cfg, errors.New("add at least one [[repositories]] entry") }
 used := map[string]bool{}
 for i:=range cfg.Repositories {
  r:=&cfg.Repositories[i]
  r.Name = strings.TrimSpace(r.Name)
  r.URL = strings.TrimSpace(r.URL)
  r.Directory = strings.TrimSpace(r.Directory)
  r.Branch = strings.TrimSpace(r.Branch)
  r.Remote = strings.TrimSpace(r.Remote)
  if r.Name=="" || r.URL=="" || r.Directory=="" || r.Branch=="" {
   return cfg, fmt.Errorf("repositories[%d]: name, url, directory and branch are required",i+1)
  }
  if r.Remote=="" { r.Remote="origin" }
  if !filepath.IsAbs(r.Directory) { r.Directory=filepath.Join(filepath.Dir(path),r.Directory) }
  r.Directory=filepath.Clean(r.Directory)
  key:=strings.ToLower(r.Directory)
  if used[key] { return cfg, fmt.Errorf("duplicate local directory: %s",r.Directory) }
  used[key]=true
 }
 return cfg,nil
}
