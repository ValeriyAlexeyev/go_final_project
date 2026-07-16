package api

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

var (
	errEmptyRepeat   = errors.New("правило повторения не указано")
	errInvalidRepeat = errors.New("некорректное правило повторения")
)

// NextDate вычисляет следующую дату выполнения задачи.
//
// now — дата, относительно которой ищется следующее выполнение.
// dstart — исходная дата задачи в формате YYYYMMDD.
// repeat — правило повторения: d, y, w или m.
func NextDate(now time.Time, dstart, repeat string) (string, error) {
	start, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата %q: %w", dstart, err)
	}

	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errEmptyRepeat
	}

	// Убираем время, оставляем только календарную дату.
	now = normalizeDate(now)

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errInvalidRepeat
	}

	switch parts[0] {
	case "d":
		return nextDateByDays(now, start, parts)

	case "y":
		return nextDateByYear(now, start, parts)

	case "w":
		return nextDateByWeekdays(now, start, parts)

	case "m":
		return nextDateByMonthDays(now, start, parts)

	default:
		return "", fmt.Errorf("%w: неизвестный тип %q", errInvalidRepeat, parts[0])
	}
}

// nextDateByDays обрабатывает правило d N.
func nextDateByDays(now, start time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("%w: правило d должно иметь вид d N", errInvalidRepeat)
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("%w: интервал должен быть числом", errInvalidRepeat)
	}

	if interval < 1 || interval > 400 {
		return "", fmt.Errorf(
			"%w: интервал должен быть от 1 до 400",
			errInvalidRepeat,
		)
	}

	date := start

	for {
		date = date.AddDate(0, 0, interval)

		if afterNow(date, now) {
			return date.Format(DateFormat), nil
		}
	}
}

// nextDateByYear обрабатывает правило y.
func nextDateByYear(now, start time.Time, parts []string) (string, error) {
	if len(parts) != 1 {
		return "", fmt.Errorf("%w: правило y не содержит параметров", errInvalidRepeat)
	}

	date := start

	for {
		date = date.AddDate(1, 0, 0)

		if afterNow(date, now) {
			return date.Format(DateFormat), nil
		}
	}
}

// nextDateByWeekdays обрабатывает правило:
//
//	w 1
//	w 1,4,5
//
// В Go воскресенье имеет номер 0, поэтому преобразуем его в 7.
func nextDateByWeekdays(
	now time.Time,
	start time.Time,
	parts []string,
) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf(
			"%w: правило w должно иметь вид w 1,2,3",
			errInvalidRepeat,
		)
	}

	weekdays, err := parseNumberSet(parts[1], 1, 7)
	if err != nil {
		return "", fmt.Errorf("некорректные дни недели: %w", err)
	}

	date := start

	// Максимально достаточно проверить несколько лет вперёд,
	// но при корректном правиле подходящий день найдётся за 7 суток.
	for i := 0; i < 366*10; i++ {
		date = date.AddDate(0, 0, 1)

		weekday := int(date.Weekday())
		if weekday == 0 {
			weekday = 7
		}

		if weekdays[weekday] && afterNow(date, now) {
			return date.Format(DateFormat), nil
		}
	}

	return "", errors.New("не удалось определить следующую дату")
}

// nextDateByMonthDays обрабатывает правила:
//
//	m 4
//	m 1,15,25
//	m -1
//	m -2
//	m 3 1,3,6
//	m 1,-1 2,8
func nextDateByMonthDays(
	now time.Time,
	start time.Time,
	parts []string,
) (string, error) {
	if len(parts) != 2 && len(parts) != 3 {
		return "", fmt.Errorf(
			"%w: правило m должно иметь вид m <дни> [месяцы]",
			errInvalidRepeat,
		)
	}

	days, err := parseMonthDays(parts[1])
	if err != nil {
		return "", err
	}

	months := make(map[int]bool, 12)

	if len(parts) == 3 {
		months, err = parseNumberSet(parts[2], 1, 12)
		if err != nil {
			return "", fmt.Errorf("некорректные месяцы: %w", err)
		}
	} else {
		for month := 1; month <= 12; month++ {
			months[month] = true
		}
	}

	date := start

	// Проверяем даты последовательно. Ограничение защищает
	// функцию от бесконечного цикла при ошибке алгоритма.
	for i := 0; i < 366*20; i++ {
		date = date.AddDate(0, 0, 1)

		if !afterNow(date, now) {
			continue
		}

		if !months[int(date.Month())] {
			continue
		}

		if matchesMonthDay(date, days) {
			return date.Format(DateFormat), nil
		}
	}

	return "", errors.New("не удалось определить следующую дату")
}

// parseNumberSet разбирает строку вида 1,4,5.
func parseNumberSet(
	value string,
	minValue int,
	maxValue int,
) (map[int]bool, error) {
	result := make(map[int]bool)

	if strings.TrimSpace(value) == "" {
		return nil, errors.New("список значений пуст")
	}

	items := strings.Split(value, ",")

	for _, item := range items {
		if item == "" {
			return nil, errors.New("обнаружено пустое значение")
		}

		number, err := strconv.Atoi(item)
		if err != nil {
			return nil, fmt.Errorf("%q не является числом", item)
		}

		if number < minValue || number > maxValue {
			return nil, fmt.Errorf(
				"значение %d должно быть от %d до %d",
				number,
				minValue,
				maxValue,
			)
		}

		if result[number] {
			return nil, fmt.Errorf("значение %d указано повторно", number)
		}

		result[number] = true
	}

	return result, nil
}

// parseMonthDays разбирает дни месяца от 1 до 31,
// а также специальные значения -1 и -2.
func parseMonthDays(value string) ([]int, error) {
	if strings.TrimSpace(value) == "" {
		return nil, errors.New("дни месяца не указаны")
	}

	items := strings.Split(value, ",")
	result := make([]int, 0, len(items))
	seen := make(map[int]bool)

	for _, item := range items {
		if item == "" {
			return nil, errors.New("обнаружен пустой день месяца")
		}

		day, err := strconv.Atoi(item)
		if err != nil {
			return nil, fmt.Errorf("%q не является днём месяца", item)
		}

		if day != -1 && day != -2 && (day < 1 || day > 31) {
			return nil, fmt.Errorf(
				"день месяца %d должен быть от 1 до 31, -1 или -2",
				day,
			)
		}

		if seen[day] {
			return nil, fmt.Errorf("день месяца %d указан повторно", day)
		}

		seen[day] = true
		result = append(result, day)
	}

	sort.Ints(result)

	return result, nil
}

// matchesMonthDay проверяет соответствие даты одному из дней.
//
// -1 — последний день месяца;
// -2 — предпоследний день месяца.
func matchesMonthDay(date time.Time, allowedDays []int) bool {
	lastDay := daysInMonth(date.Year(), date.Month())

	for _, allowedDay := range allowedDays {
		switch allowedDay {
		case -1:
			if date.Day() == lastDay {
				return true
			}

		case -2:
			if date.Day() == lastDay-1 {
				return true
			}

		default:
			if date.Day() == allowedDay {
				return true
			}
		}
	}

	return false
}

func daysInMonth(year int, month time.Month) int {
	// Нулевой день следующего месяца — последний день текущего.
	return time.Date(
		year,
		month+1,
		0,
		0,
		0,
		0,
		0,
		time.UTC,
	).Day()
}

func afterNow(date, now time.Time) bool {
	return normalizeDate(date).After(normalizeDate(now))
}

func normalizeDate(date time.Time) time.Time {
	return time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(
			w,
			"метод не поддерживается",
			http.StatusMethodNotAllowed,
		)
		return
	}

	nowValue := r.FormValue("now")
	dateValue := r.FormValue("date")
	repeatValue := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowValue == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowValue)
		if err != nil {
			http.Error(w, "некорректный параметр now: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateValue, repeatValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(next))
}
