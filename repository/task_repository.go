package repository

import (
	"fmt"
	"go-rest-api/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ITaskRepository interface {
	// GetAllTasksはログインしているユーザー自身が作成したタスクの一覧を取得するメソッド
	// タスクの一覧を配列に格納するために第1引数としてモデルタスクのスライス([]model.Task)のポインタを渡す
	// 第2引数はログインしてるユーザーのuserIdを渡す
	// 返り値はerrorインターフェース型
	GetAllTasks(tasks *[]model.Task, userId uint) error
	// GetTaskByIdは引数で渡すtaskIdに一致するタスクを取得するメソッド
	GetTaskById(task *model.Task, userId uint, taskId uint) error
	// CreateTaskでタスクの新規作成
	CreateTask(task *model.Task) error
	// UpdateTaskで引数で渡すtaskIdのタスクの内容の更新
	UpdateTask(task *model.Task, userId uint, taskId uint) error
	// DeleteTaskで引数で渡すtaskIdのタスクのオブジェクトの削除
	DeleteTask(userId uint, taskId uint) error
}

// taskRepositoryという構造体を定義
type taskRepository struct {
	//フィールドとしてDBを作っておきます。
	db *gorm.DB
}

// NewTaskRepositoryという名前のコンストラクターも作っておきます。
func NewTaskRepository(db *gorm.DB) ITaskRepository {
	// 外側でインスタンス化されたdbを引数で受け取り、taskRepository構造体の実体を作成
	// その実体のアドレスを取得してリターンで返す
	return &taskRepository{db}
}

// GetAllTasksの実装
// taskRepositoryをpointerレシーバーとして受け取る形でGetAllTasksというメソッドを定義
// 引数と返り値の型は、interfaceの型と一緒にする必要がある
func (tr *taskRepository) GetAllTasks(tasks *[]model.Task, userId uint) error {
	// タスクの一覧の中でユーザーIDのフィールド(user_id)が引数で渡されたユーザーID(userId)に一致するタスクの一覧を取得
	// Order("created_at")でタスクの作成日時が一番新しいものが末尾に来る順番でデータを取得する
	if err := tr.db.Joins("User").Where("user_id=?", userId).Order("created_at").Find(tasks).Error; err != nil {
		// エラーが発生した場合はエラーを返し、
		return err
	}
	// 成功した場合はNILをリターンで返す
	return nil
}

// タスクの一覧の中でユーザーIDの値(user_id)が引数で受け取るユーザーID(userId)に一致するタスクの一覧を抽出
// さらに、その中でタスクの主キーが引数で受け取ったタスクID(taskId)に一致するtaskを取得
// そして、取得したタスクオブジェクト(task)を引数で受け取っていたポインタアドレスが指し示す先(*model.Task)のメモリー領域に書き込む
func (tr *taskRepository) GetTaskById(task *model.Task, userId uint, taskId uint) error {
	if err := tr.db.Joins("User").Where("user_id=?", userId).First(task, taskId).Error; err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) CreateTask(task *model.Task) error {
	// tr.db.Createでtaskのポインタを引数で渡す
	if err := tr.db.Create(task).Error; err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) UpdateTask(task *model.Task, userId uint, taskId uint) error {
	// tr.db.Modelでtaskオブジェクトのポインターを渡す
	// そして、Clauses(clause.Returning{})のキーワードをつけると
	// 更新した後のタスクのオブジェクトをこのタスクのポインタが指し示す先(*model.Task)に書き込んでくれるようになります。
	// そして、Whereでタスクの主キーであるID(id)が引数で受け取れるタスクID(taskId)に一致する
	// かつユーザーIDが引数で受け取るユーザーID(user_id)に一致するタスクに対してUpdateの処理をかけていきます。
	// そして、ここではtitleの値(title)を引き継いで受け取れるタスクオブジェクトのタイトルの値(task.Title)で更新するようにしています。
	result := tr.db.Model(task).Clauses(clause.Returning{}).Where("id=? AND user_id=?", taskId, userId).Update("title", task.Title)
	// 処理の返り値をresultという変数に代入して、result.Errorでエラーを取得
	if result.Error != nil {
		// エラーが発生した場合は、エラーをリターンで返す
		return result.Error
	}
	// 実際に更新されたレコードの数を取得することができ、
	// その数が1より小さい0の場合は更新が行なわれなかったことを意味してる
	if result.RowsAffected < 1 {
		// その場合はfmt.Errorfでobject does not existとエラーメッセージを付けてエラーをリターンで返す
		return fmt.Errorf("object does not exist")
	}
	return nil
}

func (tr *taskRepository) DeleteTask(userId uint, taskId uint) error {
	// tr.db.Whereで引数で渡されたタスクID(taskId)とユーザーID(userId)に一致するタスクをDELETE
	result := tr.db.Where("id=? AND user_id=?", taskId, userId).Delete(&model.Task{})
	if result.Error != nil {
		// エラーが発生した場合は、エラーをリターンで返す
		return result.Error
	}
	// RowsAffectedが0の場合もエラーを返す
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}
