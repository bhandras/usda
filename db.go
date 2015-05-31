package usda

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)

type DB struct {
	FoodGroups     []FoodGroup
	Nutrients      []NutrientDef
	NutrientValues []NutrientValue
	WeightValues   []WeightValue
	Foods          map[string]*FoodDesc
}

type ParseError struct {
	idx  int
	kind reflect.Kind
}

func (e *ParseError) Error() string {
	return "usda: Parse error at position " + string(e.idx) + ", type != " + e.kind.String()
}

func (db *DB) Read(dir string) error {
	db.Foods = make(map[string]*FoodDesc)

	// food groups
	if err := scanFile(path.Join(dir, "FD_GROUP.txt"), 2, func(parts []string) error {
		foodGroup := &FoodGroup{}
		if err := fillStruct(foodGroup, parts); err != nil {
			return err
		}
		db.FoodGroups = append(db.FoodGroups, *foodGroup)
		return nil
	}); err != nil {
		return err
	}

	// food descriptions
	if err := scanFile(path.Join(dir, "FOOD_DES.txt"), 14, func(parts []string) error {
		foodDesc := &FoodDesc{}
		if err := fillStruct(foodDesc, parts); err != nil {
			return err
		}
		db.Foods[foodDesc.FoodId] = foodDesc
		return nil
	}); err != nil {
		return err
	}

	// nutrient definitions
	if err := scanFile(path.Join(dir, "NUTR_DEF.txt"), 6, func(parts []string) error {
		nutrientDef := &NutrientDef{}
		if err := fillStruct(nutrientDef, parts); err != nil {
			return err
		}
		db.Nutrients = append(db.Nutrients, *nutrientDef)
		return nil
	}); err != nil {
		return err
	}

	// nutrient data
	if err := scanFile(path.Join(dir, "NUT_DATA.txt"), 18, func(parts []string) error {
		nutrientValue := &NutrientValue{}
		if err := fillStruct(nutrientValue, parts); err != nil {
			return err
		}
		db.NutrientValues = append(db.NutrientValues, *nutrientValue)
		return nil
	}); err != nil {
		return err
	}

	// weight data
	if err := scanFile(path.Join(dir, "WEIGHT.txt"), 7, func(parts []string) error {
		weightValue := &WeightValue{}
		if err := fillStruct(weightValue, parts); err != nil {
			return err
		}
		db.WeightValues = append(db.WeightValues, *weightValue)
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func scanFile(path string, columns int, lf func([]string) error) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts, err := splitLine(line, columns)

		if err != nil {
			return err
		}

		err = lf(parts)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func splitLine(line string, columns int) ([]string, error) {
	parts := strings.Split(line, "^")
	if len(parts) != columns {
		return nil, fmt.Errorf("len(parts) = %d, expected = %d", len(parts), columns)
	}

	for i := 0; i < len(parts); i++ {
		parts[i] = strings.Trim(parts[i], "~")
	}

	return parts, nil
}

func fillStruct(m interface{}, parts []string) error {
	typ := reflect.TypeOf(m)

	if typ.Kind() != reflect.Ptr {
		return fmt.Errorf("%s must be a pointer to a struct: ", typ.Kind())
	}

	typ = typ.Elem()

	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("m must be a pointer to a struct: ", typ.Kind())
	}

	ps := reflect.ValueOf(m)
	s := ps.Elem()

	for i := 0; i < s.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("field")
		if tag == "" {
			// non tagged struct field
			continue
		}

		f := s.Field(i)
		idx, err := strconv.Atoi(tag)

		if err != nil {
			return fmt.Errorf("usda: Field tag \"%S\" type != int", f.String())
		}

		if idx < 0 || len(parts) < idx {
			return fmt.Errorf("usda: Field tag out of range: %s:%d", f.String(), idx)
		}

		if parts[idx] != "" && f.IsValid() && f.CanSet() {
			switch f.Kind() {
			case reflect.Int:
				value, err := strconv.ParseInt(parts[idx], 10, 32)
				if err != nil {
					return &ParseError{idx, f.Kind()}
				}
				f.SetInt(value)
			case reflect.String:
				f.SetString(parts[idx])
			case reflect.Float32:
				value, err := strconv.ParseFloat(parts[idx], 32)
				if err != nil {
					return &ParseError{idx, f.Kind()}
				}
				f.SetFloat(value)
			}
		}
	}

	return nil
}
