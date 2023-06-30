package main

import (
	"go-rest-api/controller"
	"go-rest-api/db"
	"go-rest-api/repository"
	"go-rest-api/router"
	"go-rest-api/usecase"
	"go-rest-api/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// データベースをインスタンス化
	// データベースパッケージの中で作っておいたNewDBを実行して
	// 作成されたインスタンスをdbという変数に格納
	db := db.NewDB()
	// NewUserValidatorとNewTaskValidatorでコンストラクターを実行して構造体のインスタンスを作成します。
	// インスタンス(userValidator、taskValidator)をusecaseのコンストラクターに渡していきます。
	userValidator := validator.NewUserValidator()
	taskValidator := validator.NewTaskValidator()
	// レポジトリで作っておいたコンストラクターを起動
	// repositoryパッケージの中で作っておいたNewUserRepositoryコンストラクターを起動
	// 外側でインスタンス化してるデーターベース(db)を引数として注入
	userRepository := repository.NewUserRepository(db)
	// taskRepositoryのコンストラクターを起動
	taskRepository := repository.NewTaskRepository(db)
	// usecaseのコンストラクターも起動
	// usecaseのパッケージで作っておいたNewUserUsecaseコンストラクターを起動
	// 引数として外側でインスタンス化しておいたuserRepositoryを引数として注入
	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	// taskUsecaseのコンストラクターのNewTaskUsecaseも起動
	taskUsecase := usecase.NewTaskUsecase(taskRepository, taskValidator)
	// controllerのコンストラクターも起動
	// controllerパッケージの中で作っておいたNewUserControllerコンストラクターを起動
	// 外側でインスタンス化してるuserUsecaseのインスタンスを引数として注入
	userController := controller.NewUserController(userUsecase)
	// NewTaskControllerを使ってtaskControllerのコンストラクターも起動
	taskController := controller.NewTaskController(taskUsecase)
	// routerパッケージの中に作っておいたNewRouter関数を呼び出す
	// 外側でインスタンス化してるuserControllerを引数として注入
	// taskControllerをNewRouterの第2引数に追加
	e := router.NewRouter(userController, taskController)
	// echoのインスタンス(e)を使ってサーバーを起動
	// e.Startでサーバーを起動し、port番号を8080番にして、
	// エラーが発生した場合は、e.Loggerの機能を使ってログ情報出力した後にプログラムを強制終了

	// CORSミドルウェアの設定を追加(CreateReactApp→Viteに変更の場合)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowCredentials: true,
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders:     []string{"X-Requested-With", "Content-Type", "X-CSRF-Token"},
	}))

	e.Logger.Fatal(e.Start(":8080"))
}
