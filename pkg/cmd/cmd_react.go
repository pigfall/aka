package cmd

import (
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

func (c *InitReactUILibraryCmd) executeTpl(tpl *template.Template, tplFileName string) error {
	f, err := os.Create(strings.TrimSuffix(tplFileName, ".tpl"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tpl.ExecuteTemplate(f, tplFileName, map[string]any{"Name": c.Name})
}
