package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strings"

	"github.com/Knetic/govaluate"
)

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func calculate(expression string) (interface{}, error) {
	matched, err := regexp.MatchString(`^[\s\d+\-*/().]+$`, expression)
	if err != nil || !matched {
		return nil, fmt.Errorf("Expression is not valid")
	}

	expression = strings.ReplaceAll(expression, ",", ".")

	expression, err = simplifyExpression(expression)
	if err != nil {
		return nil, err
	}
	log.Printf("expression: %s\n", expression)

	expressionTree, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Expression is not valid: %w", err)
	}

	result, err := expressionTree.Evaluate(nil)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Expression is not valid: %w", err)
	}

	return result, nil
}

func simplifyExpression(expression string) (string, error) {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(expression, ""), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := calculate(req.Expression)

	// # Как сформировать тело ответа?
	// json.Marshal(v any) ([]byte, error)
	// example:
	// resp := Response{Error: err.Error()}
	// resp := Response{Result: 6}
	// byteArr, err = json.Marshal(resp)

	// # Как отправить ответ пользователю?
	// w : interface http.ResponseWriter — используется для отправки ответа пользователю на запрос

	// w.WriteHeader(statusCode int)
	// example:
	// w.WriteHeader(422)

	// w.Write([]byte) (int, error)
	// example:
	// bytesWritten, err = w.Write(byteArr)

	var httpCode int = http.StatusOK
	var resp Response

	if err != nil {
		log.Printf("error: %s\n", err.Error())

		resp = Response{Error: err.Error()}
		httpCode = http.StatusInternalServerError

		if strings.Contains(err.Error(), "Expression is not valid") {
			// вернуть ошибку с кодом 422, если ошибка Expression is not valid: {"error": "Expression is not valid"}
			httpCode = http.StatusUnprocessableEntity
		}
	} else if result == math.Inf(1) {
		// вернуть ошибку с кодом 500 в случае ошибки деления на ноль: {"error": "division by zero"}
		resp = Response{Error: "division by zero"}
		httpCode = http.StatusInternalServerError
	} else {
		log.Printf("result: %s\n", result)

		// вернуть успешный результат с кодом 200: {"result":6}
		resp = Response{Result: result}
	}

	w.WriteHeader(httpCode)
	byteArr, err := json.Marshal(resp)
	if err != nil {
		log.Println(err.Error())
	}
	w.Write(byteArr)
}

func main() {
	http.HandleFunc("/api/v1/calculate", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
