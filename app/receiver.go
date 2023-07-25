package receiver

type (
	ReceiveHandler interface {
		func Receive(c echo.Context) error
	}

	CSVReceiveHandler struct {}
)
var CSVReceiver := CSVReceiveHandler{}

func (c CSVReceiveHandler) Receive(c echo.Context) error {
	// ファイルを取得する
	styleFile, err := c.FormFile("StyleFile")
	if err != nil {
		// エラー処理
	}

	cardFile, err := c.FormFile("CardFile")
	if err != nil {
		// エラー処理
	}

	// ファイルの内容を読み込む
	styleSrc, err := styleFile.Open()
	if err != nil {
		// エラー処理
	}
	defer styleSrc.Close()

	cardSrc, err := cardFile.Open()
	if err != nil {
		// エラー処理
	}
	defer cardSrc.Close()

	styleBytes, err := ioutil.ReadAll(styleSrc)
	if err != nil {
		// エラー処理
	}

	cardBytes, err := ioutil.ReadAll(cardSrc)
	if err != nil {
		// エラー処理
	}

	// 取得したファイルデータを使って何らかの処理を行う
	// ...

	// レスポンスを返す
	return c.String(http.StatusOK, "ファイルを受け取りました")
}

func CSVReceiver(data io.Reader) error {
	return nil
}
