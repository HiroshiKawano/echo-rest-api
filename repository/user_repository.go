package repository

import (
	"go-rest-api/model"

	"gorm.io/gorm"
)

// Repositoryインターフェイス
// GO言語では、インターフェースはメソッドの一覧になっていて、
// GetUserByEmailとCreateUserという2つのメソッドを定義
type IUserRepository interface {
	// GetUserByEmailは、第1引数でユーザーオブジェクトのポインタ、
	// 第2引数で検索したいメールstringで受け取る
	// 返り値の型はerrorインターフェイス型
	GetUserByEmail(user *model.User, email string) error
	// CreateUserも、ユーザーオブジェクトのポインタを引数で受け取り、返り値の型はerrorインターフェース型
	CreateUser(user *model.User) error
}

// userRepository構造体の定義
type userRepository struct {
	// 構造体の要素として、gorm.DBのPOINTER型でDBという名前のフィールドを作る
	db *gorm.DB
}

// リポジトリにDBのインスタンスを依存注入するために、リポジトリーの方にもコンストラクターを作る
// NewUserRepositoryというコンストラクターを作り、外側でインスタンス化されたDBを引数で受け取る
// 返り値の型は、IUserRepositoryのinterface型にする
func NewUserRepository(db *gorm.DB) IUserRepository {
	// db(DBのインスタンス)を要素にして、userRepository構造体の実体を作成し、ポインターをリターンで返す
	return &userRepository{db}
}

// userRepository型をpointerレシーバーとして受け取る形でGetUserByEmailメソッドを定義
// 引数の型と返り値の型はinterfaceで定義している内容と全く同じにする必要がある
func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {
	// ur.db.Whereでデータベースの中でメールの値が引数で受け取った値に一致するユーザーを探す
	// そのユーザーが存在する場合は、引数で受け取ったユーザーオブジェクト(user)のアドレスが指し示す先の値の内容を
	//検索したユーザーオブジェクトの内容で書き換える
	if err := ur.db.Where("email=?", email).First(user).Error; err != nil {
		// エラーが発生した場合は、リターンでエラーを返し
		return err
	}
	// 成功した場合はNILを返す
	return nil
}

// userRepository型をPOINTERレシーバーとして受け取る形でCreateUserメソッドを作る
// 引数の型と返り値の型はinterfaceで定義している内容と全く同じにする必要がある
func (ur *userRepository) CreateUser(user *model.User) error {
	// ur.db.Createで引数で受け取っていたユーザーオブジェクト(user)のポインタを渡す
	// ユーザーの作成に成功した場合は、このポインターが指し示す先の値が新しく作成されたユーザーの情報で書き換えられる
	if err := ur.db.Create(user).Error; err != nil {
		// エラーが発生した場合は、リターンでエラー
		return err
	}
	// 成功した場合はNILを返す
	return nil
}
