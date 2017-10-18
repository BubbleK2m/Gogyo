package gogyo

import "fmt"

type MealName string

type Meal struct {
	Name MealName
}

type Menu map[int][]Meal

func NewMeal(name MealName) *Meal {
	meal := new(Meal)
	meal.Name = name

	return meal
}

func (meal Meal) String() string {
	return fmt.Sprintf("%s", meal.Name)
}
