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

const version = "0.1.12"

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
  if !busy.CompareAndSwap(false,true){return}
  go func(){
   defer busy.Store(false)
   logLine("[ADD] "+repo.Name+" -> "+repo.Directory)

   if err:=checkDestination(repo);err!=nil{
    logLine("[ERROR] Destination check failed: "+err.Error())
    return
   }
   // Register only after clone succeeds. Never create broken cards on errors.
   if _,err:=os.Stat(repo.Directory);os.IsNotExist(err){
    cmd:=exec.Command("git","-c","credential.helper=","-c","credential.helper=manager",
     "clone","--branch",repo.Branch,"--",repo.URL,repo.Directory)
    cmd.Env=auth.Env() // GCM may open its browser authentication flow
    logLine("[CLONE] Connecting to GitHub. Complete sign-in if prompted.")
    out,err:=cmd.CombinedOutput()
    if err!=nil{
     logLine(fmt.Sprintf("[ERROR] Clone failed: %v\n%s",err,strings.TrimSpace(string(out))))
     logLine("[INFO] Repository was not added to config.toml.")
     return
    }
    logLine("[OK] Clone finished: "+repo.Name)
   }else if err==nil{
    gitRoot:=exec.Command("git","-C",repo.Directory,"rev-parse","--show-toplevel")
    out,e:=gitRoot.Output()
    if e!=nil{
     // Existing empty folder is a valid clone destination.
     cmd:=exec.Command("git","-c","credential.helper=","-c","credential.helper=manager",
      "clone","--branch",repo.Branch,"--",repo.URL,repo.Directory)
     cmd.Env=auth.Env()
     cloneOutput,cloneErr:=cmd.CombinedOutput()
     if cloneErr!=nil{
      logLine(fmt.Sprintf("[ERROR] Clone failed: %v\n%s",cloneErr,strings.TrimSpace(string(cloneOutput))))
      return
     }
     logLine("[OK] Clone finished: "+repo.Name)
    }else if strings.TrimSpace(string(out))!=""{
     logLine("[OK] Using existing matching clone: "+repo.Name)
    }
   }else{
    logLine("[ERROR] Cannot inspect destination: "+err.Error())
    return
   }
   if err:=config.Add(configPath,repo);err!=nil{
    logLine("[ERROR] Clone is ready, but saving config.toml failed: "+err.Error())
    return
   }
   logLine("[OK] Saved "+repo.Name+" in "+configPath)
   mw.Synchronize(func(){
    busy.Store(false)
    refresh()
   })
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

 installTool:=func(which string){
  // Run in a visible console so the user can review installer prompts/errors.
  script:=filepath.Join(filepath.Dir(exe),"..","scripts","install-tools.ps1")
  // Executable usually lives in dist. For development builds in other
  // locations, check an adjacent scripts directory as well.
  if _,err:=os.Stat(script);err!=nil{
   script=filepath.Join(filepath.Dir(exe),"scripts","install-tools.ps1")
  }
  if _,err:=os.Stat(script);err!=nil{
   walk.MsgBox(mw,"Install Tools","Installer script not found next to the Dock project. Run from the source checkout or use the documented download links.",walk.MsgBoxIconWarning)
   return
  }
  command:=fmt.Sprintf("& '%s' -Tool %s",strings.ReplaceAll(script,"'","''"),which)
  cmd:=exec.Command("powershell.exe","-NoProfile","-NoExit","-Command",command)
  if err:=cmd.Start();err!=nil{
   walk.MsgBox(mw,"Install Tools","Cannot launch PowerShell installer: "+err.Error(),walk.MsgBoxIconError)
  }
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
    Menu{Text:"Install &Tools",Items:[]MenuItem{
     Action{Text:"Install &Required (Git + GCM)",OnTriggered:func(){installTool("Required")}},
     Action{Text:"Install Git for &Windows",OnTriggered:func(){installTool("Git")}},
     Action{Text:"Install Git &Credential Manager",OnTriggered:func(){installTool("GCM")}},
     Action{Text:"Install GitHub CLI (&gh)",OnTriggered:func(){installTool("GH")}},
    }},
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
