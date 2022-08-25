package linear

import (
	"database/sql"
	"fmt"

	"encoding/json"

	"errors"

	_ "github.com/gonum/blas"
	_ "github.com/lib/pq"
	"github.com/sjwhitworth/golearn/base"
	"gonum.org/v1/gonum/mat"
)

const (
	host     = "20.50.138.179"
	port     = 5432
	user     = "postgres"
	password = "austral1a"
	dbname   = "mlearning"
)

var (
	NotEnoughDataError  = errors.New("not enough rows to support this many variables")
	NoTrainingDataError = errors.New("you need to Fit() before you can Predict")
)

type LinearRegression struct {
	Fitted                 bool             `json:"Fitted"`
	Disturbance            float64          `json:"Dist"`
	RegressionCoefficients []float64        `json:"RegCoe"`
	Attrs                  []base.Attribute `json:"Attrs"`
	Cls                    base.Attribute   `json:"Cls"`
}
type FA struct {
	Name      string
	Precision int
}

type LinearRegression2 struct {
	Fitted                 bool      `json:"Fitted"`
	Disturbance            float64   `json:"Dist"`
	RegressionCoefficients []float64 `json:"RegCoe"`
	Attrs                  []FA      `json:"Attrs"`
	Cls                    FA        `json:"Cls"`
}

func NewLinearRegression() *LinearRegression {
	return &LinearRegression{Fitted: false}
}

func mainlinear() {
	data, err := base.ParseCSVToInstances("./Advertising.csv", true) // true: means first line of csv is headers.
	//data2, err := base.ParseCSVToInstances("./Advertaising2.csv", true) // true: means first line of csv is headers.
	//fmt.Println(data2)
	if err != nil {
		fmt.Println(err)
		return
	}
	attrArray := data.AllAttributes()

	data.AddClassAttribute(attrArray[4]) //setting final column as class attribute, note that there cannot be more than one class attribute for linear regression.
	trainData, testData := base.InstancesTrainTestSplit(data, 0.015)
	//fmt.Print(testData)
	lr := NewLinearRegression()

	err = lr.Fit(trainData)
	if err != nil {
		fmt.Println(err)
	}
	predictions, err := lr.Predict(testData)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(predictions)

	c1lr := new(LinearRegression2)
	c1lr.Fitted = lr.Fitted
	c1lr.Disturbance = lr.Disturbance
	c1lr.RegressionCoefficients = make([]float64, len(lr.RegressionCoefficients))
	copy(c1lr.RegressionCoefficients, lr.RegressionCoefficients)
	c1lr.Attrs = make([]FA, len(lr.Attrs))
	for x := 0; x < len(lr.Attrs); x++ {
		c1lr.Attrs[x].Name = lr.Attrs[x].GetName()
		c1lr.Attrs[x].Precision = lr.Attrs[x].(*base.FloatAttribute).Precision
	}
	c1lr.Cls.Name = lr.Cls.GetName()
	c1lr.Cls.Precision = lr.Cls.(*base.FloatAttribute).Precision

	js, _ := json.Marshal(c1lr)

	fmt.Println("==============================")
	fmt.Println(lr)
	fmt.Println(string(js))
	fmt.Println(c1lr)
	fmt.Println("==============================")

	// Postgres
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	sqlStatement := `INSERT INTO model (id, contenido) VALUES ($1, $2)`

	_, err = db.Exec(sqlStatement, 3, js)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var js2 []byte
	sqlStatement = `SELECT CONTENIDO FROM MODEL WHERE ID=$1`
	err = db.QueryRow(sqlStatement, 1).Scan(&js2)
	if err != nil {
		panic(err)
	}

	c2lr := new(LinearRegression2)
	json.Unmarshal(js2, c2lr)
	fmt.Println("js2 unmarshaleado")
	fmt.Println(c2lr)

	lr2 := NewLinearRegression()
	lr2.Fitted = c2lr.Fitted
	lr2.Disturbance = c2lr.Disturbance
	lr2.RegressionCoefficients = make([]float64, len(c2lr.RegressionCoefficients))
	copy(lr2.RegressionCoefficients, c2lr.RegressionCoefficients)
	//lr2.Attrs = make([]base.Attribute, len(c2lr.Attrs))
	fmt.Println(len(c2lr.Attrs))
	for x := 0; x < len(c2lr.Attrs); x++ {
		lr2.Attrs = append(lr2.Attrs, &base.FloatAttribute{Name: c2lr.Attrs[x].Name, Precision: c2lr.Attrs[x].Precision})

	}
	lr2.Cls = &base.FloatAttribute{Name: c2lr.Cls.Name, Precision: c2lr.Cls.Precision}
	fmt.Println("==============================")
	fmt.Println(lr2)

	predictions2, err := lr2.Predict(testData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(predictions2)

}

func (lr *LinearRegression) Fit(inst base.FixedDataGrid) error {

	// Retrieve row size
	_, rows := inst.Size()

	// Validate class Attribute count
	classAttrs := inst.AllClassAttributes()
	if len(classAttrs) != 1 {
		return fmt.Errorf("Only 1 class variable is permitted")
	}
	classAttrSpecs := base.ResolveAttributes(inst, classAttrs)

	// Retrieve relevant Attributes
	allAttrs := base.NonClassAttributes(inst)
	attrs := make([]base.Attribute, 0)
	for _, a := range allAttrs {
		if _, ok := a.(*base.FloatAttribute); ok {
			attrs = append(attrs, a)
		}
	}

	cols := len(attrs) + 1

	if rows < cols {
		return NotEnoughDataError
	}

	// Retrieve relevant Attribute specifications
	attrSpecs := base.ResolveAttributes(inst, attrs)

	// Split into two matrices, observed results (dependent variable y)
	// and the explanatory variables (X) - see http://en.wikipedia.org/wiki/Linear_regression
	observed := mat.NewDense(rows, 1, nil)
	explVariables := mat.NewDense(rows, cols, nil)

	// Build the observed matrix
	inst.MapOverRows(classAttrSpecs, func(row [][]byte, i int) (bool, error) {
		val := base.UnpackBytesToFloat(row[0])
		observed.Set(i, 0, val)
		return true, nil
	})

	// Build the explainatory variables
	inst.MapOverRows(attrSpecs, func(row [][]byte, i int) (bool, error) {
		// Set intercepts to 1.0
		explVariables.Set(i, 0, 1.0)
		for j, r := range row {
			explVariables.Set(i, j+1, base.UnpackBytesToFloat(r))
		}
		return true, nil
	})

	n := cols
	qr := new(mat.QR)
	qr.Factorize(explVariables)
	var q, reg mat.Dense
	qr.QTo(&q)
	qr.RTo(&reg)

	var transposed, qty mat.Dense
	transposed.CloneFrom(q.T())
	qty.Mul(&transposed, observed)

	regressionCoefficients := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		regressionCoefficients[i] = qty.At(i, 0)
		for j := i + 1; j < n; j++ {
			regressionCoefficients[i] -= regressionCoefficients[j] * reg.At(i, j)
		}
		regressionCoefficients[i] /= reg.At(i, i)
	}

	lr.Disturbance = regressionCoefficients[0]
	lr.RegressionCoefficients = regressionCoefficients[1:]
	lr.Fitted = true
	lr.Attrs = attrs
	lr.Cls = classAttrs[0]
	return nil
}

func (lr *LinearRegression) Predict(X base.FixedDataGrid) (base.FixedDataGrid, error) {
	if !lr.Fitted {
		return nil, NoTrainingDataError
	}

	ret := base.GeneratePredictionVector(X)

	attrSpecs := base.ResolveAttributes(X, lr.Attrs)
	ClsSpec, err := ret.GetAttribute(lr.Cls)
	if err != nil {
		return nil, err
	}

	X.MapOverRows(attrSpecs, func(row [][]byte, i int) (bool, error) {
		var prediction float64 = lr.Disturbance
		for j, r := range row {
			prediction += base.UnpackBytesToFloat(r) * lr.RegressionCoefficients[j]
		}

		ret.Set(ClsSpec, i, base.PackFloatToBytes(prediction))
		return true, nil
	})

	return ret, nil
}
