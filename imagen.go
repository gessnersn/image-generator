package main

import (
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/zdebeer99/goexpression"
	"github.com/joho/godotenv"
)

var baseTpl = template.Must(template.ParseFiles("base.html"))

func main() {
	port := getPortFromEnv()
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/imagen/", imagenHandler)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func getPortFromEnv() string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return port
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(
		writer,
		"This is the index page, you can also explore the '/imagen' subpage!",
	)
}

func imagenHandler(writer http.ResponseWriter, request *http.Request) {
	special := "form input:\n"
	formError := request.ParseForm()
	check(formError)
	width, widthExists := request.Form["width"]
	height, heightExists := request.Form["height"]
	red, redExists := request.Form["r"]
	green, greenExists := request.Form["g"]
	blue, blueExists := request.Form["b"]

	imagePath := ""
	if widthExists && heightExists && redExists && greenExists && blueExists {
		imagePath = computeImage(width, height, red, green, blue)
		imagePath = "/assets/outimage.png"
	}

	for key, value := range request.Form {
		special += fmt.Sprintf("%v: %v\n", key, value)
	}
	data := struct {
		Title string
		Image string
	}{
		Title: "Imagen - Image Generator",
		Image: imagePath,
	}
	templateError := baseTpl.Execute(writer, data)
	check(templateError)
}

func computeImage(width []string, height []string, red []string, green []string, blue []string) string {
	outPath := "assets/outimage.png"

	widthString := strings.Join(width, "")
	widthInt, _ := strconv.Atoi(widthString)
	heightString := strings.Join(height, "")
	heightInt, _ := strconv.Atoi(heightString)
	redString := strings.Join(red, "")
	greenString := strings.Join(green, "")
	blueString := strings.Join(blue, "")

	newImage := Image{
		widthInt, heightInt, redString, greenString, blueString,
	}
	f, err := os.Create(outPath)
	check(err)
	defer f.Close()
	err = png.Encode(f, newImage)
	check(err)
	return outPath
}

type Image struct {
	width  int
	height int
	red    string
	green  string
	blue   string
}

func (i Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (i Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.width, i.height)
}

func (i Image) At(x, y int) color.Color {
	redExpression := i.red
	greenExpression := i.green
	blueExpression := i.blue

	// unfortunately "github.com/zdebeer99/goexpression" does not support ^ (exponent) expression
	context := map[string]interface{}{
		"x": x,
		"y": y,
	}

	var redValue, greenValue, blueValue float64
	redValue = goexpression.Eval(redExpression, context)
	greenValue = goexpression.Eval(greenExpression, context)
	blueValue = goexpression.Eval(blueExpression, context)

	col := color.RGBA{uint8(redValue), uint8(greenValue), uint8(blueValue), 255}
	return col
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
