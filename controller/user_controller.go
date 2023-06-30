package controller

import (
	"go-rest-api/model"
	"go-rest-api/usecase"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type IUserController interface {
	// 引数としてはechoで定義されているContext型を受け取れるようにしておきます。
	// CsrfTokenというメソッドを追加していきます。
	SignUp(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
	CsrfToken(c echo.Context) error
}

// userControllerという構造体を定義
// userController型がIUserControllerのinterfaceを満たすためには
// 3つメソッドを実装する必要がある
type userController struct {
	// フィールドとしてusecaseパッケージ内の
	// IUserUsecaseインターフェース型の値をuuという名前で定義
	uu usecase.IUserUsecase
}

// controllerに対してusecaseを依存注入したいので、controllerの中にもコンストラクターを追加
// NewUserControllerというコンストラクターを追加して、
// 引数のところで外側でインスタンス化されるusecaseを引数として注入できるようにする
// 返り値の方は、IUserControllerのinterface型にする
func NewUserController(uu usecase.IUserUsecase) IUserController {
	// 受け取ったusecaseのインスタンスを使ってuserControllerの構造体の実体を生成して
	//そのアドレスをリターンで返す
	return &userController{uu}
}

// userControllerをucという名前でポインターレシーバーとして受け取るSignUpメソッドを定義
// 引数の型と返り値の型はinterfaceで書かれている内容と全く一緒にする必要があります
func (uc *userController) SignUp(c echo.Context) error {
	// クライアントから受け取るリクエストBODYの値を構造体に変換する処理
	// 0値で初期化されたUser構造体のオブジェクトを作成
	user := model.User{}
	// echo.Context(echoのコンテキスト)に準備されてるBindメソッドを実行
	// この時に引数にUserオブジェクトのポインターを渡します
	// そうするとクライアントから送られてくるリクエストボディーの値(model.User)を
	// ユーザーオブジェクトのポインターが指し示す先の値に格納
	if err := c.Bind(&user); err != nil {
		// 変換作業に失敗した場合は、c.JSONでクライアントBADリクエストのステータスと
		// エラーの内容をJSONで返す
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// バインドに成功した場合は、uc.uu.SignUpメソッドを呼び出す
	userRes, err := uc.uu.SignUp(user)
	if err != nil {
		// サインアップに失敗した場合は、c.JSONでクライアントに
		// StatusInternalServerError(ステータスコード)とエラーメッセージをクライアントに返す
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// 成功した場合は、c.JSONでStatusCreatedステータスと新しく作成したユーザー情報を返す
	return c.JSON(http.StatusCreated, userRes)
}

// Loginメソッドの実装
func (uc *userController) LogIn(c echo.Context) error {
	user := model.User{}
	// クライアントのリクエストボディーを構造体にバインド
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// uc.uu.Loginメソッドは、JWTtokenを生成してくれますので、
	// tokenStringという名前で受け取る
	tokenString, err := uc.uu.Login(user)
	if err != nil {
		// エラーが発生してしまった場合は、StatusInternalServerErrorのステータスとエラーメッセージをクライアントに返す
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// 成功した場合は、取得したJWT tokenをサーバーサイドでCookieに設定
	// new関数を使ってhttpパッケージに定義されてるクッキー構造体を新しく作成
	cookie := new(http.Cookie)
	// cookieのNAMEとしてtokenという名前を付ける
	cookie.Name = "token"
	// cookieの値は先ほど作成したJWT tokenを代入
	cookie.Value = tokenString
	// クッキーの有効期限は24時間にしています。
	cookie.Expires = time.Now().Add(24 * time.Hour)
	// パスはINDEXの/
	cookie.Path = "/"
	// ドメインは環境変数で設定しているAPIドメインの値を割り当てる
	cookie.Domain = os.Getenv("API_DOMAIN")
	// ※postmanで動作確認する時は一旦falseにする
	cookie.Secure = true
	//cookie.Secure = false
	// HttpOnly属性はtrueにしてクライアントのJAVASCRIPTから
	// tokenの値が読み取れないようにしておきます。
	cookie.HttpOnly = true
	// 今回は、フロントエンドとバックエンドのドメインが違うクロスドメイン間での
	// クッキーの送受信になりますので、SameSiteNoneModeにしています。
	cookie.SameSite = http.SameSiteNoneMode
	// c.SetCookieを使って今作成したクッキーの内容をHTTPレスポンスに含めるようにする
	c.SetCookie(cookie)
	// 最後にリターンでc.NoContentでStatusOKステータスをクライアントに返す
	return c.NoContent(http.StatusOK)
}

// LogOutメソッド
func (uc *userController) LogOut(c echo.Context) error {
	// new関数を使ってhttpパッケージに定義されてるクッキー構造体を新しく作成
	cookie := new(http.Cookie)
	// cookieのNAMEとしてtokenという名前を付ける
	cookie.Name = "token"
	// 値をクリアしたいのでからの文字列
	cookie.Value = ""
	// そして有効期限をNowにして有効期限がすぐにですね切れるようにする
	cookie.Expires = time.Now()
	// パスはINDEXの/
	cookie.Path = "/"
	// ドメインはAPIドメイン
	cookie.Domain = os.Getenv("API_DOMAIN")
	// ※postmanで動作確認する時は一旦falseにする
	cookie.Secure = true
	//cookie.Secure = false
	// HttpOnly属性はtrue
	cookie.HttpOnly = true
	// 今回は、フロントエンドとバックエンドのドメインが違うクロスドメイン間での
	// クッキーの送受信になりますので、SameSiteNoneModeにしています。
	cookie.SameSite = http.SameSiteNoneMode
	// c.SetCookieを使って今作成したクッキーの内容をHTTPレスポンスに含めるようにする
	c.SetCookie(cookie)
	// 最後にリターンでc.NoContentでStatusOKステータスをクライアントに返す
	return c.NoContent(http.StatusOK)
}

// CsrfTokenは、echoのコンテキストの中で、"csrf"のキーワードを使って取得することができます。
// string型に型アサーションをしてからJSONでクライアントにcsrf_tokenをレスポンスで返すようにしています。
func (uc *userController) CsrfToken(c echo.Context) error {
	token := c.Get("csrf").(string)
	return c.JSON(http.StatusOK, echo.Map{
		"csrf_token": token,
	})
}
