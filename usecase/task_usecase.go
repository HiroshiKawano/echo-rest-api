package usecase

import (
	"go-rest-api/model"
	"go-rest-api/repository"
	"go-rest-api/validator"
)

type ITaskUsecase interface {
	GetAllTasks(userId uint) ([]model.TaskResponse, error)
	GetTaskById(userId uint, taskId uint) (model.TaskResponse, error)
	CreateTask(task model.Task) (model.TaskResponse, error)
	UpdateTask(task model.Task, userId uint, taskId uint) (model.TaskResponse, error)
	DeleteTask(userId uint, taskId uint) error
}

type taskUsecase struct {
	// taskUsecase構造体はtrというフィールド名で、repositoryパッケージ内のITaskRepositoryインターフェースの値を格納
	tr repository.ITaskRepository
	// taskUsecase構造体のフィールドにITaskValidatorのtvというフィールドを追加
	tv validator.ITaskValidator
}

// NewTaskUsecaseのコンストラクターは引数で外側でインスタンス化されるタスクリポジトリー(tr)を受け取り、
// その値を使ってtaskUsecase構造体の実体を生成
// NewTaskUsecaseのコンストラクターに外側でインスタンス化されるITaskValidatorを注入できるように
// するために引数のところにtv validator.ITaskValidatorを追加します。
func NewTaskUsecase(tr repository.ITaskRepository, tv validator.ITaskValidator) ITaskUsecase {
	// &でアドレスを取得してリターンで返す
	// そしてタスクユースケースをインスタンス化するフィールドのところにtvを追加
	return &taskUsecase{tr, tv}
}

// GetAllTasksは、引数でユーザーID(userId)を受け取り、
// 返り値の1つ目の型として、modelパッケージで定義したTaskResponse構造体の配列の型を指定
// そして、2つ目の返り値の型はerrorインターフェース型
func (tu *taskUsecase) GetAllTasks(userId uint) ([]model.TaskResponse, error) {
	// 取得するタスク一覧を格納するためのTask構造体のスライスを定義
	tasks := []model.Task{}
	//taskリポジトリのGetAllTasksを呼び出しtasksのアドレスとuserIdを引数で渡す
	if err := tu.tr.GetAllTasks(&tasks, userId); err != nil {
		// エラーが返ってきた場合は、1つ目の返り値としてnilスライス、2つ目の返り値としてエラーを返す
		return nil, err
	}
	// 取得に成功した場合は、クライアントへのレスポンス用のTaskResponse構造体を0値で作成
	resTasks := []model.TaskResponse{}
	// for rangeでタtasksからはタスクを一つ一つ取り出し、タTaskResponse構造体を新しく作る
	for _, v := range tasks {
		t := model.TaskResponse{
			ID:        v.ID,
			Title:     v.Title,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
		// 作成した新しい構造体をresTasksのスライスにappendで追加
		resTasks = append(resTasks, t)
	}
	// 最後にreturnでresTasksとnilを返す
	return resTasks, nil
}

func (tu *taskUsecase) GetTaskById(userId uint, taskId uint) (model.TaskResponse, error) {
	// 取得するTaskを格納するための構造体をまずは作成し
	task := model.Task{}
	// tu.tr.GetTaskByIdで、この空の構造体のポインタ(&task)を第1引数として渡していきます。
	// そして、userIdとtaskIdを引数で渡していきます。
	if err := tu.tr.GetTaskById(&task, userId, taskId); err != nil {
		// エラーが発生した場合は、TaskResponse構造体を0値でインスタンス化したものとerrをreturnで返す
		return model.TaskResponse{}, err
	}
	// 成功した場合は、第1引数で渡したポインタ(&task)が指し示す先の値が取得したタスクの値で書き換えられますので、
	// ID,Title,CreatedAt,UpdatedAtの値を取り出して
	// 新しくタスクレスポンス構造体の実体(model.TaskResponse)を作成してreturnで返す
	resTask := model.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
	// 作成した構造体とnilをreturnで返す
	return resTask, nil
}

// タスクリポジトリのCreateTaskを呼び出す前にタスクのバリデーションを実行
func (tu *taskUsecase) CreateTask(task model.Task) (model.TaskResponse, error) {
	// tu.tv.TaskValidateで引数としてバリデーションを行いたいtaskのオブジェクトを渡します。
	if err := tu.tv.TaskValidate(task); err != nil {
		// そして、バリデーションに失敗した場合は、returnでエラーを返す
		return model.TaskResponse{}, err
	}
	// taskリポジトリ内のCreateTaskを呼び出し、引数としてtaskオブジェクトのアドレスを渡す
	if err := tu.tr.CreateTask(&task); err != nil {
		// CreateTaskでエラーが発生した場合は、TaskResponse構造体の0値の実体とerrをreturnで返す
		return model.TaskResponse{}, err
	}
	// 成功した場合は、引数で渡したアドレスが指し示す先の値が新規作成したタスクの値で書き換わっていますので
	// ID,Title,CreatedAt,UpdatedAtの値を取り出して
	// 新しくタスクレスポンス構造体の実体(model.TaskResponse)を作成してreturnで返す
	resTask := model.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
	// 成功した場合は、第2引数のエラーはnilを返す
	return resTask, nil
}

func (tu *taskUsecase) UpdateTask(task model.Task, userId uint, taskId uint) (model.TaskResponse, error) {
	// tu.tv.TaskValidateでバリデーションを掛けたいtaskオブジェクトを引数で渡しておきます。
	if err := tu.tv.TaskValidate(task); err != nil {
		return model.TaskResponse{}, err
	}
	// tu.tr.UpdateTaskでtaskオブジェクトのアドレス,userId,taskIdを渡していきます。
	if err := tu.tr.UpdateTask(&task, userId, taskId); err != nil {
		// エラーが発生した場合は、TaskResponseの0値のインスタンスとerrをreturnで返す
		return model.TaskResponse{}, err
	}
	// 成功した場合は、第1引数で渡したtaskのアドレスが指し示す先のメモリ領域のタスクの値が更新後のタスクで書きかえられていますので、
	// ID,Title,CreatedAt,UpdatedAtの値を取り出して
	// 新しくタスクレスポンス構造体の実体(model.TaskResponse)を作成してreturnで返す
	resTask := model.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
	// その作成した構造体をreturnで返すのとerrの値としてnilを返す
	return resTask, nil
}

func (tu *taskUsecase) DeleteTask(userId uint, taskId uint) error {
	// tu.tr.DeleteTaskでuserIdとtaskIdを渡していきます。
	if err := tu.tr.DeleteTask(userId, taskId); err != nil {
		return err
	}
	// そして成功した場合は、returnでnilを返す
	return nil
}
