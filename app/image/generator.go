package image

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	generator "github.com/lapis2411/card-generator"
)

type (
	ImageGenerateHandler interface {
		Generate(c echo.Context) error
	}

	CSVReceiveHandler struct{}
	ImageResponse     struct {
		ImageData string `json:"imageData"`
	}
)

var CSVReceiver = CSVReceiveHandler{}

func (CSVReceiveHandler) Generate(c echo.Context) error {
	styleBytes, err := readReceiveFile(c, "StyleFile")
	if err != nil {
		return fmt.Errorf("failed to receive style csv file: %w", err)
	}

	cardBytes, err := readReceiveFile(c, "CardFile")
	if err != nil {
		return fmt.Errorf("failed to receive card csv file: %w", err)
	}
	if err != nil {
		return fmt.Errorf("failed to receive csv data: %w", err)
	}
	cards, err := generator.MakeCards(styleBytes, cardBytes)
	imgs := cards.PrintImages() // []*image.RGBAを返却することを想定
	jsons := []ImageResponse{}
	for _, img := range imgs {
		b, err := encodeImageBase64(img)
		if err != nil {
			return fmt.Errorf("failed to encode image to base64: %w", err)
		}
		jsons = append(jsons, ImageResponse{ImageData: fmt.Sprintf("data:image/png;base64,%s", b)})
	}
	return c.JSON(http.StatusOK, jsons)
}

func readReceiveFile(c echo.Context, name string) ([]byte, error) {
	file, err := c.FormFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to receive file: %w", err)
	}
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()
	bs, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("failed to read style csv file: %w", err)
	}
	return bs, nil
}

func encodeImageBase64(img *image.RGBA) (string, error) {

	buffer := new(strings.Builder)
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	if err := png.Encode(encoder, img); err != nil {
		return "", fmt.Errorf("failed to encode image to PNG: %w", err)
	}
	encoder.Close()
	return buffer.String(), nil
}

// APIに対してファイルをアップロード
// アップロードされたファイルを元にサーバー側で画像を生成し、返却
//これを単一のAPIとしてJSで表示してあげれば良いらしい

// func main() {
// 	e := echo.New()

// 	e.GET("/", func(c echo.Context) error {
// 		// ファイルアップロードページを表示
// 		return c.HTML(http.StatusOK, `
// 			<input type="file" id="csvfile">
// 			<button onclick="upload()">Upload</button>
// 			<img id="output">
// 			<script>
// 				async function upload() {
// 					const file = document.getElementById('csvfile').files[0];
// 					const text = await file.text();
// 					const res = await fetch('/generate_image', {
// 						method: 'POST',
// 						body: text,
// 						headers: { 'Content-Type': 'text/plain' }
// 					});
// 					const dataUrl = await res.text();
// 					document.getElementById('output').src = dataUrl;
// 				}
// 			</script>
// 		`)
// 	})

// 	e.POST("/generate_image", func(c echo.Context) error {
// 		// CSVデータの取得
// 		csvData := c.Request().Body
// 		defer csvData.Close()

// 		// CSVの解析
// 		reader := csv.NewReader(csvData)
// 		records, err := reader.ReadAll()
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, "Unable to parse CSV")
// 		}

// 		// CSVデータを基に画像生成（ここでは最初の行の最初の値を画像の青の強さとする）
// 		blueIntensity, err := strconv.Atoi(records[0][0])
// 		if err != nil || blueIntensity < 0 || blueIntensity > 255 {
// 			return echo.NewHTTPError(http.StatusBadRequest, "Invalid CSV data")
// 		}

// 		img := image.NewRGBA(image.Rect(0, 0, 200, 200))
// 		blue := color.RGBA{0, 0, uint8(blueIntensity), 255}
// 		draw.Draw(img, img.Bounds(), &image.Uniform{blue}, image.Point{}, draw.Src)

// 		// 画像をbase64形式でエンコード
// 		buffer := new(strings.Builder)
// 		encoder := base64.NewEncoder(base64.StdEncoding, buffer)
// 		err = png.Encode(encoder, img)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, "Unable to encode image to PNG")
// 		}
// 		encoder.Close()

// 		// base64形式の画像をレスポンスとして返す
// 		return c.String(http.StatusOK, fmt.Sprintf("data:image/png;base64,%s", buffer.String()))
// 	})

// 	e.Logger.Fatal(e.Start(":8080"))
// }
