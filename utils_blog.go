/*
landers Apps
Author: landers
Github: github.com/landers1037
*/

package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/adrg/frontmatter"
	"github.com/go-yaml/yaml"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const RoutinePerCPU = 4
const MaxRoutine = 20

type Meta struct {
	Name       string
	ID         int
	Title      string
	Date       string
	DatePlus   string
	Tags       []string
	Categories []string
}

type MdData struct {
	Meta Meta
	Abstract string
	Body string
}

// blog parser
// 为了减少缓存的顶级文件读取器
func ReadMd(file string) ([]byte, error) {
	fileRaw, e := ioutil.ReadFile(file)
	if e != nil {
		return nil, e
	}
	return fileRaw, nil
}

// 一级解析器 生成不带meta的md字节
func ParseMd(fileRaw []byte) []byte {
	reg := regexp.MustCompile(`(?s)(^---\n[\s\S].*?\n---\n)(.*)`)
	res := reg.Find(fileRaw)
	res = reg.ReplaceAll(fileRaw, []byte("$2"))

	return res
}


// 二级解析器 生成不带摘要的md字节
func ParseMdAbs(fileRaw []byte) []byte {
	reg := regexp.MustCompile(`(?s).*<!--more-->\n`)
	res := reg.ReplaceAll(fileRaw, []byte(""))

	return res
}

// 三级解析器 生成摘要信息
func ParseAbs(fileRaw []byte) []byte {
	// 摘要不存在时 使用默认摘要
	reg := regexp.MustCompile(`(?s)(.*)(<!--more-->\n)`)
	res := reg.Find(fileRaw)
	if len(res) <= 0 {
		return []byte("<code>Sorry</code>该文章暂无概述💊")
	}
	res = reg.ReplaceAll(res, []byte("$1"))

	return res
}

// meta解析器
// 当前使用正则
func ParseMeta(fileRaw []byte) Meta {
	reg := regexp.MustCompile(`(?s)^---\n([\s\S].*?)\n---\n`)
	res := reg.Find(fileRaw)
	res = reg.ReplaceAll(res, []byte("$1"))
	// start
	regTitle, _ := regexp.Compile("title: (.*)")
	resTitle := regTitle.ReplaceAll(regTitle.Find(res), []byte("$1"))

	regName, _ := regexp.Compile("name: (.*)")
	resName := regName.ReplaceAll(regName.Find(res), []byte("$1"))

	regId, _ := regexp.Compile("id: (.*)")
	resId := regId.ReplaceAll(regId.Find(res), []byte("$1"))

	regDate, _ := regexp.Compile("date: (.*)")
	resDate := regDate.ReplaceAll(regDate.Find(res), []byte("$1"))

	regTags, _ := regexp.Compile("tags: (.*)")
	resTags := regTags.ReplaceAll(regTags.Find(res), []byte("$1"))

	regCate, _ := regexp.Compile("categories: (.*)")
	resCate := regCate.ReplaceAll(regCate.Find(res), []byte("$1"))

	// 解析 存储

	name := strings.TrimSpace(string(resName))
	id, e := strconv.Atoi(string(resId))
	if e != nil {
		id = 0
	}
	title := strings.TrimSpace(string(resTitle))
	date := strings.Split(string(resDate), " ")[0]
	dateplus := string(resDate)
	tags := strings.Fields(string(resTags))
	cates := strings.Fields(string(resCate))

	return Meta{
		Name:       name,
		ID:         id,
		Title:      title,
		Date:       date,
		DatePlus:   dateplus,
		Tags:       tags,
		Categories: cates,
	}
}

// 使用正则提取yaml
func ParseMetaYaml(fileRaw []byte) Meta {
	type yamlData struct {
		Name       string `yaml:"name"`
		ID         int	`yaml:"id"`
		Title      string `yaml:"title"`
		DatePlus   string `yaml:"date"`
		Tags       string `yaml:"tags"`
		Categories string `yaml:"categories"`
	}
	reg := regexp.MustCompile(`(?s)^---\n([\s\S].*?)\n---\n`)
	res := reg.Find(fileRaw)
	res = reg.ReplaceAll(res, []byte("$1"))
	var ym yamlData
	e := yaml.Unmarshal(res, &ym)
	if e != nil {
		fmt.Println("解析文件meta信息失败")
		os.Exit(1)
	}

	name := strings.TrimSpace(string(ym.Name))
	id, e := strconv.Atoi(string(ym.ID))
	if e != nil {
		id = 0
	}
	title := strings.TrimSpace(string(ym.Title))
	date := strings.Split(string(ym.DatePlus), " ")[0]
	dateplus := string(ym.DatePlus)
	tags := strings.Fields(string(ym.Tags))
	cates := strings.Fields(string(ym.Categories))

	return Meta{
		Name:       name,
		ID:         id,
		Title:      title,
		Date:       date,
		DatePlus:   dateplus,
		Tags:       tags,
		Categories: cates,
	}
}

// 使用yaml front解析yaml头部信息
func ParseYamlFront(fileRaw []byte) Meta {
	type yamlData struct {
		Name       string `yaml:"name"`
		ID         int	`yaml:"id"`
		Title      string `yaml:"title"`
		DatePlus   string `yaml:"date"`
		Tags       []string `yaml:"tags"`
		Categories []string `yaml:"categories"`
	}

	var ym yamlData
	_, _ = frontmatter.Parse(bytes.NewReader(fileRaw), &ym)
	name := strings.TrimSpace(ym.Name)
	id, e := strconv.Atoi(string(ym.ID))
	if e != nil {
		id = 0
	}
	title := strings.TrimSpace(ym.Title)
	date := strings.Split(ym.DatePlus, " ")[0]
	dateplus := ym.DatePlus
	tags := ym.Tags
	cates := ym.Categories

	return Meta{
		Name:       name,
		ID:         id,
		Title:      title,
		Date:       date,
		DatePlus:   dateplus,
		Tags:       tags,
		Categories: cates,
	}
}

// blog new
func NewMd(name string) error {
	// 生成带元信息头部的新md文件
	// 默认以title名称作为唯一文件名
	// check file exist
	_, e := os.Stat(name + ".md")
	if e != nil {
		if os.IsExist(e) {
			return errors.New(fmt.Sprintf("file %s exist", name + ".md"))
		}else {
			// file not exist
			f, e := os.OpenFile(name + ".md", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
			if e != nil {
				return e
			}
			var metainfo string
			if res, e := GetYamlFront();e != nil {
				metainfo = res
			}
			metainfo = "---\ntitle: %s\nname: %s\ndate: %s\nid: 0\ntags: \ncategories: \nabstract: \n---\n<!--more-->"
			dateString := time.Now().Format("2006-01-02 15:04:05")
			metainfo = fmt.Sprintf(metainfo,name, name, dateString)
			_, e = f.WriteString(metainfo)
			if e != nil {
				return e
			}
			return nil
		}
	}else {
		return errors.New(fmt.Sprintf("file %s exist", name + ".md"))
	}
}

// 自定义yaml文件头部
func EditYamlFront() error {
	// 这里限制了文件的来源为/etc/jjtool.default
	// 调用shell来编辑所以无论如何能保证存在
	cmd := exec.Command("bash", "-c", "vi /etc/jjtool.default")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// 读取自定义的yaml头部
func GetYamlFront() (string, error) {
	_, e := os.Stat("/etc/jjtool.default")
	if e != nil {
		// 不存在时写入默认的
		metainfo := "---\ntitle: %s\nname: %s\ndate: %s\nid: 0\ntags: \ncategories: \nabstract: \n---\n<!--more-->"
		_ = ioutil.WriteFile("/etc/jjtool.default", []byte(metainfo), 0644)
		return "", e
	}
	data, e := ioutil.ReadFile("/etc/jjtool.default")
	return fmt.Sprintf("%s", data), e
}

// by db
// 根据数据库id 生成带自增id的md
func NewMdByDB(name string, db string) error {
	return nil
}

// md文件byte转html原生字符串
func MarkdownToHtml(md []byte) string {
	return fmt.Sprintf("%s", md)
}

// 生成最终的要存入数据库前的结构体数据
func GenMdData(file string) MdData {
	fileRaw, e := ReadMd(file)
	if e != nil {
		fmt.Printf("文件%s解析失败 %s\n",file, e.Error())
		panic("异常退出")
	}
	md := ParseMd(fileRaw)
	mdAbs := ParseAbs(md)
	meta := ParseYamlFront(fileRaw)
	htmlAbs := MarkdownToHtml(mdAbs)
	htmlBody := MarkdownToHtml(md)

	return MdData{
		Meta:     meta,
		Abstract: htmlAbs,
		Body:     htmlBody,
	}
}

func GetFileFromPath(fpath string) []string {
	files, e := ioutil.ReadDir(fpath)
	if e != nil {
		fmt.Println("文件目录读取失败")
		return []string{}
	}else {
		var ff []string
		for _, f := range files {
			ext := strings.Split(f.Name(), ".")
			index := len(ext) - 1
			if len(ext) >= 2 && ext[index] == "md" {
				ff = append(ff, path.Join(fpath, f.Name()))
			}
		}
		return ff
	}
}

// 基于协程的md渲染器
// 传入文件数组 解析文件后返回总的内存数据
// 根据cpu数量限制协程数
func ProcessAllMd(fileList []string, userRoutine int) []MdData {
	cpus := runtime.NumCPU()
	// 默认假设一个cpu的可用协程为50
	max := RoutinePerCPU * cpus

	fileLen := len(fileList)
	var AllMd []MdData
	fmt.Printf("开始分配协程: %d\n", userRoutine)

	ch := make(chan string)
	chData := make(chan MdData)
	// 基于文件数目分配 协程不足时分段读取
	if fileLen <= userRoutine {
		for _, file := range fileList {
			go func(file string) {
				d := GenMdData(file)
				time.Sleep(1500 * time.Millisecond)
				ch<-d.Meta.Name
				chData<-d
			}(file)
		}
	}else if fileLen > userRoutine && userRoutine <= max && userRoutine >0 {
		// 动态分配协程
		// 创建指定数目的文章切片
		var mdFileList [][]string
		for i:=0;i<userRoutine;i++ {
			mdFileList = append(mdFileList, []string{})
		}

		for i:=0;i<fileLen;i++ {
			mdFileList[i%userRoutine] = append(mdFileList[i%userRoutine], fileList[i])
		}
		// 开启解析
		for _, mdList := range mdFileList {
			go func(mdList []string) {
				for _, file := range mdList {
					d := GenMdData(file)
					time.Sleep(1500 * time.Millisecond)
					ch<-d.Meta.Name
					chData<-d
				}
			}(mdList)
		}

	}else {
		fmt.Printf("不支持的输入线程数 你的系统仅支持最大%d的线程数\n", max)
		os.Exit(1)
	}

	for _ = range fileList {
		rec := <-ch
		data := <-chData
		AllMd = append(AllMd, data)
		fmt.Printf("%s 解析完毕\n", rec)
	}
	defer close(ch)
	defer close(chData)
	// 开始写入数据库

	return AllMd
}


func ProcessAllMd2(fileList []string) []MdData {
	cpus := runtime.NumCPU()
	// 默认假设一个cpu的可用协程为50
	max := RoutinePerCPU * cpus
	fileLen := len(fileList)
	var AllMd []MdData
	fmt.Printf("开始分配协程: %d\n", fileLen)

	ch := make(chan int)
	chData := make(chan MdData)
	// 基于文件数目分配 协程不足时分段读取
	if fileLen <= max {
		for i, file := range fileList {
			go func(i int, file string) {
				d := GenMdData(file)
				time.Sleep(1500 * time.Millisecond)
				ch<-i
				chData<-d
			}(i, file)
		}
	}else {
		fmt.Printf("暂不支持此数量协程: %d\n", fileLen)
	}

	for _ = range fileList {
		rec := <-ch
		data := <-chData
		AllMd = append(AllMd, data)
		fmt.Printf("%d 解析完毕\n", rec)
	}
	defer close(ch)
	defer close(chData)
	// 开始写入数据库

	return AllMd
}

// 数据库操作相关

// 生成新数据库
// 需要指定数据库名称 默认为blog.db
// 因为新添加的文章id递增 为保持逻辑一致性 id0为pin id越大表示文章越新
func CreateDB(data []MdData, db string) error {
	var dbName string
	// 保证文件不存在
	_, e := os.Stat(db)
	if e == nil {
		return errors.New("数据库已经存在")
	}
	if db == "" {
		dbName = "blog.db"
	}else {
		dbName = db
	}
	dbCon, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		return err
	}
	dbCon.SingularTable(true)
	// 开始根据数据创建数据库
	fmt.Println("开始创建表")
	dbCon.CreateTable(DB_BLOG_POST{})
	dbCon.CreateTable(DB_BLOG_MESSAGES{})
	dbCon.CreateTable(DB_BLOG_TAGS{})
	dbCon.CreateTable(DB_BLOG_CATES{})
	dbCon.CreateTable(DB_BLOG_COMMENTS{})
	dbCon.CreateTable(DB_BLOG_LIKES{})
	dbCon.CreateTable(DB_BLOG_SHARE{})
	dbCon.CreateTable(DB_BLOG_VIEWS{})
	dbCon.CreateTable(DB_BLOG_ADMIN{})
	dbCon.CreateTable(DB_BLOG_SUBSCRIBE{})

	// 开始存库
	// 这里假设md文件一定是存在id且按顺序的， 没有的会随机存储 后期可以更新id
	//sort.SliceStable(data, func(i, j int) bool {
	//	return data[i].Meta.ID < data[j].Meta.ID
	//})
	// 注意pin文章是后面自己更新定义的
	fmt.Println("开始写入文章数据")
	for _, d := range data {
		post := DB_BLOG_POST{
			ID:         d.Meta.ID,
			Name:       d.Meta.Name,
			Title:      d.Meta.Title,
			Date:       d.Meta.Date,
			DatePlus:   d.Meta.DatePlus,
			Update:     "",
			Abstract:   d.Abstract,
			Content:    d.Body,
			Tags:       strings.Join(d.Meta.Tags, " "),
			Categories: strings.Join(d.Meta.Categories, " "),
			Pin:        0,
		}
		dbCon.Model(&DB_BLOG_POST{}).Create(&post)

		// tags
		for _, t := range d.Meta.Tags {
			tag := DB_BLOG_TAGS{
				Tag:   t,
				Name:  d.Meta.Name,
			}
			dbCon.Model(&DB_BLOG_TAGS{}).Create(&tag)
		}

		// cates
		for _, c := range d.Meta.Categories {
			cate := DB_BLOG_CATES{
				Cate:  c,
				Name:  d.Meta.Name,
			}
			dbCon.Model(&DB_BLOG_CATES{}).Create(&cate)
		}
	}

	defer dbCon.Close()
	fmt.Println("确保数据库句柄关闭完成")
	fmt.Println("如需添加置顶文章 请使用更新操作 更新pin=1")
	return nil
}

// 增量更新数据库
// 全量更新文章
func UpdatDB(file string, db string) error {
	return nil
}

// 更新任意表数据
func UpdateTable(table, where, data, db string) error {
	dbc, e := gorm.Open("sqlite3", db)
	if e != nil {
		return e
	}
	return dbc.Table(table).Where(where).Update(data).Error
}

// 更新pin文章
func UpdatePin(name string, db string) error {
	fmt.Println("你指定的文章是" + name)
	fmt.Println("开始尝试连接数据库" + db)
	dbCon, e := gorm.Open("sqlite3", db)
	if e != nil {
		return e
	}
	// 因为全局只能有一个pin 所以先查询
	var lastPin DB_BLOG_POST
	if dbCon.Model(DB_BLOG_POST{}).Where(&DB_BLOG_POST{Pin: 1}).First(&lastPin).Error == nil {
		lastPin.Pin = 0
		dbCon.Save(&lastPin)
	}
	e = dbCon.Model(DB_BLOG_POST{}).Where(&DB_BLOG_POST{Name: name}).Update("pin", 1).Error
	if e != nil {
	 	return e
	}
	return nil
}
// 更新文章的id

// 更新标签

// 更新分类

// 删除留言

// 更新likes

// 更新share

// 更新comments

// 博文联合迁移
// 读取数据库内容 然后重新生成符合meta的新文章
func MigrateMd() error {
	return nil
}