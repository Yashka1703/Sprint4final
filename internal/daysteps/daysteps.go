package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// Парсинг строки
func parsePackage(data string) (int, time.Duration, error) {
	dataToSlice := strings.Split(data, ",")
	if len(dataToSlice) != 2 {
		return 0, 0, fmt.Errorf("некорректные входные данные")
	}

	steps, err := strconv.Atoi(dataToSlice[0])
	if err != nil || steps <= 0 {
		return 0, 0, fmt.Errorf("некорректное количество шагов")
	}

	timer, err := time.ParseDuration(dataToSlice[1])
	if err != nil || timer <= 0 {
		return 0, 0, fmt.Errorf("некорректное время")
	}
	return steps, timer, nil
}

// Вычисление дневной активности
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("некорректный формат данных")
		return ""
	}
	if steps <= 0 || duration <= 0 {
		log.Println("некорректные шаги или время")
		return ""
	}

	distant := float64(steps) * stepLength / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return ""
	}
	result := fmt.Sprintf("Количество шагов: %d.\n"+
		"Дистанция составила %.2f км.\n"+
		"Вы сожгли %.2f ккал.\n",
		steps,
		distant,
		calories,
	)
	return result
}
