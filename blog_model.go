/*
landers Apps
Author: landers
Github: github.com/landers1037
*/

package main

// 数据库模型
// 数据库模型
type Model struct {
	PrimaryID int `gorm:"primary_key" json:"primary_id"`
	//CreatedOn int `json:"created_on"`
	//ModifiedOn int `json:"modified_on"`
}

type DB_BLOG_POST struct {
	Model
	ID int `gorm:"not null" json:"id"`
	Name string `gorm:"unique;not null" json:"name"`
	Title string `json:"title"`
	Date string `json:"date"`
	DatePlus string `json:"date_plus"`
	Update string `json:"update"`
	Abstract string `json:"abstract"`
	Content string `json:"content"`
	Tags string `json:"tags"`
	Categories string `json:"categories"`
	Pin int `json:"pin"`
}

type DB_BLOG_TAGS struct {
	Model
	Tag string `json:"tag"`
	Name string `json:"name"`
}

type DB_BLOG_CATES struct {
	Model
	Cate string `json:"cate"`
	Name string `json:"name"`
}

// 总views name=ALL/all
type DB_BLOG_VIEWS struct {
	Model
	Name string `gorm:"unique;not null" json:"name"`
	View int `json:"view"`
}

type DB_BLOG_LIKES struct {
	Model
	Name string `gorm:"unique;not null" json:"name"`
	Like int `json:"like"`
}

type DB_BLOG_SHARE struct {
	Model
	Name string `gorm:"unique;not null" json:"name"`
	Share int `json:"share"`
}

type DB_BLOG_COMMENTS struct {
	Model
	Name string `gorm:"not null" json:"name"`
	Date string `json:"date"`
	Comment int `json:"comment"`
}

type DB_BLOG_MESSAGES struct {
	Model
	User string `json:"user"`
	Date string `json:"date"`
	Message string `gorm:"not null" json:"message"`
}

type DB_BLOG_ADMIN struct {
	Model
	UserName string `gorm:"unique;not null" json:"user_name"`
	PassWd string `json:"passwd"`
	Date string `json:"date"`
}

type DB_BLOG_SUBSCRIBE struct {
	Model
	Mail string `json:"mail"`
	SubscribeDate string `json:"subscribe_date"`
	Period string `json:"period"`
}

type DB_BLOG_ZHUANLAN struct {
	Model
	Name string `gorm:"unique;not null" json:"name"`
	Title string `json:"title"`
	Date string `json:"date"`
	Posts string `json:"posts"`
}