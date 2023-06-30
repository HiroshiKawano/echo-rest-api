package validator

import (
	"go-rest-api/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ITaskValidator interface {
	// ITaskValidatorという名前のインターフェースを定義し、TaskValidateメソッドを定義。
	TaskValidate(task model.Task) error
}

type taskValidator struct{}

// taskValidator構造体のインスタンスを生成する為のコンストラクターを
// NewTaskValidatorという名前で作成
func NewTaskValidator() ITaskValidator {
	// taskValidator構造体を作成した後に&でアドレスを取得してリターンで返す
	return &taskValidator{}
}

// taskValidatorをポインターレシーバー(*)として受け取る形でTaskValidateメソッドを定義しまして、
// 引数でバリデーションで評価したいtaskのオブジェクトを受け取れるようにしています。
func (tv *taskValidator) TaskValidate(task model.Task) error {
	// そして、validationパッケージで定義されてるValidateStruct関数を実行しまして、
	// 第1引数にタスクオブジェクトのアドレス(&task)を渡します。
	// そして第2引数でtaskのTitleに対するvalidationを実装しています。
	// validation.RequiredでTitleに値が存在するかチェックすることができます。
	// 値が存在しない場合は、title is requiredというエラーメッセージを返す
	// そして、2つ目のvalidationとして、文字数が最小1最大10になっているかチェック
	return validation.ValidateStruct(&task,
		validation.Field(
			&task.Title,
			validation.Required.Error("title is required"),
			validation.RuneLength(1, 10).Error("limited max 10 char"),
		),
	)
}
