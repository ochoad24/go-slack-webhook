package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Struct para parsear la respuesta JSON
type WebhookData struct {
	Job    string `json:"event_title"`
	Text   string `json:"text"`
	Url    string `json:"job_details_url"`
	Action string `json:"action"`
	Status string
	Emoji  string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file \n" + err.Error())
	}

	// Crea una instancia de Gin
	// gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{os.Getenv("PROXY")})

	// Define una ruta para manejar las solicitudes POST en "/webhook"
	r.POST("/webhook", func(c *gin.Context) {
		// Declara una variable para almacenar los datos de la solicitud
		var requestData WebhookData

		// Parsea los datos JSON de la solicitud en la estructura
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Muestra los datos en la consola
		requestData.Status = "Éxito"
		requestData.Emoji = ":white_check_mark:"
		if strings.Contains(requestData.Text, "failed") {
			requestData.Status = "Error"
			requestData.Emoji = ":warning:"
		}
		sendMessage(requestData)
		// Responde con un mensaje de éxito
		c.JSON(200, gin.H{"message": "Datos recibidos correctamente"})
	})

	// Inicia el servidor en el puerto 8080
	r.Run(os.Getenv("PORT"))
}

func sendMessage(requestData WebhookData) {
	// Test Job started
	if requestData.Action == "job_start" {
		return
	}

	// Definir la URL del webhook de Slack
	webhookURL := os.Getenv("WEBHOOKURL")
	// Crear el mensaje que deseas enviar
	message := fmt.Sprintf("%s\nTarea: %s \nEstado: %s\n<%s|Ver Job>", requestData.Emoji, requestData.Job, requestData.Status, requestData.Url)

	// Crear un mapa con los datos del mensaje
	payload := map[string]string{
		"text":       message,
		"icon_emoji": os.Getenv("ICON"),
		"username":   os.Getenv("USERNAME"),
		"channel":    os.Getenv("CHANNEL"),
	}

	// Codificar el mapa en formato JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// Enviar la solicitud POST al webhook de Slack
	_, err = http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		panic(err)
	}

}
