package usecase

import (
	"go-rest-api/model"
	"go-rest-api/repository"
	"go-rest-api/validator"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Usecaseインターフェイス
// SignUpとLoginというメソッドを定義
type IUserUsecase interface {
	// それぞれ引数でUserオブジェクトを値渡しで受け取り、
	// SignUpの返り値の1つ目はmodelで定義したUserResponse型を割り当て、2つ目はerrorインターフェイス型
	SignUp(user model.User) (model.UserResponse, error)
	// Loginの1つ目の返り値は、JWTtokenを返すためにstring型を割り当て、2つ目はerrorインターフェイス型
	Login(user model.User) (string, error)
}

// Usecase構造体
type userUsecase struct {
	// urフィールドとしてrepositoryを追加
	// Usecaseのソースコードはrepositoryのインターフェイスにだけに依存させるので、
	// repositoryパッケージの中で定義されているIUserRepositoryの型を使用
	ur repository.IUserRepository
	// userUsecaseの構造体のフィールドにvalidator.IUserValidatorのuvというフィールドを追加
	uv validator.IUserValidator
}

// Usecaseにrepositoryを依存性注入するためのコンストラクタ
// 外部でインスタンス化されるRepositoryを引数で受け取れるようにする
// UsecaseのソースコードはRepositoryのインターフェースだけに依存しますので、
// 引数の型としてパッケージに定義されてるIUserRepositoryのインターフェースのですね型を指定
// 返り値の型は、IUserUsecaseのインターフェース型
// userUsecaseのコンストラクターの引数のところにもユーザーバリデーター(uv)を追加します。
func NewUserUsecase(ur repository.IUserRepository, uv validator.IUserValidator) IUserUsecase {
	// 引数で受け取れるRepositoryのインスタンスをフィールドとして、
	// userUsecaseの構造体の実体を作成します。
	// 作成した実体のポインタを&で取得してリターンで返す
	// 返り値で返しているuserUsecase型がIUserUsecaseのインターフェースを満たすためには、
	// interfaceで定義されているすべてのメソッドを実装する必要があります。
	// UserUsecaseのインスタンスを作成するところで、ここで受け取ったuvの値を使用する
	return &userUsecase{ur, uv}
}

// SignUpの実装
/* ユーザーユースケースをポインターレシーバーとしてSignUpというメソッドを追加しています。
*  引数の方と返り値の方はINTERFACEで定義されてるものと同じ
*
 */
func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	// 入力値のバリデーターをで追加
	// uu.uv.UserValidateでuserオブジェクトのバリデーションを掛けるようにしておきます。
	if err := uu.uv.UserValidate(user); err != nil {
		// バリデーションに失敗した場合はreturnでエラーを返すようにしておきます。
		return model.UserResponse{}, err
	}
	// パスワードのハッシュ化
	// user.Passwordでユーザーが入力したパスワードを取り出し
	// 10は暗号化の複雑さ
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	// エラーが発生した場合は、0値でUSERRESPONSE構造体のインスタンスを作成したものと発生したエラーを返す
	if err != nil {
		return model.UserResponse{}, err
	}
	// 新規ユーザの作成
	// 引数で受け取ったユーザーオブジェクトのメール情報とハッシュ化したパスワードを使って
	// 新しくユーザーオブジェクトを作成しnewUser変数に格納
	newUser := model.User{Email: user.Email, Password: string(hash)}
	// リポジトリー層のCreateUserをユースケースの方から呼び出し、
	//フィールド名uu（POINTERレシーバーとして受け取るユーザーユースケース）
	// フィールド名ur(ユーザーユースケース構造体の中ではユーザーレポジトリ)
	// リポジトリー内で定義されているCreateUserのメソッドを呼び出す
	// CreateUserはユーザーオブジェクトのポインタを引数で受け取りますので、
	// &newUserでnewUserのポインタを取得し引数で渡す
	// エラーの場合は、ユーザーレスポンス0値のインスタンスとエラーを返す
	if err := uu.ur.CreateUser(&newUser); err != nil {
		return model.UserResponse{}, err
	}
	// 作成された新規ユーザーのIDとメールアドレスを返す
	// CreateUserに成功した場合は、ポインタで渡したニューユーザーの
	//オブジェクトの内容が新しく作成したユーザーの内容に書き換わっていますので、
	// そこからIDメールの情報を取り出してユーザーレスポンスの新しい構造体の実体を作成して
	// resUser変数に格納してからリターン
	resUser := model.UserResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
	}
	// 成功した場合はエラーが発生していないので、エラーの返り値はNILにします。
	return resUser, nil
}

// Loginの実装
// userUsecase型をポインターレシーバーとして受け取る
// 引数としてUserモデルを受け取り、返り値としてGWTのstringとerrorを返す
func (uu *userUsecase) Login(user model.User) (string, error) {
	// uu.uv.UserValidateでuserオブジェクトにバリデーションをかける
	if err := uu.uv.UserValidate(user); err != nil {
		// バリデーションに失敗した場合は、返り値がstringとerrorになってますので、
		// 空の文字列とerrをリターンで返す
		return "", err
	}
	// クライアントから送られてきたメールがデータベースに存在するか
	// メールで検索するUserのオブジェクトを格納するための空のUserオブジェクトを作っておく
	storedUser := model.User{}
	// uu.ur.GetUserByEmailでUserRepositoryに定義されているGetUserByEmail関数を呼び出し
	// 第1引数にstoredUserのポインター第2引数には検索したいメールを渡す
	// 引数で受け取とるUserオブジェクトにメール情報が入ってますので、user.Emailで取り出して渡す
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		// エラーが発生した場合は空の文字列と、発生したエラーを返す
		return "", err
	}
	// クライアントが要求するメールが存在する場合は、続いてパスワードの検証を行う
	// CompareHashAndPassword関数を使い、データベース内に保存されてるハッシュ化されたパスワードと
	// クライアントから送られてきた平文のパスワードが一致するか検証を行う
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		// エラーが発生した場合は空の文字列と、発生したエラーを返す
		return "", err
	}
	// パスワードが一致する場合は、jwt tokenの生成を行う
	// HS256というアルゴリズムを指定、ペイロードの設定として、ユーザーIDとjwtの有効期限を設定
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// user_idはstoredUserのIDの値を割り当て
		"user_id": storedUser.ID,
		// jwt tokenの有効期限は、12時間に設定
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	})
	// NewWithClaimsの返り値に対してSignedStringメソッドを実行し実際にjwt token(tokenString)を生成
	// 引数のところでjwtのSECRETキーを渡す
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		// エラーが発生した場合は空の文字列と、発生したエラーを返す
		return "", err
	}
	// 成功した場合はjwt tokenとnilを返す
	return tokenString, nil
}
