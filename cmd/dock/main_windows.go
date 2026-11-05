package main

import (
 "fmt"
 "os"
 "os/exec"
 "path/filepath"
 "strings"
 "sync/atomic"

 "github.com/amdkna/tool-project-dock/internal/config"
 "github.com/amdkna/tool-project-dock/internal/auth"
 "github.com/lxn/walk"
 . "github.com/lxn/walk/declarative"
)

const version = "0.1.10"

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
   args:=[]string{"-c","credential.helper=","-c","credential.helper=manager","-C",repo.Directory}
   if action=="PULL" {
    args=append(args,"pull","--ff-only",repo.Remote,repo.Branch)
   } else {
    args=append(args,"push",repo.Remote,repo.Branch)
   }
   logLine("["+repo.Name+"] "+action+" started...")
   git:=exec.Command("git",args...)
   git.Env=append(auth.Env(),"GIT_TERMINAL_PROMPT=0")
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

 openExplorer:=func(repo config.Repository){
  if info,err:=os.Stat(repo.Directory);err!=nil||!info.IsDir(){
   walk.MsgBox(mw,"Dock","Folder does not exist: "+repo.Directory,walk.MsgBoxIconWarning);return
  }
  if err:=exec.Command("explorer.exe",repo.Directory).Start();err!=nil{
   walk.MsgBox(mw,"Dock","Cannot open folder: "+err.Error(),walk.MsgBoxIconError)
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
   explorer,err:=walk.NewPushButton(box)
   if err!=nil{logLine("[ERROR] Browse button: "+err.Error());return}
   explorer.SetText("Browse")
   explorer.Clicked().Attach(func(){openExplorer(repo)})
   name,err:=walk.NewLabel(box)
   if err!=nil {logLine("[ERROR] Label: "+err.Error());return}
   name.SetText(repo.Name+"   |   "+repo.Branch)
  }
  repoList.RequestLayout()
  logLine(fmt.Sprintf("[OK] Loaded %d repositories",len(cfg.Repositories)))
 }

 addNew:=func(){
  if busy.Load(){walk.MsgBox(mw,"Dock","Wait until the Git operation finishes.",walk.MsgBoxIconInformation);return}
  repo,ok:=addRepositoryDialog(mw)
  if !ok{return}
  // Persist and display the new repository before cloning, so failed
  // authentication/network operations never silently drop the new entry.
  if err:=config.Add(configPath,repo);err!=nil{
   walk.MsgBox(mw,"Cannot save project","Could not update config.toml: "+err.Error(),walk.MsgBoxIconError)
   return
  }
  refresh()
  logLine("[OK] Saved "+repo.Name+" in "+configPath)
  if !busy.CompareAndSwap(false,true){return}
  go func(){
   defer busy.Store(false)
   // Folder picker selects a parent directory; clone to parent/repo-name.
   if info,err:=os.Stat(repo.Directory);err==nil{
    if !info.IsDir(){logLine("[ERROR] Target exists and is not a directory: "+repo.Directory);return}
    check:=exec.Command("git","-C",repo.Directory,"rev-parse","--is-inside-work-tree")
    if out,err:=check.Output();err!=nil||strings.TrimSpace(string(out))!="true"{
     logLine("[ERROR] Target folder exists but is not a Git checkout: "+repo.Directory);return
    }
    remote:=exec.Command("git","-C",repo.Directory,"remote","get-url","origin")
    remoteURL,err:=remote.Output()
    if err!=nil||strings.TrimSpace(string(remoteURL))!=repo.URL{
     logLine("[ERROR] Existing clone remote does not match supplied GitHub URL.");return
    }
   }else if !os.IsNotExist(err){
    logLine("[ERROR] Could not inspect destination: "+err.Error());return
   }else{
    logLine("[CLONE] "+repo.Name+" -> "+repo.Directory)
    cmd:=exec.Command("git","-c","credential.helper=","-c","credential.helper=manager","clone","--branch",repo.Branch,"--",repo.URL,repo.Directory)
    cmd.Env=append(auth.Env(),"GIT_TERMINAL_PROMPT=0")
    out,err:=cmd.CombinedOutput()
    if err!=nil{logLine(fmt.Sprintf("[ERROR] Clone failed: %v %s",err,strings.TrimSpace(string(out))));return}
    logLine("[OK] Clone finished: "+repo.Name)
   }
   logLine("[OK] Repository ready: "+repo.Name)
  }()
 }

 saveToken:=func(){
  if busy.Load(){walk.MsgBox(mw,"Dock","Finish the current operation first.",walk.MsgBoxIconInformation);return}
  username,token,ok:=tokenDialog(mw)
  if !ok{return}
  if err:=auth.StoreGitHub(username,token);err!=nil{
   walk.MsgBox(mw,"GitHub token",err.Error(),walk.MsgBoxIconError);return
  }
  walk.MsgBox(mw,"GitHub token","Saved using Git Credential Manager and Windows Credential Manager.",walk.MsgBoxIconInformation)
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
    Action{Text:"&GitHub Token...",OnTriggered:saveToken},
   }},
   Menu{Text:"&Info",Items:[]MenuItem{
    Action{Text:"ver. "+version,Enabled:false},
   }},
  },
  Children:[]Widget{
   Composite{Layout:HBox{},Children:[]Widget{
    Label{Text:"Configured repositories"},
    HSpacer{},
    PushButton{Text:"Add New...",OnClicked:addNew},
   }},
   Composite{AssignTo:&repoList,Layout:VBox{MarginsZero:true,Spacing:5}},
   Label{Text:"Activity log"},
   TextEdit{AssignTo:&output,ReadOnly:true,VScroll:true},
  },
 }
 if err=ui.Create();err!=nil {walk.MsgBox(nil,"Dock",err.Error(),walk.MsgBoxIconError);return}
 refresh()
 mw.Run()
}
