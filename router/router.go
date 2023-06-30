package router

import (
	"go-rest-api/controller"
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// routerの中でUserControllerを使えるようにするために、
// 引数でユーザーコントローラー(uc)を受け取れるようにしておきます。
// routerの中でタスクコントローラーを使用できるようにするために、
// 引数のところにタスクコントローラーを追加しておきます。
func NewRouter(uc controller.IUserController, tc controller.ITaskController) *echo.Echo {
	// echo.Newでエコーのインスタンスを作成
	e := echo.New()
	// e.Useで、CORSのmiddlewareを追加しまして、新ORIGINSのところにアクセスをですね。
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// AllowOriginsのところにアクセスを許可するフロントエンドのドメインを追加
		// ここでは、Reactのlocalhost:3000番と環境変数のFE_URLの値をAllowOriginsに追加しています。
		// こちらは後ほどフロントエンドのアプリケーションをVercelにデプロイした時に
		// 取得できるドメインを環境変数FE_URLに設定していきます。
		// AllowHeadersで許可するヘッダーの一覧を入力していきます。
		// ここではechoのHeaderXCSRFTokenを含めることによって、header経由でCsrfTokenを受け取れるようにしています。
		// そして許可をしたいメソッド("GET", "PUT", "POST", "DELETE")を追加
		// クッキーの送受信を可能にするために、AllowCredentialsをtrueに設定
		AllowOrigins: []string{"http://localhost:3000", os.Getenv("FE_URL")},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowCredentials: true,
	}))
	// e.Useで、CSRFのmiddlewareを設定しています。
	// CsrfTokenを格納するcookieの設定を行いまして、
	// まずCookiePathとしてINDEX、CookieDomainとしてAPI_DOMAIN、CookieHTTPOnlyをtrueにする
	// Postmanで動作確認する時だけは、SameSiteDefaultModeにしておきます。
	// CsrfTokenの有効期限は、デフォルトでは24時間に設定されてるんですけども、
	// この値を変更したい場合は、CookieMaxAgeを設定することで有効期限を秒単位で設定することができます。
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",
		CookieDomain:   os.Getenv("API_DOMAIN"),
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteNoneMode,
		//CookieSameSite: http.SameSiteDefaultMode,
		//CookieMaxAge:   60,
	}))
	// echoのインスタンスに対してエンドポイントを追加
	// POSTメソッドでsignupのエンドポイントにリクエストがあった際は、
	// ユーザーコントローラー(uc)のSignUpメソッドを呼び出す
	e.POST("/signup", uc.SignUp)
	// POSTメソッドでLOGINのエンドポイントにリクエストがあった場合は、
	// ユーザーコントローラー(uc)のLogOutメソッドを呼び出す
	e.POST("/login", uc.LogIn)
	// エンドポイントにPOSTメソッドでリクエストがあった際は、
	// ユーザーコントローラー(uc)のLogOutメソッドを呼び出す
	e.POST("/logout", uc.LogOut)
	// e.GETでCSRFのエンドポイント(/csrf)にリクエストがあった際は
	//ユーザーコントローラー(uc)のCsrfTokenのメソッドを呼び出すようにしておきます。
	e.GET("/csrf", uc.CsrfToken)
	// ECHOインスタンスのeに対して新しくグループを作っていきます。
	//タスク関係のエンドポイントをグループ化して、こちらをtという変数に格納しておきます。
	t := e.Group("/tasks")
	// タスクのグループに対して、JWTのミドルウェアを適用するようにしておきます。
	// Useキーワードを使うことで、エンドポイントにミドルウェアを追加することができます。
	// ここではECHOのJWTというミドルウェアを適用し、SigningKeyのところにJWTを生成した時と同じSECRETキーを指定します。
	// TokenLookupのところでは、クライアントから送られてくるJWTtokenがどこに格納されてるかのを指定する必要があります。
	// cookieの中にtokenという名前でJWTtokenを格納するように実装してるので、
	// TokenLookupのところでcookie:tokenを指定します。
	t.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
	}))
	// タスク関係のエンドポイントを追加
	// GetAllTasksのエンドポイントにリクエストがあった際は、
	// タスクコントローラーのGetAllTasksを呼び出すようにしています。
	// パラメーター付きでリクエストがあった際は、タスクコントローラーのGetTaskById
	// POSTメソッドのリクエストの場合はCreateTask、
	// PUTメソッドの場合はUpdateTask
	// DELETEの場合はタスクコントローラーのDeleteTaskを呼び出すようにしておきます。
	t.GET("", tc.GetAllTasks)
	t.GET("/:taskId", tc.GetTaskById)
	t.POST("", tc.CreateTask)
	t.PUT("/:taskId", tc.UpdateTask)
	t.DELETE("/:taskId", tc.DeleteTask)
	//NewRouter関数の返り値としてechoインスタンス(e)を返す
	return e
}
