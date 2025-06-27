package tasks_service

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)
func TaskDone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	taskId := r.URL.Query().Get("id")
	if len(taskId) == 0 {
		writeError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(taskId)
	if err != nil {
		writeError(w, "Ошибка изменения id в int", http.StatusBadRequest)
		log.Println(err)
		return
	}

	var task Task

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	err = database.DB.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		writeError(w, "Задача не найдена", http.StatusNotFound)
		return
	} else if err != nil {
		writeError(w, "Ошибка сканирование с БД", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if task.Repeat == "" {
		query = "DELETE FROM scheduler WHERE id = ?"
		result, err := database.DB.Exec(query, task.ID)
		if err != nil {
			writeError(w, "Ошибка удаления задачи", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			writeError(w, "ошибка получения количества затронутых строк", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if rowsAffected == 0 {
			writeError(w, "Задача не найдена", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if _, err = w.Write([]byte("{}")); err != nil {
			log.Fatalf("Ошибка вывода")
		}
	} else {
		now := time.Now().Truncate(24 * time.Hour)

		dateTime, err := time.Parse(tools.DateFormat, task.Date)
		if err != nil {
			writeError(w, "дата представлена в формате, отличном от 20060102", http.StatusBadRequest)
			log.Println(err)
			return
		}

		date, err := tools.NextDate(now, dateTime.Format(tools.DateFormat), task.Repeat)
		if err != nil {
			writeError(w, "ошибка работы NextDate", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		query = "UPDATE scheduler SET date = ? WHERE id = ?"
		result, err := database.DB.Exec(query, date, task.ID)
		if err != nil {
			writeError(w, "ошибка обновление даты в бд", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			writeError(w, "ошибка получения количества затронутых строк", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if rowsAffected == 0 {
			writeError(w, "Задача не найдена", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if _, err = w.Write([]byte("{}")); err != nil {
			log.Fatalf("Ошибка вывода")
		}
	}
}
