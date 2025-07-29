package httpadapter

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"path/filepath"
	"strings"
	"ticketera/internal/adapters/storage"
	"ticketera/internal/domain"
)

// LogoHandler recibe el logo, lo guarda en base64 en MongoDB y responde con la data URL
func LogoHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(2 << 20); err != nil {
		http.Error(w, "Error al leer el formulario", http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("logo")
	if err != nil {
		http.Error(w, "No se recibió archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	var img image.Image
	var decodeErr error
	if ext == ".png" {
		img, decodeErr = png.Decode(file)
	} else if ext == ".jpg" || ext == ".jpeg" {
		img, decodeErr = jpeg.Decode(file)
	} else {
		http.Error(w, "Formato no soportado. Solo PNG o JPG", http.StatusBadRequest)
		return
	}
	if decodeErr != nil {
		http.Error(w, "No se pudo decodificar la imagen", http.StatusBadRequest)
		return
	}
	// Redimensionar si es mayor a 200px
	bounds := img.Bounds()
	maxDim := 200
	width, height := bounds.Dx(), bounds.Dy()
	scale := 1.0
	if width > maxDim || height > maxDim {
		if width > height {
			scale = float64(maxDim) / float64(width)
		} else {
			scale = float64(maxDim) / float64(height)
		}
	}
	newW := int(float64(width) * scale)
	newH := int(float64(height) * scale)
	// Redimensionar usando image/draw
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	// Usar NearestNeighbor para simplicidad (puedes cambiar a CatmullRom si quieres más calidad)
	drawNearestNeighbor(dst, img)
	// Codificar a base64
	var buf bytes.Buffer
	mime := "image/png"
	if ext == ".png" {
		png.Encode(&buf, dst)
	} else {
		mime = "image/jpeg"
		jpeg.Encode(&buf, dst, nil)
	}
	b64 := "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	repo := storage.NewMongoRepository("mongodb://localhost:27017")
	coll := repo.Client.Database("ticketera").Collection("logo")
	coll.DeleteMany(r.Context(), map[string]interface{}{})
	logo := domain.Logo{FileURL: b64}
	_, err = coll.InsertOne(r.Context(), logo)
	if err != nil {
		http.Error(w, "Error al guardar en base de datos", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": b64})
}

// drawNearestNeighbor realiza un escalado simple de la imagen src a dst
func drawNearestNeighbor(dst *image.RGBA, src image.Image) {
	dx := dst.Bounds().Dx()
	dy := dst.Bounds().Dy()
	sx := src.Bounds().Dx()
	sy := src.Bounds().Dy()
	for y := 0; y < dy; y++ {
		for x := 0; x < dx; x++ {
			srcX := int(float64(x) * float64(sx) / float64(dx))
			srcY := int(float64(y) * float64(sy) / float64(dy))
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}
}
