package image

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	generator "github.com/lapis2411/card-generator"
)

type (
	ImageGenerateHandler interface {
		Generate(c echo.Context) error
	}

	ImageResponse struct {
		ImageData string `json:"imageData"`
	}
	ImageBase64 struct {
		image string
	}
)

// PrintsJsons make base64 png for response in json format by style and card csv
func PrintsJsons(styleBytes, cardBytes []byte) ([]ImageResponse, error) {
	b64s, err := base64Images(styleBytes, cardBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate images: %w", err)
	}
	jsons := make([]ImageResponse, len(b64s))
	for i, b64 := range b64s {
		jsons[i] = b64.ToPNGImageResponse()
	}
	return jsons, nil
}

func (i ImageBase64) ToPNGImageResponse() ImageResponse {
	s := fmt.Sprintf("data:image/png;base64,%s", i.image)
	return ImageResponse{ImageData: s}
}

func base64Images(styles, cards []byte) ([]ImageBase64, error) {
	cts, err := generator.MakeCards(styles, cards)
	fp := os.Getenv("FontPath")
	gen := generator.NormalSizeCardGenerator(fp)
	cds, err := gen.Generate(cts)
	if err != nil {
		return nil, fmt.Errorf("failed to generate cards: %w", err)
	}
	l := generator.NewA4Layout()
	cvs, err := l.Arrange(cds)
	if err != nil {
		return nil, fmt.Errorf("failed to arrange cards: %w", err)
	}
	imgs := cvs.ToImageRGBA()
	base64s := make([]ImageBase64, len(imgs))
	for i, img := range imgs {
		b64i, err := encodeImageBase64(img)
		if err != nil {
			return nil, fmt.Errorf("failed to encode image to base64: %w", err)
		}
		base64s[i].image = b64i
	}
	return base64s, nil
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
