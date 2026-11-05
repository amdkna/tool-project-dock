package config
import (
 "errors"
 "fmt"
 "os"
 "path/filepath"
 "regexp"
 "strconv"
 "strings"
 "github.com/BurntSushi/toml"
)
type Config struct {Repositories []Repository `toml:"repositories"`}
type Repository struct {
 Name string `toml:"name"`
 URL string `toml:"url"`
 Directory string `toml:"directory"`
 Branch string `toml:"branch"`
 Remote string `toml:"remote"`
}
const DefaultTOML = `# Dock by Webigram
# Single-backslash Windows paths work in double-quoted directory values.
[[repositories]]
name = "Sample Project"
url = "https://github.com/example/project.git"
directory = "C:\devs\sample-project\"
branch = "main"
remote = "origin"
`

var windowsDirectory = regexp.MustCompile(`(?m)^([ \t]*directory[ \t]*=[ \t]*)"([A-Za-z]:\\[^"\r\n]*)"([ \t]*(?:#.*)?)(\r?)$`)

func normalizeWindowsDirectories(data string) string {
 return windowsDirectory.ReplaceAllStringFunc(data,func(line string) string{
  m:=windowsDirectory.FindStringSubmatch(line)
  if len(m)!=5{return line}
  path:=m[2]
  raw:=false
  for i:=0;i<len(path);{
   if path[i]!='\\'{i++;continue}
   j:=i
   for j<len(path)&&path[j]=='\\'{j++}
   if (j-i)%2!=0{raw=true;break}
   i=j
  }
  if !raw{return line}
  value:="'"+path+"'"
  if strings.Contains(path,"'"){value=strconv.Quote(path)}
  return m[1]+value+m[3]+m[4]
 })
}
func Ensure(path string) error{
 _,err:=os.Stat(path)
 if err==nil{return nil}
 if !errors.Is(err,os.ErrNotExist){return err}
 f,err:=os.OpenFile(path,os.O_WRONLY|os.O_CREATE|os.O_EXCL,0600)
 if errors.Is(err,os.ErrExist){return nil}
 if err!=nil{return err}
 defer f.Close()
 _,err=f.WriteString(DefaultTOML)
 return err
}
func Load(path string)(Config,error){
 var cfg Config
 content,err:=os.ReadFile(path)
 if err!=nil{return cfg,fmt.Errorf("read config: %w",err)}
 meta,err:=toml.Decode(normalizeWindowsDirectories(string(content)),&cfg)
 if err!=nil{return cfg,fmt.Errorf("parse config: %w",err)}
 if u:=meta.Undecoded();len(u)>0{return cfg,fmt.Errorf("unknown configuration field: %s",u[0].String())}
 if len(cfg.Repositories)==0{return cfg,errors.New("add at least one [[repositories]] entry")}
 used:=map[string]bool{}
 for i:=range cfg.Repositories{
  r:=&cfg.Repositories[i]
  r.Name=strings.TrimSpace(r.Name)
  r.URL=strings.TrimSpace(r.URL)
  r.Directory=strings.TrimSpace(r.Directory)
  r.Branch=strings.TrimSpace(r.Branch)
  r.Remote=strings.TrimSpace(r.Remote)
  if r.Name==""||r.URL==""||r.Directory==""||r.Branch==""{return cfg,fmt.Errorf("repositories[%d]: name, url, directory and branch are required",i+1)}
  if r.Remote==""{r.Remote="origin"}
  if !filepath.IsAbs(r.Directory){r.Directory=filepath.Join(filepath.Dir(path),r.Directory)}
  r.Directory=filepath.Clean(r.Directory)
  key:=strings.ToLower(r.Directory)
  if used[key]{return cfg,fmt.Errorf("duplicate local directory: %s",r.Directory)}
  used[key]=true
 }
 return cfg,nil
}
