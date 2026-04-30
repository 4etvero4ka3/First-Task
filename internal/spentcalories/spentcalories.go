package spentcalories

import (
	"errors"
	"fmt"
	"log"
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

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("недопустимый формат данных: ожидаются три части, разделенные запятыми")
	}

	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("недопустимое значение шага: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	activity := strings.TrimSpace(parts[1])

	durationStr := strings.TrimSpace(parts[2])
	dur, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("недопустимое значение продолжительности: %w", err)
	}
	if dur <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	if strings.Contains(durationStr, "-") && !strings.HasPrefix(durationStr, "-") {
		return 0, "", 0, fmt.Errorf("недопустимое значение продолжительности: отрицательные интервалы")
	}

	return steps, activity, dur, nil
}

func distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient

	meters := float64(steps) * stepLen

	return meters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	d := distance(steps, height)

	hours := duration.Hours()

	return d / hours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, dur, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var distanceKm, speed, calories float64

	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, dur)
		if err != nil {
			log.Println(err)
			return "", err
		}
		distanceKm = distance(steps, height)
		speed = meanSpeed(steps, height, dur)

	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, dur)
		if err != nil {
			log.Println(err)
			return "", err
		}
		distanceKm = distance(steps, height)
		speed = meanSpeed(steps, height, dur)

	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activity)
	}

	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity, dur.Hours(), distanceKm, speed, calories)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	speed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	calories := (weight * speed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("длительность должна быть больше 0")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / minInH
	calories *= walkingCaloriesCoefficient
	return calories, nil
}
