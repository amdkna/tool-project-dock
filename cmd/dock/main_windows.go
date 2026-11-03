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

const version = "0.1.2"

func main() {
 executable, err := os.Executable()
 if err != nil { walk.MsgBox(nil, "Dock", err.Error(), walk.MsgBoxIconError); return }
 configPath := filepath.Join(filepath.Dir(executable), "config.toml")
 if err := config.Ensure(configPath); err != nil { walk.MsgBox(nil, "Dock", "Cannot create config: "+err.Error(), walk.MsgBoxIconError); return }

 var mw *walk.MainWindow
 var output *walk.TextEdit
 var busy atomic.Bool

 logLine := func(s string) {
  if mw == nil || output == nil { return }
  mw.Synchronize(func() {
   previous := output.Text()
   if len(previous) > 60000 { previous = previous[len(previous)-45000:] }
   output.SetText(previous + s + "\r\n")
  })
 }

 runGit := func(action string) {
  if !busy.CompareAndSwap(false, true) { walk.MsgBox(mw, "Dock", "A Git operation is already running.", walk.MsgBoxIconInformation); return }
  go func() {
   defer busy.Store(false)
   cfg, err := config.Load(configPath)
   if err != nil { logLine("[ERROR] "+err.Error()); return }
   if info, err := os.Stat(cfg.Project.Directory); err != nil || !info.IsDir() {
    logLine("[ERROR] Project directory not found: "+cfg.Project.Directory)
    return
   }
   args := []string{"-C", cfg.Project.Directory}
   if action == "PULL" { args = append(args, "pull", "--ff-only", cfg.Git.Remote, cfg.Git.Branch) } else { args = append(args, "push", cfg.Git.Remote, cfg.Git.Branch) }
   logLine("["+action+"] git "+strings.Join(args, " "))
   cmd := exec.Command("git", args...)
   cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
   result, err := cmd.CombinedOutput()
   if len(result) > 0 { logLine(strings.TrimRight(string(result), "\r\n")) }
   if err != nil { logLine(fmt.Sprintf("[ERROR] %s failed: %v", action, err)); return }
   logLine("[OK] "+action+" completed")
  }()
 }

 openConfig := func() {
  cmd := exec.Command("notepad.exe", configPath)
  if err := cmd.Start(); err != nil { walk.MsgBox(mw, "Dock", "Cannot open config: "+err.Error(), walk.MsgBoxIconError) }
 }

 _, err = (MainWindow{
  AssignTo: &mw,
  Title: "Dock by Webigram | ver. " + version,
  MinSize: Size{Width: 580, Height: 400},
  Size: Size{Width: 720, Height: 480},
  Layout: VBox{MarginsZero: false, Spacing: 10},
  MenuItems: []MenuItem{
   Menu{Text: "&File", Items: []MenuItem{
    Action{Text: "&Config...", OnTriggered: openConfig},
   }},
   Menu{Text: "&Info", Items: []MenuItem{
    Action{Text: "ver. " + version, Enabled: false},
   }},
  },
  Children: []Widget{
   Label{Text: "Dock by Webigram - Git workspace controls"},
   Composite{Layout: HBox{}, Children: []Widget{
    PushButton{Text: "PULL", OnClicked: func(){runGit("PULL")}},
    PushButton{Text: "PUSH", OnClicked: func(){runGit("PUSH")}},
   }},
   Label{Text: "Activity log"},
   TextEdit{AssignTo: &output, ReadOnly: true, VScroll: true},
   Label{Text: "Configure your repository under File > Config before using Git commands."},
  },
 }).Run()
 if err != nil { walk.MsgBox(nil, "Dock", err.Error(), walk.MsgBoxIconError) }
}
