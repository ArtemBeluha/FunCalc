package main

import (
 "encoding/json"
 "fmt"
 "log"
 "net/http"
 "regexp"

 "github.com/Knetic/govaluate"
)

type Request struct {
 Expression string `json:"expression"`
}

type Response struct {
 Result interface{} `json:"result,omitempty"`
 Error string      `json:"error,omitempty"`
}

func calculate(expression string) (interface{}, error) {
 matched, err := regexp.MatchString(`^[\d+\-*/().]+$`, expression)
 if err != nil || !matched {
  return nil, fmt.Errorf("Expression is not valid")
 }

 expression = strings.ReplaceAll(expression, ",", ".")

 expression, err = simplifyExpression(expression)
 if err != nil {
  return nil, err
 }


 expressionTree, err := govaluate.NewEvaluableExpression(expression)
 if err != nil {
  return nil, fmt.Errorf("invalid expression: %w", err)
 }

 result, err := expressionTree.Evaluate(nil)
 if err != nil {
  if strings.Contains(err.Error(), "division by zero") {
   return nil, fmt.Errorf("division by zero")
  }
  return nil, fmt.Errorf("error evaluating expression: %w", err)
 }

 return result, nil
}

func simplifyExpression(expression string) (string, error){
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
 if err != nil {
  resp := Response{Error: err.Error()}
  if strings.Contains(err.Error(), "invalid expression") {
   http.Error(w, resp.Error, http.StatusUnprocessableEntity) //422
  } else {
   http.Error(w, resp.Error, http.StatusInternalServerError) //550
  }
  return
 }

 resp := Response{Result: result}
 json.NewEncoder(w).Encode(resp)
}
