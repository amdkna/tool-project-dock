package main

import (
 "fmt"
 "os/exec"
 "path/filepath"
 "strings"

 "github.com/amdkna/tool-project-dock/internal/config"
 "github.com/lxn/walk"
 . "github.com/lxn/walk/declarative"
)

// addRepositoryDialog asks for an HTTPS GitHub URL, branch and parent folder.
func addRepositoryDialog(owner *walk.MainWindow)(config.Repository,bool){
 var dlg *walk.Dialog
 var urlEdit,folderEdit,branchEdit *walk.LineEdit
 var result config.Repository
 var accepted bool
 _,err:=(Dialog{
  AssignTo:&dlg,
  Title:"Add a GitHub repository",
  MinSize:Size{Width:560,Height:250},
  Layout:VBox{},
  Children:[]Widget{
   Label{Text:"GitHub HTTPS URL"},
   LineEdit{AssignTo:&urlEdit,Text:"https://github.com/owner/repository.git"},
   Label{Text:"Local parent folder (Dock creates the repository subfolder)"},
   Composite{Layout:HBox{},Children:[]Widget{
    LineEdit{AssignTo:&folderEdit},
    PushButton{Text:"Browse...",OnClicked:func(){
     chooser:=walk.FileDialog{Title:"Choose the parent folder for the repository"}
     yes,err:=chooser.ShowBrowseFolder(dlg)
     if err!=nil{walk.MsgBox(dlg,"Folder chooser",err.Error(),walk.MsgBoxIconError);return}
     if yes{folderEdit.SetText(chooser.FilePath)}
    }},
   }},
   Label{Text:"Branch"},
   LineEdit{AssignTo:&branchEdit,Text:"main"},
   Composite{Layout:HBox{},Children:[]Widget{
    HSpacer{},
    PushButton{Text:"Add",OnClicked:func(){
     rawURL:=strings.TrimSpace(urlEdit.Text())
     name,err:=config.GitHubName(rawURL)
     if err!=nil{walk.MsgBox(dlg,"Invalid GitHub URL",err.Error(),walk.MsgBoxIconWarning);return}
     folder:=strings.TrimSpace(folderEdit.Text())
     branch:=strings.TrimSpace(branchEdit.Text())
     if folder==""||branch==""{walk.MsgBox(dlg,"Missing information","Select a parent folder and enter a branch.",walk.MsgBoxIconWarning);return}
     result=config.Repository{Name:name,URL:rawURL,Directory:filepath.Join(folder,name),Branch:branch,Remote:"origin"}
     accepted=true
     dlg.Accept()
    }},
    PushButton{Text:"Cancel",OnClicked:func(){dlg.Cancel()}},
   }},
  },
 }).Run(owner)
 if err!=nil{walk.MsgBox(owner,"Dock",fmt.Sprintf("Cannot open Add New dialog: %v",err),walk.MsgBoxIconError);return config.Repository{},false}
 return result,accepted
}

func tokenDialog(owner *walk.MainWindow)(string,string,bool){
 var dlg *walk.Dialog
 var userEdit,tokenEdit *walk.LineEdit
 var username,token string
 var accepted bool
 _,err:=(Dialog{
  AssignTo:&dlg,Title:"GitHub token - Windows Credential Manager",
  MinSize:Size{Width:660,Height:310},
  Layout:VBox{},
  Children:[]Widget{
   Label{Text:"GitHub username"},
   LineEdit{AssignTo:&userEdit},
   Label{Text:"Fine-grained personal access token (masked; not stored in TOML)"},
   LineEdit{AssignTo:&tokenEdit,PasswordMode:true},
   Label{Text:"Permissions for fine-grained tokens: Repository Contents = Read and write."},
   Label{Text:"Repository access: select only the repositories you will pull and push."},
   Label{Text:"Metadata = Read-only (automatically granted by GitHub). No extra permissions needed."},
   Composite{Layout:HBox{},Children:[]Widget{
    PushButton{Text:"Create GitHub Token",OnClicked:func(){
     constURL:="https://github.com/settings/personal-access-tokens/new"
     if err:=exec.Command("rundll32.exe","url.dll,FileProtocolHandler",constURL).Start();err!=nil{
      walk.MsgBox(dlg,"Open GitHub",err.Error(),walk.MsgBoxIconError)
     }
    }},
    HSpacer{},
   }},
   Label{Text:"Credentials are saved with Git Credential Manager, not in config.toml."},
   Composite{Layout:HBox{},Children:[]Widget{
    HSpacer{},
    PushButton{Text:"Save securely",OnClicked:func(){
     username=strings.TrimSpace(userEdit.Text())
     token=tokenEdit.Text()
     if username==""||strings.TrimSpace(token)==""{walk.MsgBox(dlg,"Missing credentials","Enter GitHub username and token.",walk.MsgBoxIconWarning);return}
     accepted=true
     dlg.Accept()
    }},
    PushButton{Text:"Cancel",OnClicked:func(){dlg.Cancel()}},
   }},
  },
 }).Run(owner)
 if err!=nil{walk.MsgBox(owner,"Dock",err.Error(),walk.MsgBoxIconError);return "","",false}
 return username,token,accepted
}
