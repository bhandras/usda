package usda

type FoodGroup struct {
	FoodGroupId int    `field:"0"` // 4-digit code identifying a food group.
	Desc        string `field:"1"` // Name of food group.
}

type FoodDesc struct {
	FoodId         string `field:"0"` // 5-digit Nutrient Databank number that uniquely identifies a food item.
	FoodGroupId    int    `field:"1"` // 4-digit code indicating food group to which a food item belongs.
	LongDesc       string `field:"2"` // 200-character description of food item.
	ShortDesc      string `field:"3"` // 60-character abbreviated description of food item.
	CommonNames    string `field:"4"` // Other names commonly used to describe a food, including local or regional names for various foods.
	Manufacturer   string `field:"5"` // Indicates the company that manufactured the product, when appropriate.
	RefuseDesc     string `field:"7"` // Description of inedible parts of a food item (refuse), such as seeds or bone.
	Refuse         int    `field:"8"` // Percentage of refuse.
	ScientificName string `field:"9"` // Scientific name of the food item.
}

type NutrientDef struct {
	NutrientId int    `field:"0"` // Unique 3-digit identifier code for a nutrient.
	Units      string `field:"1"` // Units of measure (mg, g, ug, and so on).
	Desc       string `field:"3"` // Name of nutrient/food component.
	Order      int    `field:"5"` // Used to sort nutrient records in the same order as various reports produced from SR.
}

type NutrientValue struct {
	FoodId     string  `field:"0"`  // 5-digit Nutrient Databank number that uniquely identifies a food item.
	NutrientId int     `field:"1"`  // Unique 3-digit identifier code for a nutrient.
	Value      float32 `field:"2"`  // Amount in 100 grams, edible portion
	Min        float32 `field:"10"` // Minimum value.
	Max        float32 `field:"11"` // Maximum value.
}

type WeightValue struct {
	FoodId     string  `field:"0"` // 5-digit Nutrient Databank number that uniquely identifies a food item.
	Seq        int     `field:"1"` // Sequence number.
	Amount     float32 `field:"2"` // Unit modifier (for example, 1 in “1 cup”).
	Desc       string  `field:"3"` // Description (for example, cup, diced, and 1-inch pieces).
	GramWeight float32 `field:"4"` // Gram weight.
}
