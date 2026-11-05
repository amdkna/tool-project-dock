package main

import (
 "context"
 "fmt"
 "os"
 "os/exec"
 "path/filepath"
 "strings"
 "time"

 "github.com/amdkna/tool-project-dock/internal/auth"
 "github.com/amdkna/tool-project-dock/internal/config"
)

// Check access without changing the remote repository. Read access is tested;
// push permissions cannot be verified safely without actually pushing.
func testRepository(repo config.Repository) error {
 if _,err:=exec.LookPath("git");err!=nil{return fmt.Errorf("Git for Windows is not installed or not on PATH")}
 if _,err:=config.GitHubName(repo.URL);err!=nil{return err}
 if strings.TrimSpace(repo.Branch)==""{return fmt.Errorf("enter a branch")}
 if err:=checkDestination(repo);err!=nil{return err}
 ctx,cancel:=context.WithTimeout(context.Background(),90*time.Second)
 defer cancel()
 // Do not disable GCM browser sign-in; a private repository may require it.
 cmd:=exec.CommandContext(ctx,"git","-c","credential.helper=","-c","credential.helper=manager",
  "ls-remote","--exit-code","--heads",repo.URL,"refs/heads/"+repo.Branch)
 cmd.Env=auth.Env()
 out,err:=cmd.CombinedOutput()
 if ctx.Err()!=nil{return fmt.Errorf("GitHub access check timed out; complete sign-in and try again")}
 if err!=nil{
  return fmt.Errorf("could not read branch %s on GitHub (check URL, branch, credentials, and network): %v\n%s",repo.Branch,err,strings.TrimSpace(string(out)))
 }
 if len(strings.TrimSpace(string(out)))==0{return fmt.Errorf("branch %q was not found",repo.Branch)}
 return nil
}

func checkDestination(repo config.Repository) error {
 parent:=filepath.Dir(repo.Directory)
 info,err:=os.Stat(parent)
 if err!=nil{return fmt.Errorf("parent folder unavailable: %s: %w",parent,err)}
 if !info.IsDir(){return fmt.Errorf("parent path is not a folder: %s",parent)}
 info,err=os.Stat(repo.Directory)
 if os.IsNotExist(err){return nil}
 if err!=nil{return err}
 if !info.IsDir(){return fmt.Errorf("destination is an existing file: %s",repo.Directory)}
 cmd:=exec.Command("git","-C",repo.Directory,"rev-parse","--show-toplevel")
 root,gitErr:=cmd.Output()
 if gitErr==nil {
  // Refuse nested folders of another repository; only its actual root is safe.
  if !strings.EqualFold(filepath.Clean(strings.TrimSpace(string(root))),filepath.Clean(repo.Directory)){
   return fmt.Errorf("destination is inside another Git repository: %s",repo.Directory)
  }
  remote:=exec.Command("git","-C",repo.Directory,"remote","get-url",repo.Remote)
  actual,err:=remote.Output()
  if err!=nil{return fmt.Errorf("existing repository has no %s remote: %w",repo.Remote,err)}
  if strings.TrimSpace(string(actual))!=repo.URL{
   return fmt.Errorf("existing Git remote differs from the requested URL: %s",strings.TrimSpace(string(actual)))
  }
  return nil
 }
 entries,err:=os.ReadDir(repo.Directory)
 if err!=nil{return err}
 if len(entries)>0{return fmt.Errorf("destination is not empty and is not a Git checkout: %s",repo.Directory)}
 return nil // git clone supports an existing empty directory
}
