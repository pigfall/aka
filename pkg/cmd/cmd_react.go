package cmd

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type InitReactUILibraryCmd struct {
	Name string
}

type InitReactAppCmd struct {
	Name string
}

func (c *InitReactUILibraryCmd) Run(cmd *cobra.Command, args []string) error {
	if c.Name == "" {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		c.Name = filepath.Base(dir)
	}

	tpl, err := template.ParseFS(assets, "assets/react/ui-library/*.tpl")
	if err != nil {
		panic(err)
	}

	files, err := assets.ReadDir("assets/react/ui-library")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".tpl" {
			continue
		}
		if err := c.executeTpl(tpl, file.Name()); err != nil {
			panic(err)
		}
	}

	os.RemoveAll("src")

	srcSub, err := fs.Sub(assets, "assets/react/ui-library/src")
	if err != nil {
		panic(err)
	}
	if err := os.CopyFS("src", srcSub); err != nil {
		panic(err)
	}

	return nil
}

func (c *InitReactAppCmd) Run(cmd *cobra.Command, args []string) error {
	if c.Name == "" {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		c.Name = filepath.Base(dir)
	}

	tpl, err := template.ParseFS(assets, "assets/react/app/*.tpl")
	if err != nil {
		panic(err)
	}

	files, err := assets.ReadDir("assets/react/app")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".tpl" {
			if file.IsDir() {
				continue
			}
			if err := c.copyFile(file.Name()); err != nil {
				panic(err)
			}
			continue
		}
		if err := c.executeTpl(tpl, file.Name()); err != nil {
			panic(err)
		}
	}

	os.RemoveAll("src")

	srcSub, err := fs.Sub(assets, "assets/react/app/src")
	if err != nil {
		panic(err)
	}
	if err := os.CopyFS("src", srcSub); err != nil {
		panic(err)
	}

	return nil
}

func (c *InitReactAppCmd) executeTpl(tpl *template.Template, tplFileName string) error {
	f, err := os.Create(strings.TrimSuffix(tplFileName, ".tpl"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tpl.ExecuteTemplate(f, tplFileName, map[string]any{"Name": c.Name})
}

func (c *InitReactAppCmd) copyFile(srcFileName string) error {
	srcFile, err := assets.Open(filepath.Join("assets/react/app", srcFileName))
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(srcFileName)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	return err
}

func (c *InitReactUILibraryCmd) executeTpl(tpl *template.Template, tplFileName string) error {
	f, err := os.Create(strings.TrimSuffix(tplFileName, ".tpl"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tpl.ExecuteTemplate(f, tplFileName, map[string]any{"Name": c.Name})
}
