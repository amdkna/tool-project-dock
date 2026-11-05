package main

import (
 "fmt"
 "os"
 "os/exec"
 "path/filepath"
 "strings"

 "github.com/amdkna/tool-project-dock/internal/config"
 "github.com/lxn/walk"
 . "github.com/lxn/walk/declarative"
)

// addRepositoryDialog chooses the exact destination or a repository subfolder.
// Test checks URL, branch, credentials and destination before cloning.
func addRepositoryDialog(owner *walk.MainWindow)(config.Repository,bool){
 var dlg *walk.Dialog
 var urlEdit,folderEdit,branchEdit *walk.LineEdit
 var directCheck *walk.CheckBox
 var testButton *walk.PushButton
 var result config.Repository
 var accepted bool

 collect:=func()(config.Repository,error){
  rawURL:=strings.TrimSpace(urlEdit.Text())
  name,err:=config.GitHubName(rawURL)
  if err!=nil{return config.Repository{},err}
  folder:=strings.TrimSpace(folderEdit.Text())
  branch:=strings.TrimSpace(branchEdit.Text())
  if folder=="" || branch==""{return config.Repository{},fmt.Errorf("choose a local folder and enter a branch")}
  info,err:=os.Stat(folder)
  if err!=nil{return config.Repository{},fmt.Errorf("selected folder is not accessible: %w",err)}
  if !info.IsDir(){return config.Repository{},fmt.Errorf("selected path must be a folder")}
  destination:=folder
  if !directCheck.Checked(){destination=filepath.Join(folder,name)}
  return config.Repository{Name:name,URL:rawURL,Directory:filepath.Clean(destination),Branch:branch,Remote:"origin"},nil
 }

 _,err:=(Dialog{
  AssignTo:&dlg,
  Title:"Add a GitHub repository",
  MinSize:Size{Width:620,Height:340},
  Layout:VBox{},
  Children:[]Widget{
   Label{Text:"GitHub HTTPS URL"},
   LineEdit{AssignTo:&urlEdit,Text:"https://github.com/owner/repository.git"},
   Label{Text:"Local destination folder"},
   Composite{Layout:HBox{},Children:[]Widget{
    LineEdit{AssignTo:&folderEdit},
    PushButton{Text:"Browse...",OnClicked:func(){
     chooser:=walk.FileDialog{Title:"Select a local folder"}
     yes,err:=chooser.ShowBrowseFolder(dlg)
     if err!=nil{walk.MsgBox(dlg,"Folder chooser",err.Error(),walk.MsgBoxIconError);return}
     if yes{folderEdit.SetText(chooser.FilePath)}
    }},
   }},
   CheckBox{AssignTo:&directCheck,Text:"Clone directly into selected folder (no repository subfolder)"},
   Label{Text:"Unchecked: selected folder / repository name. Checked: selected folder itself."},
   Label{Text:"Branch"},
   LineEdit{AssignTo:&branchEdit,Text:"main"},
   Label{Text:"Test verifies read access and destination. Push access requires an actual push."},
   Composite{Layout:HBox{},Children:[]Widget{
    HSpacer{},
    PushButton{AssignTo:&testButton,Text:"Test",OnClicked:func(){
     repo,err:=collect()
     if err!=nil{walk.MsgBox(dlg,"Test - invalid details",err.Error(),walk.MsgBoxIconWarning);return}
     testButton.SetEnabled(false)
     go func(){
      checkErr:=testRepository(repo)
      dlg.Synchronize(func(){
       testButton.SetEnabled(true)
       if checkErr!=nil{
        walk.MsgBox(dlg,"Repository test failed",checkErr.Error(),walk.MsgBoxIconWarning)
       }else{
        walk.MsgBox(dlg,"Repository test passed","GitHub branch is accessible and the local destination looks usable. Push authorization has not been tested.",walk.MsgBoxIconInformation)
       }
      })
     }()
    }},
    PushButton{Text:"Add",OnClicked:func(){
     repo,err:=collect()
     if err!=nil{walk.MsgBox(dlg,"Invalid repository",err.Error(),walk.MsgBoxIconWarning);return}
     result=repo
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
