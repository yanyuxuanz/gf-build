package main

import (
	"fmt"
	"github.com/gogf/gf/g/os/gfile"
	"os"
	"strings"
	"time"
)
//cli => cli模式
type cli struct {
	AppPath string
	AppName string
	ConfigPath string
	Directory string
	Adis []adi
}
 //adi => app dir info
type adi struct {
	Type string
	Name string
	Comment string
	Parent string
	Content []byte
}
func (client *cli) GetAppPath() string{
	if strings.Contains(client.AppName,"/") {
		client.AppName = strings.Replace(client.AppName,"/","",-1)
	}
	if strings.Contains(client.AppName,"\\") {
		client.AppName = strings.Replace(client.AppName,"\\","",-1)
	}
	return client.AppPath + "\\" + client.AppName
}

var appcli *cli
func init(){

	appcli = new(cli)
	appcli.Directory = `/
├── app	业务逻辑层	所有的业务逻辑存放目录。
│   ├── api	业务接口	接收/解析用户输入参数的入口/接口层。
│   ├── model	数据模型	数据管理层，仅用于操作管理数据，如数据库操作。
│   └── service	逻辑封装	业务逻辑封装层，实现特定的业务需求，可供不同的包调用。
├── boot	初始化包	用于项目初始化参数设置。
├── config	配置管理	所有的配置文件存放目录。
├── docfile	项目文档	DOC项目文档，如: 设计文档、脚本文件等等。
├── library	公共库包	公共的功能封装包，往往不包含业务需求实现。
├── public	静态目录	仅有该目录下的文件才能对外提供静态服务访问。
├── router	路由注册	用于路由统一的注册管理。
├── template	模板文件	MVC模板文件存放的目录。
├── vendor	第三方包	第三方依赖包存放目录(可选, 未来会被淘汰)。
├── go.mod	依赖管理	使用Go Module包管理的依赖描述文件。
└── main.go	入口文件	程序入口文件。`
appcli.AppName = "gfapps"
appcli.AppPath,_ = os.Getwd()
}

func (client *cli) setadis() {
	client.Adis = []adi{
		{Type:"dir", Name:"app",Parent:"/",Comment:"业务逻辑层	所有的业务逻辑存放目录。"},
		{Type:"dir", Name:"api",Parent:"app",Comment:"业务接口	接收/解析用户输入参数的入口/接口层。"},
		{Type:"dir", Name:"user",Parent:"app/api",Comment:"业务接口	接收/解析用户输入参数的入口/接口层。"},
		{Type:"file", Name:"user.go",Parent:"app/api/user",Comment:"业务接口	接收/解析用户输入参数的入口/接口层。",Content:[]byte(`package a_user

import (
    "`+client.AppName+`/app/service/user"
    "`+client.AppName+`/library/response"
    "github.com/gogf/gf/g/net/ghttp"
    "github.com/gogf/gf/g/util/gvalid"
)

// 用户API管理对象
type Controller struct { }

// 用户注册接口
func (c *Controller) SignUp(r *ghttp.Request) {
    if err := s_user.SignUp(r.GetPostMap()); err != nil {
        response.Json(r, 1, err.Error())
    } else {
        response.Json(r, 0, "ok")
    }
}

// 用户登录接口
func (c *Controller) SignIn(r *ghttp.Request) {
    data  := r.GetPostMap()
    rules := map[string]string {
        "passport"  : "required",
        "password"  : "required",
    }
    msgs  := map[string]interface{} {
        "passport" : "账号不能为空",
        "password" : "密码不能为空",
    }
    if e := gvalid.CheckMap(data, rules, msgs); e != nil {
        response.Json(r, 1, e.String())
    }
    if err := s_user.SignIn(data["passport"], data["password"], r.Session); err != nil {
        response.Json(r, 1, err.Error())
    } else {
        response.Json(r, 0, "ok")
    }
}

// 判断用户是否已经登录
func (c *Controller) IsSignedIn(r *ghttp.Request) {
    if s_user.IsSignedIn(r.Session) {
        response.Json(r, 0, "ok")
    } else {
        response.Json(r, 1, "")
    }
}

// 用户注销/退出接口
func (c *Controller) SignOut(r *ghttp.Request) {
    s_user.SignOut(r.Session)
    response.Json(r, 0, "ok")
}

// 检测用户账号接口(唯一性校验)
func (c *Controller) CheckPassport(r *ghttp.Request) {
    passport := r.Get("passport")
    if e := gvalid.Check(passport, "required", "请输入账号"); e != nil {
        response.Json(r, 1, e.String())
    }
    if s_user.CheckPassport(passport) {
        response.Json(r, 0, "ok")
    }
    response.Json(r, 1, "账号已经存在")
}

// 检测用户昵称接口(唯一性校验)
func (c *Controller) CheckNickName(r *ghttp.Request) {
    nickname := r.Get("nickname")
    if e := gvalid.Check(nickname, "required", "请输入昵称"); e != nil {
        response.Json(r, 1, e.String())
    }
    if s_user.CheckNickName(r.Get("nickname")) {
        response.Json(r, 0, "ok")
    }
    response.Json(r, 1, "昵称已经存在")
}`)},
		{Type:"dir", Name:"model",Parent:"app",Comment:"数据模型	数据管理层，仅用于操作管理数据，如数据库操作。"},
		{Type:"dir", Name:"service",Parent:"app",Comment:"逻辑封装	业务逻辑封装层，实现特定的业务需求，可供不同的包调用。"},
		{Type:"dir", Name:"user",Parent:"app/service",Comment:"逻辑封装	业务逻辑封装层，实现特定的业务需求，可供不同的包调用。"},
		{Type:"file", Name:"user.go",Parent:"app/service/user",Comment:"逻辑封装	业务逻辑封装层，实现特定的业务需求，可供不同的包调用。",Content:[]byte(`package s_user

import (
    "errors"
    "fmt"
    "github.com/gogf/gf/g"
    "github.com/gogf/gf/g/net/ghttp"
    "github.com/gogf/gf/g/os/gtime"
    "github.com/gogf/gf/g/util/gvalid"
)

const (
    USER_SESSION_MARK = "user_info"
)

var (
    // 表对象
    table = g.DB().Table("user").Safe()
)

// 用户注册
func SignUp(data g.MapStrStr) error {
    // 数据校验
    rules := []string {
        "passport @required|length:6,16#账号不能为空|账号长度应当在:min到:max之间",
        "password2@required|length:6,16#请输入确认密码|密码长度应当在:min到:max之间",
        "password @required|length:6,16|same:password2#密码不能为空|密码长度应当在:min到:max之间|两次密码输入不相等",
    }
    if e := gvalid.CheckMap(data, rules); e != nil {
        return errors.New(e.String())
    }
    if _, ok := data["nickname"]; !ok {
        data["nickname"] = data["passport"]
    }
    // 唯一性数据检查
    if !CheckPassport(data["passport"]) {
        return errors.New(fmt.Sprintf("账号 %s 已经存在", data["passport"]))
    }
    if !CheckNickName(data["nickname"]) {
        return errors.New(fmt.Sprintf("昵称 %s 已经存在", data["nickname"]))
    }
    // 记录账号创建/注册时间
    if _, ok := data["create_time"]; !ok {
        data["create_time"] = gtime.Now().String()
    }
    if _, err := table.Filter().Data(data).Save(); err != nil {
        return err
    }
    return nil
}

// 判断用户是否已经登录
func IsSignedIn(session *ghttp.Session) bool {
    return session.Contains(USER_SESSION_MARK)
}

// 用户登录，成功返回用户信息，否则返回nil; passport应当会md5值字符串
func SignIn(passport, password string, session *ghttp.Session) error {
    record, err := table.Where("passport=? and password=?", passport, password).One()
    if err != nil {
        return err
    }
    if record == nil {
        return errors.New("账号或密码错误")
    }
    session.Set(USER_SESSION_MARK, record)
    return nil
}

// 用户注销
func SignOut(session *ghttp.Session) {
    session.Remove(USER_SESSION_MARK)
}

// 检查账号是否符合规范(目前仅检查唯一性),存在返回false,否则true
func CheckPassport(passport string) bool {
    if i, err := table.Where("passport", passport).Count(); err != nil {
        return false
    } else {
        return i == 0
    }
}

// 检查昵称是否符合规范(目前仅检查唯一性),存在返回false,否则true
func CheckNickName(nickname string) bool {
    if i, err := table.Where("nickname", nickname).Count(); err != nil {
        return false
    } else {
        return i == 0
    }
}`)},
		{Type:"dir", Name:"boot",Parent:"/",Comment:"初始化包	用于项目初始化参数设置。"},
		{Type:"file", Name:"boot.go",Parent:"boot",Comment:"初始化包	用于项目初始化参数设置。",Content:[]byte(`package boot

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/os/glog"
)

func init() {
	v := g.View()
	c := g.Config()
	s := g.Server()

	// 模板引擎配置
	v.AddPath("template")
	v.SetDelimiters("${", "}")

	// glog配置
	logpath := c.GetString("setting.logpath")
	glog.SetPath(logpath)
	glog.SetStdoutPrint(true)

	// Web Server配置
	s.SetServerRoot("public")
	s.SetLogPath(logpath)
	s.SetNameToUriType(ghttp.NAME_TO_URI_TYPE_ALLLOWER)
	s.SetErrorLogEnabled(true)
	s.SetAccessLogEnabled(true)
	s.SetPort(8199)
}
`)},
		{Type:"dir", Name:"config",Parent:"/",Comment:"配置管理	所有的配置文件存放目录。"},
		{Type:"toml", Name:"config.toml",Parent:"config",Comment:"配置管理	所有的配置文件存放目录。",Content:[]byte(`# 应用系统设置
[setting]
    logpath = "/tmp/log/gf-demos"

# 数据库连接
[database]
    [[database.default]]
        host = "127.0.0.1"
        port = "3306"
        user = "root"
        pass = "12345678"
        name = "test"
        type = "mysql"`)},
		{Type:"dir", Name:"docfile",Parent:"/",Comment:"项目文档	DOC项目文档，如: 设计文档、脚本文件等等。"},
		{Type:"dir", Name:"library",Parent:"/",Comment:"公共库包	公共的功能封装包，往往不包含业务需求实现。"},
		{Type:"dir", Name:"response",Parent:"library",Comment:"公共库包	公共的功能封装包，往往不包含业务需求实现。"},
		{Type:"file", Name:"response.go",Parent:"library/response",Comment:"公共库包	公共的功能封装包，往往不包含业务需求实现。",Content:[]byte(`package response

import (
    "github.com/gogf/gf/g"
    "github.com/gogf/gf/g/net/ghttp"
)

// 标准返回结果数据结构封装。
// 返回固定数据结构的JSON:
// err:  错误码(0:成功, 1:失败, >1:错误码);
// msg:  请求结果信息;
// data: 请求结果,根据不同接口返回结果的数据结构不同;
func Json(r *ghttp.Request, err int, msg string, data...interface{}) {
    responseData := interface{}(nil)
    if len(data) > 0 {
        responseData = data[0]
    }
    r.Response.WriteJson(g.Map{
        "err"  : err,
        "msg"  : msg,
        "data" : responseData,
    })
    r.Exit()
}`)},
		{Type:"dir", Name:"public",Parent:"/",Comment:"静态目录	仅有该目录下的文件才能对外提供静态服务访问。"},
		{Type:"dir", Name:"router",Parent:"/",Comment:"路由注册	用于路由统一的注册管理。"},
		{Type:"dir", Name:"template",Parent:"/",Comment:"模板文件	MVC模板文件存放的目录。"},

		{Type:"mod", Name:"go.mod",Parent:"/",Comment:"依赖管理	使用Go Module包管理的依赖描述文件。"},
		{Type:"file", Name:"main.go",Parent:"/",Comment:"入口文件	程序入口文件。",Content:[]byte(`package main

import (
    _ "`+client.AppName+`/boot"
    
    "github.com/gogf/gf/g"
"github.com/gogf/gf/g/net/ghttp"
)

func main() {
	s := g.Server()
    s.BindHandler("/", func(r *ghttp.Request) {
        r.Response.Write("哈喽世界！")
    })
    s.Run()
}`)},
	}
}
func main() {
	fmt.Println(">Welcome to use app-cli to build your Application!")
	fmt.Printf(">Please enter the Path of your project working directory,default : %s \n",appcli.AppPath)
	fmt.Print(">Enter the Path of your working path press Enter to next:")
	fmt.Scanln(&appcli.AppPath)
	L1:
	fmt.Printf(">Please enter the name of the project you want to build,default : %s \n",appcli.AppName)
	fmt.Print(">Enter your project name,press enter to next step:")
	fmt.Scanln(&appcli.AppName)
	fmt.Println(">Generating projects Named ",appcli.AppName," for you.")

	if gfile.Exists(appcli.GetAppPath()) {
		L2:
		fmt.Printf(">Warning:The project name you gave is already exists in the working directory. Do you want to overwrite it.(y or n)?")
		var ifoverwrite string
		fmt.Scanln(&ifoverwrite)
		if ifoverwrite == "n"{//不覆盖
			goto L1
		}else if ifoverwrite != "y"{
			goto L2
		}
	}
	appcli.setadis()
	appcli.GenerateApp()
	fmt.Println(">Press any key to close this control console.")
	var anykey string
	fmt.Scanln(&anykey)
	fmt.Println("bye~~")
	time.Sleep(time.Second*1)

}

func (client *cli) GenerateApp(){
	fmt.Println(">Show the Project directory structure on the below:\n",appcli.Directory)
	fmt.Println(">Generating......")
	Rootpath := client.GetAppPath()
	gfile.Mkdir(Rootpath)
	os.Chdir(Rootpath)
	var fullpath string
	for _,adi := range client.Adis{
		fullpath = adi.Name
		if adi.Parent != "/"{
			fullpath = adi.Parent+"/"+adi.Name
		}
		if adi.Type == "dir"{
			gfile.Mkdir(fullpath)
		}else {
			file,err := gfile.Create(fullpath)
			defer file.Close()
			if err != nil {
				fmt.Println(err)
			}
			file.Write(adi.Content)
		}
	}
	fmt.Println(">Project Generating completed!")
}