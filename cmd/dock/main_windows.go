package main

import (
 "fmt"
 "os"
 "os/exec"
 "path/filepath"
 "strings"
 "sync/atomic"

 "github.com/amdkna/tool-project-dock/internal/config"
 "github.com/lxn/walk"
 . "github.com/lxn/walk/declarative"
)

const version = "0.1.3"

func main() {
 exe,err:=os.Executable()
 if err!=nil {walk.MsgBox(nil,"Dock",err.Error(),walk.MsgBoxIconError);return}
 configPath:=filepath.Join(filepath.Dir(exe),"config.toml")
 if err:=config.Ensure(configPath);err!=nil {
  walk.MsgBox(nil,"Dock","Cannot initialize config: "+err.Error(),walk.MsgBoxIconError);return
 }

 var mw *walk.MainWindow
 var output *walk.TextEdit
 var repoList *walk.Composite
 var busy atomic.Bool

 logLine:=func(s string) {
  if mw==nil || output==nil {return}
  mw.Synchronize(func(){
   previous:=output.Text()
   if len(previous)>60000 {previous=previous[len(previous)-45000:]}
   output.SetText(previous+s+"\r\n")
  })
 }

 runGit:=func(repo config.Repository,action string) {
  if !busy.CompareAndSwap(false,true) {
   walk.MsgBox(mw,"Dock","Another Git operation is running.",walk.MsgBoxIconInformation);return
  }
  go func(){
   defer busy.Store(false)
   info,err:=os.Stat(repo.Directory)
   if err!=nil || !info.IsDir() {
    logLine(fmt.Sprintf("[%s] Missing local clone: %s (clone first)",repo.Name,repo.Directory));return
   }
   cmd:=exec.Command("git","-C",repo.Directory,"remote","get-url",repo.Remote)
   remoteURL,err:=cmd.Output()
   if err!=nil {logLine(fmt.Sprintf("[%s] Cannot read Git remote: %v",repo.Name,err));return}
   if strings.TrimSpace(string(remoteURL))!=repo.URL {
    logLine(fmt.Sprintf("[%s] Local remote URL differs from config URL. Cancelled.",repo.Name));return
   }
   args:=[]string{"-C",repo.Directory}
   if action=="PULL" {
    args=append(args,"pull","--ff-only",repo.Remote,repo.Branch)
   } else {
    args=append(args,"push",repo.Remote,repo.Branch)
   }
   logLine("["+repo.Name+"] "+action+" started...")
   git:=exec.Command("git",args...)
   git.Env=append(os.Environ(),"GIT_TERMINAL_PROMPT=0")
   out,err:=git.CombinedOutput()
   if len(out)>0 {logLine(strings.TrimRight(string(out),"\r\n"))}
   if err!=nil {logLine(fmt.Sprintf("[ERROR] %s %s: %v",repo.Name,action,err));return}
   logLine(fmt.Sprintf("[OK] %s %s completed",repo.Name,action))
  }()
 }

 openConfig:=func(){
  if err:=exec.Command("notepad.exe",configPath).Start();err!=nil {
   walk.MsgBox(mw,"Dock","Cannot open config: "+err.Error(),walk.MsgBoxIconError)
  }
 }

 var refresh func()
 refresh=func(){
  if busy.Load() {walk.MsgBox(mw,"Dock","Finish Git operation before refresh.",walk.MsgBoxIconInformation);return}
  cfg,err:=config.Load(configPath)
  if err!=nil {
   walk.MsgBox(mw,"Config error",err.Error(),walk.MsgBoxIconError)
   logLine("[ERROR] Refresh: "+err.Error());return
  }
  for repoList.Children().Len()>0 {repoList.Children().At(0).Dispose()}
  for _,repository:=range cfg.Repositories {
   repo:=repository
   box,err:=walk.NewGroupBox(repoList)
   if err!=nil {logLine("[ERROR] Project box: "+err.Error());return}
   box.SetLayout(walk.NewHBoxLayout())
   pull,err:=walk.NewPushButton(box)
   if err!=nil {logLine("[ERROR] PULL button: "+err.Error());return}
   pull.SetText("PULL")
   pull.Clicked().Attach(func(){runGit(repo,"PULL")})
   push,err:=walk.NewPushButton(box)
   if err!=nil {logLine("[ERROR] PUSH button: "+err.Error());return}
   push.SetText("PUSH")
   push.Clicked().Attach(func(){runGit(repo,"PUSH")})
   name,err:=walk.NewLabel(box)
   if err!=nil {logLine("[ERROR] Label: "+err.Error());return}
   name.SetText(repo.Name+"   |   "+repo.Branch)
  }
  repoList.RequestLayout()
  logLine(fmt.Sprintf("[OK] Loaded %d repositories",len(cfg.Repositories)))
 }

 ui:=MainWindow{
  AssignTo:&mw,
  Title:"Dock by Webigram | ver. "+version,
  MinSize:Size{Width:630,Height:420},
  Size:Size{Width:850,Height:600},
  Layout:VBox{MarginsZero:false,Spacing:9},
  MenuItems:[]MenuItem{
   Menu{Text:"&File",Items:[]MenuItem{
    Action{Text:"&Config...",OnTriggered:openConfig},
    Action{Text:"&Refresh",OnTriggered:func(){refresh()}},
   }},
   Menu{Text:"&Info",Items:[]MenuItem{
    Action{Text:"ver. "+version,Enabled:false},
   }},
  },
  Children:[]Widget{
   Label{Text:"Configured repositories  |  File > Config to edit  |  File > Refresh to reload"},
   Composite{AssignTo:&repoList,Layout:VBox{MarginsZero:true,Spacing:5}},
   Label{Text:"Activity log"},
   TextEdit{AssignTo:&output,ReadOnly:true,VScroll:true},
  },
 }
 if err=ui.Create();err!=nil {walk.MsgBox(nil,"Dock",err.Error(),walk.MsgBoxIconError);return}
 refresh()
 mw.Run()
}
