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

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("недопустимый формат данных: ожидаются две части, разделенные запятой")
	}

	stepsStrRaw := parts[0]
	stepsStrTrimmed := strings.TrimSpace(stepsStrRaw)

	if stepsStrRaw != stepsStrTrimmed {
		return 0, 0, fmt.Errorf("в значении количества шагов есть пробелы в начале или конце")
	}

	steps, err := strconv.Atoi(stepsStrRaw)
	if err != nil {
		return 0, 0, fmt.Errorf("недопустимое значение шага: %w", err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	durationStr := strings.TrimSpace(parts[1])
	dur, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("недопустимое значение продолжительности: %w", err)
	}
	if dur <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть больше 0")
	}
	return steps, dur, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, dur, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, dur)
	if err != nil {
		log.Println(err)
		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, calories)

	return result
}
