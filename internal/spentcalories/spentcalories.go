package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// Разбор строки
func parseTraining(data string) (int, string, time.Duration, error) {
	trainToSlice := strings.Split(data, ",")
	if len(trainToSlice) != 3 {
		return 0, "", 0, fmt.Errorf("wrong incoming format")
	}

	steps, err := strconv.Atoi(trainToSlice[0]) // получение количества шагов
	if err != nil {
		return 0, "", 0, fmt.Errorf("wrong incoming steps")
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("zero or negative number of steps")
	}

	timer, err := time.ParseDuration(trainToSlice[2]) // получение времени
	if err != nil {
		return 0, "", 0, fmt.Errorf("incorrect training time")
	}
	if timer <= 0 {
		return 0, "", 0, fmt.Errorf("zero or negative time")
	}
	activ := trainToSlice[1] // получение вида активности
	return steps, activ, timer, nil
}

// Вычисление общей пройденной дистанции
func distance(steps int, height float64) float64 {
	leghtStep := height * stepLengthCoefficient
	totalDistanceKm := (leghtStep * float64(steps)) / mInKm
	return totalDistanceKm
}

// Вычисление средней скорости
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()

}

// Получение информации о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activ, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}
	switch activ {
	case "Бег":
		distance := distance(steps, height)
		meanSpeed := meanSpeed(steps, height, duration)
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		result := fmt.Sprintf("Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
			activ,
			duration.Hours(),
			distance,
			meanSpeed,
			calories)
		return result, err

	case "Ходьба":
		distance := distance(steps, height)
		meanSpeed := distance / duration.Hours()
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		result := fmt.Sprintf("Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
			activ,
			duration.Hours(),
			distance,
			meanSpeed,
			calories)
		return result, err

	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}
}

// Вычисление калорий при ходьбе
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || duration <= 0 || weight <= 0 || height <= 0 {
		return 0, fmt.Errorf("wrong incoming format")
	}
	durationInMinutes := duration.Minutes()
	return (weight * meanSpeed(steps, height, duration) * durationInMinutes) / minInH, nil
}

// Вычисление калорий при беге
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || duration <= 0 || weight <= 0 || height <= 0 {
		return 0, fmt.Errorf("wrong incoming format")
	}
	durationInMinutes := duration.Minutes()
	return ((weight * meanSpeed(steps, height, duration) * durationInMinutes) / minInH) * walkingCaloriesCoefficient, nil
}
