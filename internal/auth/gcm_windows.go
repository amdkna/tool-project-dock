package auth

import (
 "fmt"
 "os"
 "os/exec"
 "strings"
)

// Env isolates Dock credentials from other Git users in Windows Credential Manager.
func Env() []string {
 return append(os.Environ(), "GCM_CREDENTIAL_STORE=wincredman", "GCM_NAMESPACE=webigram-dock")
}

// StoreGitHub sends the token only through stdin, never in args, files or logs.
func StoreGitHub(username, token string) error {
 username=strings.TrimSpace(username)
 token=strings.TrimSpace(token)
 if username=="" || token=="" {return fmt.Errorf("GitHub username and token are required")}
 if strings.ContainsAny(username,"\r\n=")||strings.ContainsAny(token,"\r\n"){return fmt.Errorf("invalid username or token")}
 cmd:=exec.Command("git","credential-manager","store")
 cmd.Env=Env()
 cmd.Stdin=strings.NewReader("protocol=https\nhost=github.com\nusername="+username+"\npassword="+token+"\n\n")
 if _,err:=cmd.CombinedOutput();err!=nil{
  return fmt.Errorf("could not save GitHub token with Git Credential Manager: %w",err)
 }
 return nil
}
