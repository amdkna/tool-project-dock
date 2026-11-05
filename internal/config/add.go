package config

import (
 "fmt"
 "net/url"
 "os"
 "path/filepath"
 "strconv"
 "strings"
)

// GitHubName accepts HTTPS repository URLs without embedded credentials.
func GitHubName(raw string)(string,error){
 parsed,err:=url.Parse(strings.TrimSpace(raw))
 if err!=nil||parsed==nil||parsed.Scheme!="https"||!strings.EqualFold(parsed.Hostname(),"github.com")||parsed.User!=nil||parsed.RawQuery!=""||parsed.Fragment!=""{
  return "",fmt.Errorf("enter an HTTPS GitHub repository URL, e.g. https://github.com/owner/repo.git")
 }
 bits:=strings.Split(strings.Trim(parsed.Path,"/"),"/")
 if len(bits)!=2||bits[0]==""||bits[1]==""{return "",fmt.Errorf("URL must have a GitHub owner and repository")}
 name:=strings.TrimSuffix(bits[1],".git")
 if name==""||name=="."||name==".."||strings.ContainsAny(name,":\\/*?\"<>|"){return "",fmt.Errorf("invalid repository name")}
 return name,nil
}

// Add saves a repository entry, preserving all existing TOML comments and settings.
func Add(path string,r Repository) error{
 cfg,err:=Load(path)
 if err!=nil{return err}
 r.Name=strings.TrimSpace(r.Name)
 r.URL=strings.TrimSpace(r.URL)
 r.Branch=strings.TrimSpace(r.Branch)
 r.Directory=filepath.Clean(strings.TrimSpace(r.Directory))
 if r.Name==""||r.Branch==""||r.Directory=="."{return fmt.Errorf("name, branch and directory are required")}
 if _,err=GitHubName(r.URL);err!=nil{return err}
 if r.Remote==""{r.Remote="origin"}
 for _,existing:=range cfg.Repositories{
  if strings.EqualFold(filepath.Clean(existing.Directory),r.Directory){
   return fmt.Errorf("directory already configured: %s",r.Directory)
  }
 }
 original,err:=os.ReadFile(path)
 if err!=nil{return err}
 directory:="'"+r.Directory+"'"
 if strings.Contains(r.Directory,"'"){directory=strconv.Quote(r.Directory)}
 addition:="\n[[repositories]]\nname = "+strconv.Quote(r.Name)+
  "\nurl = "+strconv.Quote(r.URL)+
  "\ndirectory = "+directory+
  "\nbranch = "+strconv.Quote(r.Branch)+
  "\nremote = "+strconv.Quote(r.Remote)+"\n"
 tmp,err:=os.CreateTemp(filepath.Dir(path),".dock-config-*.tmp")
 if err!=nil{return err}
 defer os.Remove(tmp.Name())
 if _,err=tmp.Write(append(original,[]byte(addition)...));err!=nil{tmp.Close();return err}
 if err=tmp.Close();err!=nil{return err}
 if _,err=Load(tmp.Name());err!=nil{return fmt.Errorf("config validation failed: %w",err)}
 return os.Rename(tmp.Name(),path)
}
