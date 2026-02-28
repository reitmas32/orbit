package settings

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config contiene variables que SÍ pueden ser modificadas desde .env
type Config struct {
	ENVIRONMENT string    `required:"true"`
	PORT        string    `required:"false" default:"8000"`
	TIMESTAMP   string    `required:"false"`
	ROOT_PATH   string    `required:"false" default:""`
	COLORED_LOG bool      `required:"false" default:"true"`
	SERVER_ID   uuid.UUID `required:"false"`

	CORE_API_URL string `required:"true"`
	CORE_API_KEY string `required:"false" default:""`

	RABBITMQ_URL      string `required:"true"`
	RABBITMQ_EXCHANGE string `required:"false" default:"orbit.chat"`

	JWT_PUBLIC_KEY string `required:"true"`
}

var Settings Config
var EnvDir = ".envs"

func LoadDotEnv() {

	err := godotenv.Load(fmt.Sprintf("%s/.env.base", EnvDir))
	if err != nil {
		log.Printf("No %s file found, using system environment variables", fmt.Sprintf("%s/.env.base.base", EnvDir))
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		log.Println("ENVIRONMENT is not set")
	}

	// Mapear el archivo .env.base correspondiente al entorno
	envFiles := map[string]string{
		"":            fmt.Sprintf("%s/.env", EnvDir),
		"local":       fmt.Sprintf("%s/.env.local", EnvDir),
		"development": fmt.Sprintf("%s/.env.dev", EnvDir),
		"production":  fmt.Sprintf("%s/.env.prod", EnvDir),
		"staging":     fmt.Sprintf("%s/.env.staging", EnvDir),
	}

	// Obtener el archivo de entorno correspondiente
	envFile, exists := envFiles[environment]
	if !exists {
		log.Printf("Environment '%s' is not supported. Must be one of: local, development, production, staging", environment)
	}

	// Cargar las variables desde el archivo correspondiente
	err = godotenv.Load(envFile)
	if err != nil {
		log.Printf("No %s file found, using system environment variables", envFile)
	} else {
		log.Printf("Loaded environment variables from %s", envFile)
	}
}

func LoadEnvs() {
	// Lista de variables inmutables que NO deben ser procesadas desde .env
	immutableVars := []string{
		"EDITION",
		"VERSION",
		"APP_NAME",
		"UUID_MODULE",
	}

	// Guardar temporalmente las variables inmutables del entorno
	// para restaurarlas después (en caso de que alguien intente modificarlas)
	immutableBackup := make(map[string]string)
	for _, varName := range immutableVars {
		if val := os.Getenv(varName); val != "" {
			immutableBackup[varName] = val
			// Eliminar temporalmente del entorno para que envconfig no las procese
			os.Unsetenv(varName)
		}
	}

	// Procesar las variables de entorno en la estructura Settings
	// Solo se procesarán las variables mutables
	err := envconfig.Process("", &Settings)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Restaurar las variables inmutables (aunque no se usarán)
	// Esto es solo para mantener el entorno limpio
	for varName := range immutableBackup {
		os.Unsetenv(varName)
	}

	Settings.TIMESTAMP = time.Now().Format("2006-01-02 15:04:05")
	Settings.SERVER_ID = uuid.MustParse("8e3de6ca-3ead-45d2-9a5b-b81b40d0c00a")

	// Imprimir las Settings si el entorno es local o development
	if Settings.ENVIRONMENT == "local" || Settings.ENVIRONMENT == "development" {
		log.Println("Settings (Mutable):")

		// Obtener el tipo y valor de la estructura Settings
		v := reflect.ValueOf(Settings)
		t := reflect.TypeOf(Settings)

		// Recorrer cada campo de la estructura
		for i := 0; i < v.NumField(); i++ {
			fieldName := t.Field(i).Name
			fieldValue := v.Field(i).Interface()
			fmt.Printf("\033[32m%s\033[0m: %v\n", fieldName, fieldValue)
		}
	}
}

// normalizeEdition normaliza el valor de la edición a minúsculas
func normalizeEdition(edition string) string {
	return strings.ToLower(strings.TrimSpace(edition))
}
