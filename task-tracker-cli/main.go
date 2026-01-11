package main

import (
	"encoding/json"
	"fmt"
	"os"
	"bufio"
)

type Task struct{
	ID int				`json:"ID"`
	Description string	`json:"Description"`
	IsDone bool			`json:"IsDone"`
}

var discard string

func check_selection(input_id int, task_slice []Task) {
	if input_id > len(task_slice) {
		fmt.Print("Wrong selection!")
		return
	}
}

func add_task(task_slice *[]Task, task_dictionary string){
	var desc string
	var is_done bool
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Enter the Description(with_no_space) - ")
	if scanner.Scan() {
		desc = scanner.Text()
		fmt.Scanf("%s", &desc)
	}
	fmt.Printf("Enter the status(true/false) - ")
	fmt.Scan(&discard)
	fmt.Scanf("%t\n", &is_done)
	newTask := Task{ID: len(*task_slice) + 1, Description: desc,IsDone: is_done}
	*task_slice = append(*task_slice, newTask)
	write_to_json(task_slice, task_dictionary)
}

func print_task(task_slice *[]Task){
	for _, task := range *task_slice{
		fmt.Printf("\n%d - %s - Done: %t", task.ID, task.Description, task.IsDone)
	}
	fmt.Println()
}

func update_task(task_slice []Task, task_dictionary string) []Task {
	var input_id int
	print_task(&task_slice)
	fmt.Printf("Enter the Id of task to update: ")
	fmt.Scanf("\n%d",&input_id)
	check_selection(input_id, task_slice)
	input_id = input_id-1
	if task_slice[input_id].IsDone{
		task_slice[input_id].IsDone = false
	} else {
		task_slice[input_id].IsDone = true
	}
	write_to_json(&task_slice, task_dictionary)
	return task_slice
}

func delete_task(task_slice []Task, task_dictionary string) []Task {
	var input_id int
	print_task(&task_slice)
	fmt.Printf("Enter the Id of task to delete: ")
	fmt.Scanf("\n%d",&input_id)
	check_selection(input_id, task_slice)
	input_id = input_id-1
	task_slice = append(task_slice[:input_id], task_slice[input_id+1:]...)
	write_to_json(&task_slice, task_dictionary)
	return task_slice
}

func write_to_json(task_slice *[]Task, task_dictionary string) {
	jsonData, err := json.MarshalIndent(*task_slice, "", "  ")
	if err != nil{
		panic(err)
	}
	err = os.WriteFile(task_dictionary, jsonData, 0644)
	if err != nil{
		panic(err)
	}
}

func main(){
	task_slice := []Task{}
	task_dictionary := "tasks.json"
	_, err1 := os.Stat(task_dictionary)
	if err1 == nil {
		file, err := os.ReadFile(task_dictionary)
		if err == nil{
			json.Unmarshal(file, &task_slice)
		}
	}
	var userChoice int
	fmt.Println(" 1. Print\n 2. Add Task\n 3. Update\n 4. Delete")
	fmt.Scan(&userChoice)
	switch userChoice {
		case 1:
			print_task(&task_slice)
		case 2:
			add_task(&task_slice, task_dictionary)
		case 3:
			task_slice = update_task(task_slice, task_dictionary)
		case 4:
			task_slice = delete_task(task_slice, task_dictionary)
		default:
			fmt.Print("Exit")
			return
	}
}