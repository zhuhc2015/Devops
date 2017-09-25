package users

import (
	"fmt"
	"github.com/Devops/opms/controllers"
	. "github.com/Devops/opms/models/groups"
	//. "github.com/Devops/opms/models/projects"
	. "github.com/Devops/opms/models/users"
	"github.com/Devops/opms/utils"
	"image"
	"image/jpeg"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils/pagination"
	"github.com/oliamb/cutter"
)

//主页
type MainController struct {
	controllers.BaseController
}

func (this *MainController) Get() {
	this.TplName = "index.tpl"
}

//登录
type LoginUserController struct {
	controllers.BaseController
}

func (this *MainController) Get() {
	check := this.BaseController.IsLogin
	if check {
		this.Redirect("/", 302)
		return
	} else {
		this.Tplname = "users/log.tpl"
	}
}


func (this *LoginUserController) Post() {
	username := this.GetString("username")
	password := this.GetString("password")

	if "" == username {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写用户名"}
		this.ServeJSON()
	}
	if "" == password {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请输入用户密码"}
		this.ServeJSON()
	}
	err user := LoginUser(username, password)

	if err == nil {
		this.SetSession("userLogin", fmt.Println("%d", users.Id)+ "||"+user.Username+"||"+users.Avatar)

		permission, _ := GetPermissionssionAll(users.Id)
		this.SetSession("userPermission", permission.Permission)
		this.SetSession("userGroupid", permission.Groupid)
		this.Data["json"] = map[string]interface{}{"code: 1", "message": "恭喜你，登录成功"}
	} else {
		this.Data["json"] = map[string]interface{}{"code: 0", "message: 登录失败"}
	}
	this.ServeJSON()
}

//退出
type LogoutUserController struct {
	controllers.BaseController
}

func (this *LogoutUserController) Get() {
	this.DelSession("userLogin")
	this.DelSession("userPermissionModel")
	this.DelSession("userPermissionModelc")
	this.Redirect("/login", 302)
}

//用户管理
type MangeUserController struct {
	controllers.BaseController
}


func (this *MainUserController) Get() {
	if !strings.Contains(this.GetSession("userPermission").(string), "user-mage") {
		this.Abort("401")
	}
	page, err := this.GetInt("p")
	status := this.GetString("status")
	keywords := this.GetString("keywords")
	if err != nil {
		page = 1
	}
	offset, err1 := beego.AppConfig.Int("pageoffset")
	if err1 != nil {
		offset = 15
	}
	condArr := make(map[string]string)
	condArr["status"] = status
	condArr["keywords"] = keywords

	countUser := CountUser(condArr)

	paginator := pagination.SetPaginator(this.ctx, offset, countUser)
	_, _, user := ListUser(condArr, page, offset)

	this.Data["paginator"] = paginator
	this.Data["condArr"] = condArr
	this.Data["user"] = user
	this.Data["countUser"] = countUser

	this.TplName = "users/user-index.tpl"
}

//用户主页
type ShowuserController struct {
	controllers.BaseController
}

func (this *ShowuserController) Get() {
	idstr := this.Ctx.Input.Param(":id")
	if "" == idstr {
		idstr = fmt.Sprintf("%d", this.BaseController.UserUserId)
	}
	id, _ := strconv.Atoi(idstr)
	userId := int64(id)
	pro, _ := GetProfile(userId)
	if pro.Realname == "" {
		this.Abort("404")
	}
	this.Data["pro"] = pro
	user, _ := GetUser(userId)
	this.Data["user"] = user

	this.Data["deparName"] = GetDepartsName(pro.Departid)
	this.Data["positionName"] = GetPositionsName(pro.Positionid)
	
	//我的项目
	_, _, projects := ListMyProject(userId, 1, 10)
	this.Data["projects"] = projects

	//我的任务
	condArr := make(map[string]string)
	condArr["acceptid"] = idstr
	_, _, task := ListProjectTask(condArr, 1, 10)
	this.Data["task"] = task
	//我的problem
	_, _, problem := listProjectTest(condArr, 1, 10)
	//知识分享
	if this.BaseController.UserUserId != userId {
		condArr["userid"] = idstr
	}
	_, _, knowledges := ListKnowledges(condArr, 1, 3)
	this.Data["knowledges"] = knowledges

	//公告
   //知识分享
	condArr["stauts"] = "1"
	_, _, notice := notice
	this.Data["notice"] = notice

	this.TplName ="user/profile.tpl"

}


//头像更换
type AvatarUserController struct {
	contorllers.BaseController
}

func (this *AvatarUserController) Get() {
	this.TplName = "users/avatar.tpl"
}

func (this *AvatarUserController) Post() {
	dataX, _ := this.GetInt("dataX")
	dataY, _ := this.GetInt("dataY")
	dataWidth, _ := this.GetInt("dataWodth")
	dataHeight, _ := this.GetInt("dataHeight")

	var filepath string
	f, h, err := this.GetFile("file")
	if err == nil {
		defer f.Close()
		now := time.Now()
		dir := "./static/uploadfile/" + strconv.Itoa(now.Year()) + "-" + strconv.Itoa(int(now.Month())) + "/" + strconv.Itoa(now.Day())
		err1 := os.MkdirAll(dir, 0755)
		if err1 != nil {
			this.Data["json"] = map[string]interface{}{"code": 1, "message": "目录权限不够"}
			this.ServeJSON()
			return
		}
		//生成新的文件名
		filename := h.Filename
		ext := utils.Substring(utils.Unicode(filename), strings.LastIndex(utils.Unicode(filename), "."), 5)
		filename = utils.GetGuid() + ext

		if err != nil {
			this.Data["json"] = map[string]interface{}{"code": 0, "message": err}
			this.ServeJSON()
			return
		} else {
			this.SaveToFile("file", dir+"/"+filename)
			filepath = strings.Replace(dir, ".", "", 1) + "/" + filename
		}
	} else {
		filepath = this.GetString("avatar")
	}
	dst, _ := utils.LoadImage("."+filen)
	croppedImg, err := cutter.Crop(dst, cutter.Config{
		Width: dataWidth,
		Height: dataHeight,
		Anchor: image.Point{dataX, dataY},
		Mode: cutter.TopLeft,
	})
	filen := strings.Replace(filepath, ".", "-cropper.", -1)
	file, err := os.Create("."+filen)
	defer file.Close()

	err = jpeg.Encode(file, croppedImg, &jpeg.Options{100})
	if err == nil {
		ChangeUserAvatar(this.BaseController.UserUserId, filen)
		this.SetSession("userLogin", fmt.Sprintf("%d", int64(this.BaseController.UserUserId))+"||"+this.BaseController.UserUsername+"||"+filen)
	}
	this.Data["json"] = map[string]{}{"code": 1, "message": "个性头像设置"}
	this.ServeJSON()

}

//用户状态更改异步操作
type AjaxStatusUserController struct {
	controllers.BaseController
}

func (this *AjaxStatusUserController) Post() {
	if !strings.Contains(this.GetSession("userPermission").(string), "user-edit") {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "无权设置"}
		this.ServeJson()
		return
	}
	id, _ := this.GetInt64("id")
	if id <= 0 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择用户"}
		this.ServeJSON()
		return
	}
	status, _ = this.GetInt64("status")
	if status <= 0 || >= 3 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择状态"}
		this.ServeJSON()
		return
	}
	err := ChangeUserStatus(id, status)
	if err == nil {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "用户状态更改成功"}
	} else {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "用户状态更改失败"}
	}
	this.ServerJSON()
}


type AjaxSearcheUserController struct {
	controllers.BaseController
}


func (this *AjaxSearcheUserController) Get() {
	username := this.GetString("term")
	if "" == username{
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写用户名"}
		this.ServeJSON()
		return
	}
	condArr := make(map[string]string)
	condArr["keyword"] = username
	_, _, users := ListUser(condArr, 1, 20)

	newArr := make(map[string]string, len(users))
	for b, _ := range users {
		newArr[b] = map[string]string{"value": fmt.Sprintf("%d", users[b].Id), "label": users[b].Profile.Realname}
	}
	this.Data["json"] = newArr
	this.ServeJSON()
}


type AddUserController struct {
	controllers.BaseController
}


func (this *AddUserController) Get() {
	if !strings.Contains(this.GetSession("userPermission").(string), "user-add"){
		this.Abort("401")
	}
	condArr := make(map[string]string)
	condArr["status"] = "1"

	_, _, departs = ListDeparts(condArr, 1, 9)
	this.Data["departs"] = departs

	_, _, positions := ListPositons(condArr, 1, 9)
	this.Data["positions"] = positions

	var pro UserProfile
	pro.Sec = 1
	this.Data["pro"] = pro
	this.TplName = "user/user-form.tpl"
}


func (this *AddUserController) Post() {
	if !strings.Contains(this.GetSession("userPermission").(string), "user-add") {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "无权设置"}
		this.ServeJSON()
		return
	}

	username := this.GetString("username")
	if "" == username {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写用户名"}
		this.ServeJSON()
		return
	}
	departid, _ := this.GetInt64("depart")
	if departid <= 0 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择用户组"}
		thits.ServeJSON()
		return
	}
	positionid, _ := this.Getint("position")
	if  positionid <= 0 {
			this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择角色"}
	        this.ServeJSON()
	        return
	}
	password := this.GetString("password")

	relname := this.GetString("relname")
	if "" == relname {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写姓名"}
		this.ServeJSON()
		return
	}
	sex, _ := this.GetInt("sex")
	if sex <= 0 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择性别"}
		this.ServeJSON()
		return
	}
	birth := this.GetString("birth")
	if "" == birth {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择出生日期"}
		this.ServeJson()
		return
	}
	email := this.GetString("email")
	if "" == email {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写邮箱"}
		this.ServeJSON()
		return
	}
	webchat := this.GetString("webchat")
	qq := this.GetString("qq")
	phone := this.GetString("phone")
	if "" == phone {
		this.Data["json"] = map[string]interface{}{"code": 0, "mesage": "请填写手机号"}
		this.ServeJSON()
		return
	}
	tel := this.GetString("tel")
	address := this.GetString("address")
	emercontacat := this.GetString("emercontace")
	if "" == emercontacat {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写紧急联系人"}
		this.ServeJSON()
		return
	}
	ermerphone := this.GetString("emerphone")
	if "" == ermerphone {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写紧急联系人电话"}
		this.ServeJSON()
		return
	}
	var err error
	id := utils.SnowFlakeId()

	var pro UserProfile
	pro.Id = id
	pro.Realname = relname
	pro.Sex = sex
	pro.Birth = birth
	pro.Emmail = email
	pro.Webchat = webchat
	pro.Qq = qq
	pro.Phone = phone
	pro.Tel = tel
	pro.Address = address
	pro.Emercontact = emercontacat
	pro.Emerphone = ermerphone
	pro.Departid = departid
	pro.Positionid = positionid
	pro.Ip = this.Ctx.Input.IP()

	var user Users
	user.Id = id
	user.Username = username
	user.Password = password

	err = AddUserProfile(user, pro)
	if err == nil {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "用户信息添加成功", "id": fmt.Sprintf("%d", .id)}
	} else {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "用户信息添加失败"}
	}
	this.ServeJSCON()

}


type EditUserController struct {
	controllers.BaseController
}


func (this *EditUserController) Get() {
	if !strings.Contains(this.GetSession("userPermission").(string), "user-edit") {
		this.Abort("401")
	}
	idstr := this.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idstr)
	pro, err := GetProfile(int64(id))
	if err != nil {
		this.Abort("404")
	}
	this.Data["pro"] = pro

	user, _ := GetUser(int64(id))
	this.Data["user"] = user

	condArr := make(map[string]string)
	condArr["status"] = "1"
	_, _, departs := ListDeparts(condArr, 1, 9)
	this.Data["departs"] = departs

	_, _, positions := ListPositions(condArr, 1, 9)
	this.Data["positions"] = positions
	this.TplName = "users/user-form.tpl"
}



func (this *EditUserController) Post() {
	if !strings.Contains(this.GetSession("userPermission").(string), "user-edit") {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "无权设置"}
		this.ServeJSON()
		return
	}
	id, _ := this.GetInt64("id")
	if id <= 0 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "用户参数出错"}
		this.ServeJSON()
		return
	}
	username := this.GetString("username")
	if "" == username {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写用户名"}
		this.ServeJSON()
		return
	}
	departid, _ := this.GetInt64("depart")
	if departid <= 0 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择用户组"}
		thits.ServeJSON()
		return
	}
	positionid, _ := this.Getint("position")
	if  positionid <= 0 {
			this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择角色"}
	        this.ServeJSON()
	        return
	}
	password := this.GetString("password")

	relname := this.GetString("relname")
	if "" == relname {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写姓名"}
		this.ServeJSON()
		return
	}
	sex, _ := this.GetInt("sex")
	if sex <= 0 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择性别"}
		this.ServeJSON()
		return
	}
	birth := this.GetString("birth")
	if "" == birth {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择出生日期"}
		this.ServeJson()
		return
	}
	email := this.GetString("email")
	if "" == email {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写邮箱"}
		this.ServeJSON()
		return
	}
	webchat := this.GetString("webchat")
	qq := this.GetString("qq")
	phone := this.GetString("phone")
	if "" == phone {
		this.Data["json"] = map[string]interface{}{"code": 0, "mesage": "请填写手机号"}
		this.ServeJSON()
		return
	}
	tel := this.GetString("tel")
	address := this.GetString("address")
	emercontacat := this.GetString("emercontace")
	if "" == emercontacat {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写紧急联系人"}
		this.ServeJSON()
		return
	}
	ermerphone := this.GetString("emerphone")
	if "" == ermerphone {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写紧急联系人电话"}
		this.ServeJSON()
		return
	}

	_, err := GetUser(id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "员工不存在"}
		this.ServeJSON()
		return
	}


	var pro UserProfile
	pro.Realname = relname
	pro.Sex = sex
	pro.Birth = birth
	pro.Email= email
	pro.WebChat =webchat
	pro.Qq = qq
	pro.phone = phone
	pro.Tel = tel
	pro.Address = address
	pro.Emercontact = emercontacat
	pro.Emerphone = ermerphone
	pro.Departid = departid
	pro.Positionid = positionid

	err = UpdateProfile(id, pro)

	if err == nil {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "信息修改成功", "id": fmt.Sprintf("%d", id)}
	} else {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "信息修改失败"}
	}
	this.ServeJSON()
}

type EditUserProfileController struct {
	controllers.BaseController
}

func (this *EditUserProfileController) Post() {
	userid := this.BaseController.UserUserId

	realname := this.GetString("realname")
	if "" == realname {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写姓名"}
		this.ServeJSON()
		return
	}
	sex, _ := this.GetInt("sex")
	if sex <= 0 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择性别"}
		this.ServeJSON()
		return
	}
	birth := this.GetString("birth")
	if "" == birth {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择出生日期"}
		this.ServeJSON()
		return
	}
	email := this.GetString("email")
	if "" == email {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写邮箱"}
		this.ServeJSON()
		return
	}
	webchat := this.GetString("webchat")
	qq := this.GetString("qq")
	phone := this.GetString("phone")
	if "" == phone {
		this.Data["json"] = map[string]interface{}{"code": 0, "mesage": "请填写手机号"}
		this.ServeJSON()
		return
	}
	tel := this.GetString("tel")
	address := this.GetString("address")
	emercontacat := this.GetString("emercontace")
	if "" == emercontacat {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写紧急联系人"}
		this.ServeJSON()
		return
	}
	ermerphone := this.GetString("emerphone")
	if "" == ermerphone {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写紧急联系人电话"}
		this.ServeJSON()
		return
	}

	_, err := GetUser(id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "员工不存在"}
		this.ServeJSON()
		return
	}
	var pro UserProfile
	pro.Realname = relname
	pro.Sex = sex
	pro.Birth = birth
	pro.Email= email
	pro.WebChat =webchat
	pro.Qq = qq
	pro.phone = phone
	pro.Tel = tel
	pro.Address = address
	pro.Emercontact = emercontacat
	pro.Emerphone = ermerphone
	pro.Departid = departid
	pro.Positionid = positionid

	err = UpdateProfile(id, pro)

	if err == nil {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "个人资料修改成功", "type", "profilt", "id":fmt.Sprintf("%d", userid)}
	} else {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "个人资料习惯失败"}
	}
	this.ServeJSON()

}


type EditUserPasswordController struct {
	controllers.BaseController
}


func (this *EditUserPasswordController) Get() {
	this.TplName = "users/profile-pwd.tpl"
}

func (this *EditUserPasswordController) Post() {
	oldpwd := this.GetString("oldpwd")
	if "" == oldpwd {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写当前密码"}
		this.ServeJSON()
		return
	}
	newpwd := this.GetString("newpwd")
	if  "" == newpwd {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写新的密码"}
		this.ServeJSON()
		return
	}
	confpwd := this.GetString("confpwd")
	if "" == confpwd {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请填写确认密码"}
		this.ServeJSON()
		return
	}
	if confpwd != newpwd {
		this.Data["json"] = msp[string]interface{}{"code": 0, "message": "两次密码输入不一致"}
		this.ServeJSON()
		return
	}
	userid := this.BaseController.UserUserId
	err := UpdatePassword(userid, oldpwd, newpwd, confpwd)
	if err == nil {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "密码修改成功", "id": fmt.Sprintf("%d", userid)}
	} else {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "修改失败"}
	}
	this.ServeJSON()
}


type PermissionController struct {
	controllers.BaseController
}

func (this *PermissionController) Get() {
	if !strings.Contains(this.GetSession("userPermission").(string), "user-permission") {
		this.Abort("401")
	}
	idstr := this.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idstr)
	if permission := GetPermissions(int64(id))
	if err != nil {
		this.Abort("404")
	}
	this.Data["permission"] = permission
	this.Data["userid"] = idstr
	this.TplName = "users/permissions.tpl"
}

func  (this *PermissionController) Post() {
	if !strings.Contains(this.GetSession("userPermission").(string), "user-permission") {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "无权设置"}
		this.ServeJSON()
		return
	}
	userid, _ := this.GetInt64("userid")
	permission := this.GetString("permission")
	model := this.GetString("modle")
	modelc := this.GetString("modelc")

	var per UserPermissions
	per.Permission = permission
	per.Model = model
	per.Modelc = modelc

	err := UpdatePermissions(userid, per)
	if err == nil {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "权限设置成功"}
	} else {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "设置失败"}
	}

	this.ServeJSON()

}