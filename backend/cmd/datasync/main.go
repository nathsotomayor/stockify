package main

import (
	"log"
	"stockify/internal/config"
	"stockify/internal/database"
	"stockify/internal/store"
	"stockify/internal/tasks"
)

func main() {
	log.Println("Iniciando script de sincronización de datos CLI...")

	cfg := config.Load()
	db := database.Connect(cfg.DatabaseURL)

	stockSt := store.NewStockStore(db)

	count, err := stockSt.CountStocks()
	if err != nil {
		log.Fatalf("Error al verificar el conteo de stocks en la base de datos: %v", err)
	}

	if count > 0 {
		log.Printf("La base de datos ya contiene %d registros de stocks. No se ejecutará la población.", count)
	} else {
		log.Println("La base de datos está vacía o no contiene registros de stocks. Iniciando población...")
		dataSyncSvc := tasks.NewDataSyncService(db, cfg)
		err := dataSyncSvc.RunPopulation()
		if err != nil {
			log.Fatalf("Falló la ejecución de la población de datos: %v", err)
		}
		log.Println("Población de datos completada exitosamente.")
	}

	log.Println("Script de sincronización de datos CLI finalizado.")
}
