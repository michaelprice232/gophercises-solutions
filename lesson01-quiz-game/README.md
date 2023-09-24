# quiz-game

My solution from: https://courses.calhoun.io/lessons/les_goph_01

Small app to read some quiz questions and answers from a CSV file and then ask them interactively to the user via the terminal.
Offers a timer and option to randomise the questions.

## Usage

```text
% go run main.go --help                                
Usage of main:
  -csv string
        a csv file in the format of 'question,answer' (default problems.csv) (default "problems.csv")
  -limit string
        the time limit for the quiz (default 30s) (default "30s")
  -random
        whether to randomise the questions (default false)
        

% go run main.go                                       
Enter any key to start timer (30s): 
Question #1: are you male? = yes
Question #2: 5+5 = 10
Question #3: what 2+2, sir? = 4

You answered 3 out of 3 correct!%
```
