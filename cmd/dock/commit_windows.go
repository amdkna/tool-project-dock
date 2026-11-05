package main

import (
 "fmt"
 "os"
 "os/exec"
 "path/filepath"
 "strings"
 "sync/atomic"

 "github.com/amdkna/tool-project-dock/internal/auth"
 "github.com/amdkna/tool-project-dock/internal/config"
 "github.com/lxn/walk"
 . "github.com/lxn/walk/declarative"
)

// commitMessageDialog is deliberately separate from the PUSH action.
// Closing the dialog never stages, commits or pushes anything.
func commitMessageDialog(owner *walk.MainWindow, repo config.Repository) (string,bool) {
 var dlg *walk.Dialog
 var entry *walk.TextEdit
 var message string
 var confirmed bool
 _,err:=(Dialog{
  AssignTo:&dlg,
  Title:"Commit & Push - "+repo.Name,
  MinSize:Size{Width:480,Height:220},
  Layout:VBox{},
  Children:[]Widget{
   Label{Text:"Commit message"},
   TextEdit{AssignTo:&entry},
   Label{Text:"Stages all changes except Git-ignored files, commits, then pushes."},
   Composite{Layout:HBox{},Children:[]Widget{
    HSpacer{},
    PushButton{Text:"Commit & Push",OnClicked:func(){
     value:=strings.TrimSpace(entry.Text())
     if value==""{
      walk.MsgBox(dlg,"Commit message required","Enter a commit message before continuing.",walk.MsgBoxIconWarning)
      return
     }
     message=value
     confirmed=true
     dlg.Accept()
    }},
    PushButton{Text:"Cancel",OnClicked:func(){dlg.Cancel()}},
   }},
  },
 }).Run(owner)
 if err!=nil{
  walk.MsgBox(owner,"Dock","Cannot open commit dialog: "+err.Error(),walk.MsgBoxIconError)
  return "",false
 }
 return message,confirmed
}

func commitAndPush(owner *walk.MainWindow,repo config.Repository,busy *atomic.Bool,logLine func(string)){
 if !busy.CompareAndSwap(false,true){
  walk.MsgBox(owner,"Dock","Another Git operation is running.",walk.MsgBoxIconInformation)
  return
 }
 message,ok:=commitMessageDialog(owner,repo)
 if !ok{busy.Store(false);return}
 go func(){
  defer busy.Store(false)
  prefix:="["+repo.Name+"] "
  run:=func(args ...string)(string,error){
   command:=exec.Command("git",append([]string{"-C",repo.Directory},args...)...)
   command.Env=auth.Env()
   out,err:=command.CombinedOutput()
   return strings.TrimSpace(string(out)),err
  }
  fail:=func(stage string,err error,details string){
   logLine(prefix+"[ERROR] "+stage+": "+err.Error())
   if details!=""{logLine(details)}
  }
  info,err:=os.Stat(repo.Directory)
  if err!=nil||!info.IsDir(){logLine(prefix+"[ERROR] Local directory does not exist: "+repo.Directory);return}
  root,err:=run("rev-parse","--show-toplevel")
  if err!=nil{fail("Repository check",err,root);return}
  if !strings.EqualFold(filepath.Clean(root),filepath.Clean(repo.Directory)){
   logLine(prefix+"[ERROR] Configured folder is not the Git repository root; cancelled.")
   return
  }
  actualURL,err:=run("remote","get-url",repo.Remote)
  if err!=nil{fail("Git remote check",err,actualURL);return}
  if actualURL!=repo.URL{logLine(prefix+"[ERROR] Git remote does not match config URL; cancelled.");return}
  branch,err:=run("branch","--show-current")
  if err!=nil{fail("Branch check",err,branch);return}
  if branch!=repo.Branch{
   logLine(prefix+"[ERROR] Checked out branch is "+fmt.Sprintf("%q",branch)+", but config expects "+fmt.Sprintf("%q",repo.Branch)+"; cancelled.")
   return
  }

  logLine(prefix+"Staging all non-ignored changes (git add -A)...")
  out,err:=run("add","-A")
  if err!=nil{fail("Stage",err,out);return}
  // Exit 0 = no staged changes, exit 1 = changes, exit >1 = error.
  diff:=exec.Command("git","-C",repo.Directory,"diff","--cached","--quiet","--exit-code")
  diff.Env=auth.Env()
  diffErr:=diff.Run()
  if diffErr==nil{
   logLine(prefix+"No changes to commit. Existing commits can be pushed separately.")
  } else if code,ok:=diffErr.(*exec.ExitError);ok&&code.ExitCode()==1{
   logLine(prefix+"Creating commit...")
   commitOutput,commitErr:=run("commit","-m",message)
   if commitErr!=nil{fail("Commit",commitErr,commitOutput);return}
   if commitOutput!=""{logLine(commitOutput)}
  }else{
   fail("Check staged changes",diffErr,"")
   return
  }

  logLine(prefix+"Pushing "+repo.Branch+" to "+repo.Remote+"...")
  // GCM is allowed to open an interactive authentication window when needed.
  pushOutput,pushErr:=run("-c","credential.helper=","-c","credential.helper=manager","push",repo.Remote,"HEAD:refs/heads/"+repo.Branch)
  if pushErr!=nil{
   fail("Push (local commit is retained)",pushErr,pushOutput)
   return
  }
  if pushOutput!=""{logLine(pushOutput)}
  logLine(prefix+"[OK] Commit and push workflow completed.")
 }()
}
