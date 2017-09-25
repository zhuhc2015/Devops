package routers

import (
	"github.com/Devops/opms/controllers"
	//"github.com/Devops/opms/controllers/projects"
	"github.com/Devops/opms/controllers/users"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	//用户
	beego.Router("/user/manage", &users.ManageUserController{})
	beego.Router("/user/ajax/status", &users.AjaxStatusUserController{})
	beego.Router("/user/edit/:id", &users.EditUserController{})
	beego.Router("/user/add", &users.AddUserController{})
	beego.Router("/user/avatar", &users.AvatarUserController{})
	beego.Router("/user/ajax/search", &users.AjaxSearchUserController{}) //搜索用户名匹配
	beego.Router("/user/show/:id", &users.ShowUserController{})
	beego.Router("/my/manage", &users.ShowUserController{})
	beego.Router("/user/profile", &users.EditUserProfileController{})
	beego.Router("/user/password", &users.EditUserPasswordController{})

	beego.Router("/user/permission/:id", &users.PermissionController{})

	beego.Router("/login", &users.LoginUserController{})
	beego.Router("/logout", &users.LogoutUserController{})
}
