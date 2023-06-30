package controller

import (
	"go-rest-api/model"
	"go-rest-api/usecase"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type ITaskController interface {
	// ここでは5つメソッドを定義しまして、引数はECHOフレームワークのContext、
	// 返り値としてerrorインターフェース型を設定
	GetAllTasks(c echo.Context) error
	GetTaskById(c echo.Context) error
	CreateTask(c echo.Context) error
	UpdateTask(c echo.Context) error
	DeleteTask(c echo.Context) error
}

type taskController struct {
	// tuというフィールド名でusecaseパッケージ内のITaskUsecaseインターフェースの値を格納できるようにしておきます。
	tu usecase.ITaskUsecase
}

func NewTaskController(tu usecase.ITaskUsecase) ITaskController {
	// NewTaskControllerのコンストラクタは引数で外側でインスタンス化されるタスクのusecaseを受け取りまして、
	// そして受け取った値を使ってタtaskController構造体のインスタンスを作成して
	// そのポインタ(&taskController)をreturnで返す
	return &taskController{tu}
}

func (tc *taskController) GetAllTasks(c echo.Context) error {
	// まずユーザーから送られてくるJWTtokenに組み込まれているユーザーIDの値を取り出す
	// 後ほどrouterの方に実装するJWTのミドルウェア側で送られてきたJWTtokenをデコードしてくれまして、
	// デコードした内容をechoのContextの中にuserというフィールド名を付けて自動的に格納してくれます。
	// そして、コントローラー側では、そのuserというキーワードを使ってContextからJWTをデコードした値を読み込んできます。
	user := c.Get("user").(*jwt.Token)
	// その中にはデコードされたClaimsが格納されてますので、user.Claimsで取り出して、
	claims := user.Claims.(jwt.MapClaims)
	// そしてclaimsの中にあるユーザーIDを取得してユuserIdという変数に代入するようにしています。
	userId := claims["user_id"]

	// Contextから取得した値(userId)はany型になっていますので、
	// いったんfloat64に型アサーションしてからuint型に型変換するようにしています。
	// そして、タスクユースケースのGetAllTasksメソッドにuserIdを引数として渡すようにしています。
	tasksRes, err := tc.tu.GetAllTasks(uint(userId.(float64)))
	if err != nil {
		// エラーが発生した場合は、コンテキスト.JSONでクライアントにInternalServerErrorのステータスとエラーメッセージを返す
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// 成功した場合は、コンテキスト.JSONでStatusOKと取得したタスクの一覧をレスポンス(tasksRes)で返す
	return c.JSON(http.StatusOK, tasksRes)
}

func (tc *taskController) GetTaskById(c echo.Context) error {
	// コンテキストからuserキーワードを使ってJWTtokenをデコードした内容を取得
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	// その中からuser_idの値を取得して、userIdという変数に代入
	userId := claims["user_id"]
	// さらに、リクエストパラメーターからtaskIdを取得しまして、
	id := c.Param("taskId")
	// こちらはstring型になっていますので、Atoiを使ってstring型からint型に変換します。
	taskId, _ := strconv.Atoi(id)
	// タスクユースケースのGetTaskByIdメソッドを呼び出して、第1引数にuserId第2引数にtaskIdをお渡していきます。
	taskRes, err := tc.tu.GetTaskById(uint(userId.(float64)), uint(taskId))
	if err != nil {
		// エラーが発生した場合は、StatusInternalServerError、
		// 成功した場合はStatusOKで取得したタスク(taskRes)をクライアントの方にJSONで返す
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, taskRes)
}

func (tc *taskController) CreateTask(c echo.Context) error {
	// コンテキストの中からuser_idを取得
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	// 0値でタスク構造体を作り、コンテキストBindを使うことでリクエストボディーに含まれる内容をタスク構造体に代入
	task := model.Task{}
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// さらに、taskオブジェクトのUserIdのフィールドにコンテキストから取得したuserIdの値を格納します。
	task.UserId = uint(userId.(float64))
	// そして、そのtaskオブジェクトをタスクのCreateTaskに引数として渡していきます。
	taskRes, err := tc.tu.CreateTask(task)
	if err != nil {
		// 失敗した場合は、StatusInternalServerError、
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// 成功した場合はコンテキストJSONでStatusCreatedという
	// 新規作成されたタスクのオブジェクトをクライアントに返す
	return c.JSON(http.StatusCreated, taskRes)
}

func (tc *taskController) UpdateTask(c echo.Context) error {
	// コンテキストからuser_idを取得するのとリクエストパラメーターからtaskIdを取得しておきます。
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]
	id := c.Param("taskId")
	// stringからint型に変換しておきます。
	taskId, _ := strconv.Atoi(id)

	task := model.Task{}
	// コンテキストBindを使ってリクエストオブジェクトの値をタスクオブジェクトにバインド
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// その後にタスクユースケースのUpdateTaskを第1引数をtask第2引数をuserId第3引数をtaskIdとして呼び出す
	taskRes, err := tc.tu.UpdateTask(task, uint(userId.(float64)), uint(taskId))
	if err != nil {
		// エラーが発生した場合は、StatusInternalServerError
		// 成功した場合は、更新後のタスクの値をStatusOKでクライアントにJSONで返すようにしています。
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, taskRes)
}

func (tc *taskController) DeleteTask(c echo.Context) error {
	// コンテキストからuser_idを取得し、リクエストパラメーターからtaskIdを取得しています。
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]
	id := c.Param("taskId")
	taskId, _ := strconv.Atoi(id)

	// そして、タスクユースケースのDeleteTaskを呼び出して、userIdとtaskIdを渡していきます。
	err := tc.tu.DeleteTask(uint(userId.(float64)), uint(taskId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// 成功した場合は、コンテキストのNoContentでStatusNoContentをクライアントに返すようにしておきます。
	return c.NoContent(http.StatusNoContent)
}
