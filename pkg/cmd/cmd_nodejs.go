package cmd

import (
  "path/filepath"
  "fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var(
  nodejsDownloadUrlsBuilder = map[string]func(version string)string{
    "darwin-arm64": func(version string)string{
      return fmt.Sprintf("https://nodejs.org/dist/%s/node-%s-darwin-arm64.tar.gz",version,version)
    },
    "darwin-amd64":func( version string) string {
      return fmt.Sprintf("https://nodejs.org/dist/%s/node-%s-darwin-x64.tar.gz",version,version)
    },
    "linux-arm64": func( version string) string {
      return fmt.Sprintf("https://nodejs.org/dist/%s/node-%s-linux-arm64.tar.xz",version,version)
    },
    "linux-amd64": func( version string) string {
      return fmt.Sprintf("https://nodejs.org/dist/%s/node-%s-linux-x64.tar.xz",version,version)
    },
  }
)

type NodejsInstallCmd struct{
  Version string
}


func (c *NodejsInstallCmd) Run(cmd *cobra.Command, args []string) error{
  return installNodejs(c.Version)
}

func installNodejs(version string)error{
  userHomePath,err := os.UserHomeDir()
  if err != nil{
    return err
  }

  osAndArch := fmt.Sprintf("%s-%s",runtime.GOOS, runtime.GOARCH)
  downloadURLBuilder := nodejsDownloadUrlsBuilder[osAndArch]
  if downloadURLBuilder == nil{
    return fmt.Errorf("unsupported os and arch: %s",osAndArch)
  }

  downloadURL := downloadURLBuilder(version)
  url,err := url.Parse(downloadURL)
  if err != nil{
    return fmt.Errorf("invalid download url: %s, %w",downloadURL,err)
  }
  filename := filepath.Base(url.Path)
  downloadFilepath := filepath.Join(os.TempDir(),filename)
  os.Remove(downloadFilepath)
  downloadCmd := exec.Command("curl","-o",downloadFilepath,"-L",downloadURL)
  downloadCmd.Stdout = os.Stdout
  downloadCmd.Stderr = os.Stderr

  if err := downloadCmd.Run();err != nil{
    return fmt.Errorf("download nodejs error: %w",err)
  }


  installPath := filepath.Join(userHomePath,"tools",fmt.Sprintf("nodejs-%s",version))
  os.RemoveAll(installPath)
  os.MkdirAll(installPath,os.ModePerm)
  uncompressCmd := exec.Command("tar","-xf",downloadFilepath,"--strip-components=1","-C",installPath)
  uncompressCmd.Stdout = os.Stdout
  uncompressCmd.Stderr = os.Stderr
  if err :=uncompressCmd.Run();err!=nil{
    return fmt.Errorf("uncompress error: %w",err)
  }

  shrc := []string{
    ".bashrc",
    ".zshrc",
  }
  for _,v := range shrc{
    if _,err :=os.Stat(v);err == nil{
      f,err := os.OpenFile(v,os.O_RDWR|os.O_APPEND,os.ModePerm)
      if err != nil{
        return err
      }
      defer f.Close()
      if _,err := f.WriteString(fmt.Sprintf("\nexport PATH=$PATH:$HOME/tools/nodejs-%s",version));err != nil{
        return err
      }
    }
  }

  return nil
}

