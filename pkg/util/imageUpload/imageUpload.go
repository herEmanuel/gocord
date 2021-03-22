package imageUpload

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ImageUpload(ctx *fiber.Ctx) error {

	fileForm, err := ctx.MultipartForm()
	fmt.Println(fileForm)
	if err != nil {
		return ctx.Status(500).SendString("Could not parse the multipart form")
	}

	images := fileForm.File["images"]
	if len(images) == 0 {
		return ctx.Status(400).SendString("No images provided")
	}

	for _, image := range images {

		if image.Size > 1024*1024*5 {
			return ctx.Status(400).SendString("This image is too big")
		}

		validType := false
		for _, imageType := range []string{"image/png", "image/jpeg"} {
			if imageType == image.Header["Content-Type"][0] {
				validType = true
			}
		}
		if !validType {
			return ctx.Status(400).SendString("Invalid image type")
		}

		imageName := fmt.Sprintf("%v$%v%v", strings.TrimSuffix(image.Filename, filepath.Ext(image.Filename)), time.Now().Unix(), filepath.Ext(image.Filename))
		err := ctx.SaveFile(image, "../../assets/uploads/"+imageName)
		if err != nil {
			return ctx.Status(500).SendString("Could not save the file, " + err.Error())
		}

		ctx.Locals("imagePath", "uploads/"+imageName)
	}

	return ctx.Next()
}
