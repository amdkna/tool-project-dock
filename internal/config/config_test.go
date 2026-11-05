package config
import (
 "testing"
 "github.com/BurntSushi/toml"
)
func TestWindowsDirectorySyntax(t *testing.T){
 cases:=[]struct{name,input,want string}{
  {"raw",`directory = "C:\devs\test\"`,`C:\devs\test\`},
  {"escape-like",`directory = "C:\test"`,`C:\test`},
  {"proper escaped",`directory = "C:\\devs\\test\\"`,`C:\devs\test\`},
  {"literal",`directory = 'C:\devs\test\'`,`C:\devs\test\`},
  {"apostrophe",`directory = "C:\John's project\"`,`C:\John's project\`},
  {"comment",`directory = "C:\devs\test\" # note`,`C:\devs\test\`},
 }
 for _,tc:=range cases{
  t.Run(tc.name,func(t *testing.T){
   var cfg struct{Directory string `toml:"directory"`}
   _,err:=toml.Decode(normalizeWindowsDirectories(tc.input),&cfg)
   if err!=nil{t.Fatal(err)}
   if cfg.Directory!=tc.want{t.Fatalf("got %q, want %q",cfg.Directory,tc.want)}
  })
 }
}
