/*
landers Apps
Author: landers
Github: github.com/landers1037
*/

package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"os"
)
func Entry() {
	app := &cli.App{
		Name:                   APPName,
		Usage:                  "bt -h",
		UsageText:              "bt是一个md一键式工具",
		Version:                Version + " " + Build,
		Description:            "github.com/landers1037/bt",
		Commands:               initCmds(),
		Flags:                  nil,
		EnableBashCompletion:   true,
		HideHelp:               false,
		HideHelpCommand:        false,
		HideVersion:            false,
		Action:                 nil,
		Copyright:              CopyRight,
		UseShortOptionHandling: false,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func initCmds() ([]*cli.Command) {
	return []*cli.Command{
		{
			Name: "new",
			Usage: "new example",
			UsageText: "新建一篇markdown文章",
			Aliases: []string{"n", "new", "newmd"},
			Category: "md",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name: "e",
					Usage: "new -e",
					Aliases: []string{"edit"},
					Value: false,
				},
				&cli.BoolFlag{
					Name: "s",
					Usage: "new -s",
					Aliases: []string{"show"},
					Value: false,
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("e") {
					return EditYamlFront()
				}

				if c.Bool("s") {
					data, e := GetYamlFront()
					if e != nil {
						return e
					}
					fmt.Printf("%s\n", data)
					return nil
				}
				fileName := c.Args().First()
				if fileName != "" {
					return NewMd(fileName)
				}else {
					return errors.New("未指定文件名称")
				}
			},
		},
		{
			Name: "test",
			Usage: "test file.md",
			UsageText: "测试markdown文章是否解析正确",
			Aliases: []string{"t", "test", "testmd"},
			Category: "md",
			Action: func(c *cli.Context) error {
				fileName := c.Args().First()
				if fileName != "" {
					fileRaw, e := ReadMd(fileName)
					if e != nil {
						return e
					}
					fmt.Printf("文章meta元信息如下: \n%+v\n", ParseMeta(fileRaw))
					fmt.Printf("文章正文内容如下: \n%s\n", ParseMd(fileRaw))
					return nil
				}else {
					return errors.New("未指定文件")
				}
			},
		},
		{
			Name: "db",
			Usage: "db -t N [file Path] [db name]",
			UsageText: "创建博客数据库",
			Aliases: []string{"db", "newdb", "dbcreate"},
			Category: "database",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name: "t",
					Aliases: []string{"thread"},
					Usage: "-t 10",
					Value: 10,
				},
			},
			Action: func(c *cli.Context) error {
				filePath := c.Args().First()
				var dbName string
				if c.Args().Len() >= 2 {
					dbName = c.Args().Slice()[1]
				}
				if filePath != "" {
					files := GetFileFromPath(filePath)
					mdDatas := ProcessAllMd(files, c.Int("t"))
					return CreateDB(mdDatas, dbName)
				}else {
					return errors.New("未指定文章所在路径")
				}
			},
		},
		{
			Name: "dbu",
			Usage: "dbu [file] [db name]",
			UsageText: "更新博客数据库",
			Aliases: []string{"dbu", "updatedb", "dbupdate"},
			Category: "database",
			Action: func(c *cli.Context) error {
				fileName := c.Args().First()
				var dbName string
				if c.Args().Len() >= 2 {
					dbName = c.Args().Slice()[1]
				}
				if fileName != "" {
					return  UpdatDB(fileName, dbName)
				}else {
					return errors.New("未指定文章")
				}
			},
		},
		{
			Name: "dbt",
			Usage: "dbt [table] [where] [data] [db name]",
			UsageText: "更新博客数据库表数据",
			Aliases: []string{"dbt", "updatetable", "dbtable"},
			Category: "database",
			Action: func(c *cli.Context) error {
				var table string
				var where string
				var data string
				var dbName string
				if c.Args().Len() >= 4 {
					return UpdateTable(table, where, data, dbName)
				} else {
					return errors.New("输入的参数不足")
				}
			},
		},
	}
}