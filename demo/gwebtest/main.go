package main

import (
	"github.com/nbvghost/gweb"
	"net/http"
	"net/url"
	"fmt"
	"github.com/nbvghost/gweb/conf"
	"log"
)
//拦截器
type InterceptorManager struct {
}
//拦截器 方法，如果 允许登陆 返回true
func (this InterceptorManager) Execute(context *gweb.Context) bool {
	if context.Session.Attributes.Get("admin") == nil { //判断当前 session 信息
		redirect := "" // 跳转地址
		if len(context.Request.URL.Query().Encode()) == 0 {
			redirect = context.Request.URL.Path
		} else {
			redirect = context.Request.URL.Path + "?" + context.Request.URL.Query().Encode()
		}
		http.Redirect(context.Response, context.Request, "/account/login?redirect="+url.QueryEscape(redirect), http.StatusFound)
		return false
	} else {
		return true
	}
}
type User struct {
	Name string
	Age int
}
//index路由控制器
type IndexController struct {
	gweb.BaseController
}
func (c *IndexController) Apply() {
	c.Interceptors.Add(&InterceptorManager{})//拦截器

	//默认index页面
	c.AddHandler("", func(context *gweb.Context) gweb.Result {

		return &gweb.RedirectToUrlResult{"index"}
	})
	// 如果没有地址view里的文件时的映射
	c.AddHandler("*",func(context *gweb.Context) gweb.Result {
		return &gweb.HTMLResult{}
	})
	// 添加index地址映射
	c.AddHandler("index", func(context *gweb.Context) gweb.Result {
		return &gweb.HTMLResult{}
	})


	wx := &WxController{}
	wx.Interceptors = c.Interceptors //使用 父级 拦截器
	c.AddSubController("/wx/", wx) // 添加子控制器，相关的路由定义在 WxController.Apply() 里

}
//account路由控制器
type AccountController struct {
	gweb.BaseController
}
func (c *AccountController) Apply() {


	c.AddHandler("login",  func(context *gweb.Context) gweb.Result {

		user:=&User{Name:"user name",Age:12}

		context.Session.Attributes.Put("admin",user)

		redirect := context.Request.URL.Query().Get("redirect")

		return &gweb.RedirectToUrlResult{Url:redirect}
	})

}

// /wx 路由控制器
type WxController struct {
	gweb.BaseController
}
func (c *WxController) Apply() {


	c.AddHandler(":id/path",  func(context *gweb.Context) gweb.Result {

		user:=context.Session.Attributes.Get("admin").(*User)

		return &gweb.HTMLResult{Name:"wx/path",Params:map[string]interface{}{"User":user,"Id":context.PathParams}}
	})
	c.AddHandler("info", func(context *gweb.Context) gweb.Result {

		user:=context.Session.Attributes.Get("admin").(*User)

		return &gweb.HTMLResult{Params:map[string]interface{}{"User":user}}
	})

}
func main()  {



	var kkk = new(IndexController)
	fmt.Println(kkk)
	kkk =&IndexController{}
	fmt.Println(kkk)


	kkks :=IndexController{}
	fmt.Println(kkks)

	//初始化控制器，拦截 / 路径
	index := &IndexController{}
	index.NewController("/", index)



	//初始化控制器，拦截 /account 路径
	account := &AccountController{}
	account.NewController("/account", account)

	//启动web服务器
	gweb.StartServer(true,false)


	//也可用，内置函数,gweb只是简单的做一个封装的
	//err := http.ListenAndServe(conf.Config.HttpPort, nil)
	//log.Println(err)



}


